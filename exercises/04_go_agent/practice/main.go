package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"tailscale.com/tsnet"
)

const (
	defaultTemporalHost   = "temporal-dev:7233"
	defaultApertureURL    = "http://ai"
	defaultAIModel        = "claude-haiku-4-5"
	defaultUserID         = "lab"
	defaultCheckIntervalS = "10m"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	configDir, err := os.UserConfigDir()
	if err != nil {
		logger.Error("UserConfigDir", "err", err)
		os.Exit(1)
	}

	userID := envOr("WORKSHOP_USER_ID", defaultUserID)
	hostname := fmt.Sprintf("%s-metrics-worker", userID)
	taskQueue := fmt.Sprintf("%s-health-check", userID)
	workflowID := fmt.Sprintf("%s-health-check", userID)
	scheduleID := fmt.Sprintf("%s-health-check-schedule", userID)

	tsNode := &tsnet.Server{
		Hostname: hostname,
		Dir:      filepath.Join(configDir, "workshop-tsnet", hostname),
		AuthKey:  os.Getenv("TS_AUTHKEY"),
	}
	if err := tsNode.Start(); err != nil {
		logger.Error("tsnet start", "err", err)
		os.Exit(1)
	}
	defer tsNode.Close()

	upCtx, upCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer upCancel()
	if _, err := tsNode.Up(upCtx); err != nil {
		logger.Error("tsnet up", "err", err)
		os.Exit(1)
	}
	logger.Info("joined tailnet", "hostname", hostname, "userID", userID)

	temporalHost := envOr("TEMPORAL_HOST", defaultTemporalHost)
	c, err := client.Dial(client.Options{
		HostPort: temporalHost,
		Logger:   logger,
		ConnectionOptions: client.ConnectionOptions{
			DialOptions: []grpc.DialOption{
				grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
					host, _, _ := net.SplitHostPort(addr)
					if host == "localhost" || host == "127.0.0.1" {
						return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
					}
					return tsNode.Dial(ctx, "tcp", addr)
				}),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithKeepaliveParams(keepalive.ClientParameters{
					Time:                30 * time.Second,
					Timeout:             10 * time.Second,
					PermitWithoutStream: true,
				}),
			},
		},
	})
	if err != nil {
		logger.Error("temporal dial", "err", err)
		os.Exit(1)
	}
	defer c.Close()
	logger.Info("connected to temporal", "host", temporalHost)

	metricsURL := mustEnv(logger, "METRICS_URL")
	acts := NewActivities(
		tsNode.HTTPClient(),
		metricsURL,
		envOr("APERTURE_URL", defaultApertureURL),
		envOr("AI_MODEL", defaultAIModel),
	)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()
	if sample, err := acts.FetchMetrics(pingCtx); err != nil {
		logger.Warn("metrics endpoint unreachable", "url", metricsURL, "err", err)
	} else {
		logger.Info("metrics reachable", "url", metricsURL, "sample", strings.SplitN(sample, "\n", 2)[0])
	}

	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(HealthCheckWorkflow)
	w.RegisterActivity(acts)

	interval := healthCheckInterval(logger)

	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cleanupCancel()
	if err := c.ScheduleClient().GetHandle(cleanupCtx, scheduleID).Delete(cleanupCtx); err != nil {
		logger.Debug("no existing schedule to delete (ok)", "id", scheduleID, "err", err)
	}

	schedCtx, schedCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer schedCancel()
	_, schedErr := c.ScheduleClient().Create(schedCtx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{{Every: interval}},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  HealthCheckWorkflow,
			TaskQueue: taskQueue,
		},
		TriggerImmediately: true,
		Overlap:            enumspb.SCHEDULE_OVERLAP_POLICY_SKIP,
	})
	if schedErr != nil {
		logger.Error("create schedule", "id", scheduleID, "err", schedErr)
		os.Exit(1)
	}
	logger.Info("created schedule", "id", scheduleID, "interval", interval.String(), "workflow", workflowID)

	logger.Info("worker running", "taskQueue", taskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		logger.Error("worker stopped with error", "err", err)
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(logger *slog.Logger, key string) string {
	v := os.Getenv(key)
	if v == "" {
		logger.Error("required env var not set", "key", key)
		os.Exit(1)
	}
	return v
}

func healthCheckInterval(logger *slog.Logger) time.Duration {
	raw := envOr("HEALTH_CHECK_INTERVAL", defaultCheckIntervalS)
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		logger.Warn("invalid HEALTH_CHECK_INTERVAL, using default", "value", raw, "default", defaultCheckIntervalS)
		fallback, _ := time.ParseDuration(defaultCheckIntervalS)
		return fallback
	}
	return d
}
