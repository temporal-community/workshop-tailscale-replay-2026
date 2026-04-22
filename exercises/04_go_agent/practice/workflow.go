package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func HealthCheckWorkflow(ctx workflow.Context) (HealthReport, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	})

	var act *Activities

	var metrics string
	if err := workflow.ExecuteActivity(ctx, act.FetchMetrics).Get(ctx, &metrics); err != nil {
		return HealthReport{}, err
	}

	var report HealthReport
	if err := workflow.ExecuteActivity(ctx, act.AnalyzeMetrics, metrics).Get(ctx, &report); err != nil {
		return HealthReport{}, err
	}

	return report, nil
}
