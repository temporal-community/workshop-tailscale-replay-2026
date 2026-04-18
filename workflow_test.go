package workshop_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

	workshop "github.com/temporal-community/workshop-tailscale-replay-2026"
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
	var act *workshop.Activities
	env.OnActivity(act.FetchMetrics, mock.Anything).Return("# metrics", nil)
	env.OnActivity(act.AnalyzeMetrics, mock.Anything, mock.Anything).Return("System healthy.", nil)

	env.ExecuteWorkflow(workshop.HealthCheckWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Equal("System healthy.", result)
}

func (s *workflowSuite) Test_FetchError_Propagates() {
	env := s.NewTestWorkflowEnvironment()
	var act *workshop.Activities
	env.OnActivity(act.FetchMetrics, mock.Anything).Return("", fmt.Errorf("connection refused"))

	env.ExecuteWorkflow(workshop.HealthCheckWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
}
