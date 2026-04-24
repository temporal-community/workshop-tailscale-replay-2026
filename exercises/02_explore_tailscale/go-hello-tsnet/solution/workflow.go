package main

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func GetAddressFromIP(ctx workflow.Context, input WorkflowInput) (WorkflowOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ip string
	if err := workflow.ExecuteActivity(ctx, GetIP).Get(ctx, &ip); err != nil {
		return WorkflowOutput{}, fmt.Errorf("get ip: %w", err)
	}

	var location string
	if err := workflow.ExecuteActivity(ctx, GetLocationInfo, ip).Get(ctx, &location); err != nil {
		return WorkflowOutput{}, fmt.Errorf("get location: %w", err)
	}

	return WorkflowOutput{IPAddr: ip, Location: location}, nil
}
