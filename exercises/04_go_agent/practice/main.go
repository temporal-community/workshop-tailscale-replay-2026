// ABOUTME: Worker and starter for the Go weather agent exercise.
// ABOUTME: Connects to the shared Temporal server on the Tailscale network.

package main

import (
	"context"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// TODO: Connect to the Temporal server on the Tailscale network.
	// Use the TEMPORAL_ADDRESS environment variable.
	address := os.Getenv("TEMPORAL_ADDRESS")
	if address == "" {
		address = "localhost:7233"
	}

	userID := os.Getenv("WORKSHOP_USER_ID")
	if userID == "" {
		userID = "unknown"
	}

	taskQueue := userID + "-go-agent"

	c, err := client.Dial(client.Options{
		HostPort: address,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	if len(os.Args) > 1 && os.Args[1] == "run" {
		// Start a workflow
		runWorkflow(context.Background(), c, userID, taskQueue)
		return
	}

	// Start the worker
	w := worker.New(c, taskQueue, worker.Options{})

	// TODO: Register the AgentWorkflow and activities
	// w.RegisterWorkflow(AgentWorkflow)
	// w.RegisterActivity(&Activities{})

	log.Printf("Starting Go agent worker on task queue: %s\n", taskQueue)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func runWorkflow(ctx context.Context, c client.Client, userID, taskQueue string) {
	// TODO: Start the AgentWorkflow
	log.Println("Go agent workflow execution not yet implemented")
}
