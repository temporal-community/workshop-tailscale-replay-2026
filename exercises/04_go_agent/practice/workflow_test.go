package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

func TestHealthCheckWorkflow(t *testing.T) {
	suite.Run(t, new(workflowSuite))
}

type workflowSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func (s *workflowSuite) Test_Completes() {
	env := s.NewTestWorkflowEnvironment()
	var act *Activities
	env.OnActivity(act.FetchMetrics, mock.Anything).Return("# metrics", nil)
	env.OnActivity(act.AnalyzeMetrics, mock.Anything, mock.Anything).Return(HealthReport{Summary: "healthy"}, nil)

	env.ExecuteWorkflow(HealthCheckWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	var result HealthReport
	s.NoError(env.GetWorkflowResult(&result))
	s.Equal("healthy", result.Summary)
}

func (s *workflowSuite) Test_FetchError_Propagates() {
	env := s.NewTestWorkflowEnvironment()
	var act *Activities
	env.OnActivity(act.FetchMetrics, mock.Anything).Return("", fmt.Errorf("connection refused"))

	env.ExecuteWorkflow(HealthCheckWorkflow)

	s.True(env.IsWorkflowCompleted())
	err := env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "connection refused")
}
