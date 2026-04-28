package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/http"
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
	defaultCheckIntervalS = "1m"
)

func main() {
	mode := "worker"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	userID := os.Getenv("WORKSHOP_USER_ID")
	if userID == "" {
		logger.Error("WORKSHOP_USER_ID is not set. Open a new terminal or run `source ~/.bashrc`. Instruqt sets this automatically for all workshop shells.")
		os.Exit(1)
	}
	taskQueue := fmt.Sprintf("%s-health-check", userID)
	workflowID := fmt.Sprintf("%s-health-check", userID)
	scheduleID := fmt.Sprintf("%s-health-check-schedule", userID)

	switch mode {
	case "worker":
		runWorker(logger, userID, taskQueue)
	case "starter":
		runStarter(logger, userID, taskQueue, workflowID, scheduleID)
	default:
		logger.Error("unknown mode", "mode", mode, "expected", "'worker' or 'starter'")
		os.Exit(1)
	}
}

func runWorker(logger *slog.Logger, userID, taskQueue string) {
	tsNode := startTsnet(logger, userID, "worker")
	defer tsNode.Close()

	c := dialTemporal(logger, tsNode)
	defer c.Close()

	metricsURL := mustEnv(logger, "METRICS_URL")
	acts := NewActivities(
		tsNode.HTTPClient(),
		metricsURL,
		envOr("APERTURE_URL", defaultApertureURL),
		envOr("AI_MODEL", defaultAIModel),
	)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()
	if req, err := http.NewRequestWithContext(pingCtx, http.MethodGet, metricsURL, nil); err == nil {
		if resp, err := acts.HTTP.Do(req); err != nil {
			logger.Warn("metrics endpoint unreachable", "url", metricsURL, "err", err)
		} else {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			logger.Info("metrics reachable", "url", metricsURL, "sample", strings.SplitN(string(body), "\n", 2)[0])
		}
	}

	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(HealthCheckWorkflow)
	w.RegisterActivity(acts)

	logger.Info("worker running", "taskQueue", taskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		logger.Error("worker stopped with error", "err", err)
		os.Exit(1)
	}
}

func runStarter(logger *slog.Logger, userID, taskQueue, workflowID, scheduleID string) {
	tsNode := startTsnet(logger, userID, "starter")
	defer tsNode.Close()

	c := dialTemporal(logger, tsNode)
	defer c.Close()

	interval := healthCheckInterval(logger)

	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cleanupCancel()
	if err := c.ScheduleClient().GetHandle(cleanupCtx, scheduleID).Delete(cleanupCtx); err != nil {
		logger.Debug("no existing schedule to delete (ok)", "id", scheduleID, "err", err)
	}

	schedCtx, schedCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer schedCancel()
	_, err := c.ScheduleClient().Create(schedCtx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{{Every: interval}},
		},
		RemainingActions: 5,
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  HealthCheckWorkflow,
			TaskQueue: taskQueue,
		},
		TriggerImmediately: true,
		Overlap:            enumspb.SCHEDULE_OVERLAP_POLICY_SKIP,
	})
	if err != nil {
		logger.Error("create schedule", "id", scheduleID, "err", err)
		os.Exit(1)
	}
	logger.Info("created schedule", "id", scheduleID, "interval", interval.String(), "workflow", workflowID, "remainingActions", 5)
}

func startTsnet(logger *slog.Logger, userID, mode string) *tsnet.Server {
	configDir, err := os.UserConfigDir()
	if err != nil {
		logger.Error("UserConfigDir", "err", err)
		os.Exit(1)
	}

	hostname, err := resolveNodeName(configDir, userID, mode)
	if err != nil {
		logger.Error("resolve node name", "err", err)
		os.Exit(1)
	}

	tsNode := &tsnet.Server{
		Hostname: hostname,
		Dir:      filepath.Join(configDir, "workshop-tsnet", hostname),
		AuthKey:  os.Getenv("TS_AUTHKEY"),
	}
	if err := tsNode.Start(); err != nil {
		logger.Error("tsnet start", "err", err)
		os.Exit(1)
	}

	upCtx, upCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer upCancel()
	if _, err := tsNode.Up(upCtx); err != nil {
		logger.Error("tsnet up", "err", err)
		os.Exit(1)
	}
	logger.Info("joined tailnet", "hostname", hostname, "userID", userID)

	return tsNode
}

func dialTemporal(logger *slog.Logger, tsNode *tsnet.Server) client.Client {
	temporalHost := envOr("TEMPORAL_HOST", defaultTemporalHost)
	c, err := client.Dial(client.Options{
		HostPort: "passthrough:///" + temporalHost,
		Logger:   logger,
		ConnectionOptions: client.ConnectionOptions{
			GetSystemInfoTimeout: 30 * time.Second,
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
	logger.Info("connected to temporal", "host", temporalHost)
	return c
}

// resolveNodeName returns a stable, per-machine node name of the form
// "<userID>-ex4-metrics-<mode>-<suffix>". The 5-char lowercase-alpha
// suffix is generated once on first run and then reused on every
// subsequent run (found by scanning workshop-tsnet/ for an existing
// dir with the same prefix). Worker and starter get their own nodes
// on the tailnet so you can see both of them in `tailscale status`.
func resolveNodeName(configDir, userID, mode string) (string, error) {
	root := filepath.Join(configDir, "workshop-tsnet")
	prefix := fmt.Sprintf("%s-ex4-metrics-%s-", userID, mode)

	entries, err := os.ReadDir(root)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), prefix) {
			return e.Name(), nil
		}
	}

	const letters = "abcdefghijklmnopqrstuvwxyz"
	suffix := make([]byte, 5)
	for i := range suffix {
		suffix[i] = letters[rand.IntN(len(letters))]
	}
	return prefix + string(suffix), nil
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
