// ABOUTME: Shared types and constants for the Go hello-tsnet workshop exercise.
// ABOUTME: Input/output structs for the GetAddressFromIP workflow.

package main

type WorkflowInput struct {
	Name string
}

type WorkflowOutput struct {
	IPAddr   string
	Location string
}
