// ABOUTME: Worker + starter for Exercise 2 Part 3 — hello-tsnet Go worker.
// ABOUTME: TODOs walk through joining the tailnet via tsnet and dialing Temporal through it.

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"tailscale.com/tsnet"
)

const temporalHost = "temporal-dev:7233"

func main() {
	mode := "worker"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	userID := os.Getenv("WORKSHOP_USER_ID")
	if userID == "" {
		log.Fatal("WORKSHOP_USER_ID is not set — source your shell profile or export it")
	}
	taskQueue := fmt.Sprintf("%s-hello-tsnet", userID)

	tsNode := startTsnet(userID, mode)
	defer tsNode.Close()

	c := dialTemporal(tsNode)
	defer c.Close()

	switch mode {
	case "worker":
		runWorker(c, taskQueue)
	case "starter":
		runStarter(c, userID, taskQueue)
	default:
		log.Fatalf("unknown mode %q (expected 'worker' or 'starter')", mode)
	}
}

func startTsnet(userID, mode string) *tsnet.Server {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("user config dir: %v", err)
	}
	_ = configDir // used by TODO 1b below

	nodeName, err := resolveNodeName(configDir, userID, mode)
	if err != nil {
		log.Fatalf("resolve node name: %v", err)
	}
	_ = nodeName // used by TODO 1a and 1b below

	tsNode := &tsnet.Server{
		// TODO 1a: Set Hostname to nodeName
		//          — this is the name your node will have on the tailnet.
		Hostname:

		// TODO 1b: Set Dir to filepath.Join(configDir, "workshop-tsnet", nodeName)
		//          — tsnet stores its node key here so later runs reuse the identity.
		Dir:

		// TODO 1c: Set AuthKey to os.Getenv("TS_AUTHKEY")
		//          — consumed once on first run to register the node.
		AuthKey:
	}
	if err := tsNode.Start(); err != nil {
		log.Fatalf("tsnet start: %v", err)
	}

	upCtx, upCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer upCancel()
	if _, err := tsNode.Up(upCtx); err != nil {
		log.Fatalf("tsnet up: %v", err)
	}
	log.Printf("joined tailnet as %s", nodeName)
	return tsNode
}

// resolveNodeName returns a stable, per-machine node name of the form
// "<userID>-ex2-go-<mode>-<suffix>". The 5-char lowercase-alpha suffix is
// generated once on first run and then reused on every subsequent run
// (found by scanning workshop-tsnet/ for an existing dir with the same
// prefix). Two attendees with the same WORKSHOP_USER_ID get different
// suffixes, so their tailnet hostnames don't collide.
func resolveNodeName(configDir, userID, mode string) (string, error) {
	root := filepath.Join(configDir, "workshop-tsnet")
	prefix := fmt.Sprintf("%s-ex2-go-%s-", userID, mode)

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

// ============================================================================
// TODO 2: Teach the Temporal gRPC client to dial through tsNode.
//
// The SDK opens a TCP connection to temporalHost. We want that connection to
// go over the tailnet, not the VM's normal network stack. Plug tsNode.Dial in
// as a gRPC ContextDialer:
//
//     grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
//         return tsNode.Dial(ctx, "tcp", addr)
//     })
//
// Add that option to the dialOptions slice below, then the existing
// client.Dial call will work.
// ============================================================================
func dialTemporal(tsNode *tsnet.Server) client.Client {
	dialOptions := []grpc.DialOption{
		// TODO 2: insert grpc.WithContextDialer(...) here
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	c, err := client.Dial(client.Options{
		HostPort: temporalHost,
		ConnectionOptions: client.ConnectionOptions{
			DialOptions: dialOptions,
		},
	})
	if err != nil {
		log.Fatalf("temporal dial: %v", err)
	}
	log.Printf("connected to temporal at %s via tsnet", temporalHost)
	return c
}

func runWorker(c client.Client, taskQueue string) {
	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(GetAddressFromIP)
	w.RegisterActivity(GetIP)
	w.RegisterActivity(GetLocationInfo)

	log.Printf("Starting Go worker on task queue: %s", taskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("worker stopped: %v", err)
	}
}

func runStarter(c client.Client, userID, taskQueue string) {
	workflowID := fmt.Sprintf("%s-hello-tsnet-%d", userID, os.Getpid())
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQueue,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options,
		GetAddressFromIP,
		WorkflowInput{Name: userID},
	)
	if err != nil {
		log.Fatalf("execute workflow: %v", err)
	}
	log.Printf("Started workflow %s (run %s)", we.GetID(), we.GetRunID())

	var result WorkflowOutput
	if err := we.Get(context.Background(), &result); err != nil {
		log.Fatalf("workflow result: %v", err)
	}
	log.Printf("Result: IP=%s  Location=%s", result.IPAddr, result.Location)
}

// Keep imports alive for the code attendees will write inside the TODOs.
var (
	_ = filepath.Join
	_ = (*net.TCPAddr)(nil)
)
