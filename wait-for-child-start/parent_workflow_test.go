package waitforchildstart_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	waitforchildstart "github.com/temporalio/samples-go/wait-for-child-start"
)

// Test_WaitForChild_ChildRegistered runs a test where the child workflow
// WaitingChildWorkflow is NOT mocked, and is run by registering the child
// workflow.
//
// This test will PASS as expected, demonstrating the expected behavior.
func Test_WaitForChild_ChildRegistered(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(waitforchildstart.WaitingChildWorkflow)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("ProvideName", waitforchildstart.NameParams{Name: "Hallie Wade"})
	}, 1*time.Hour)

	env.ExecuteWorkflow(waitforchildstart.WaitingParentWorkflow, waitforchildstart.WaitingParentWorkflowParams{
		WaitForChildToStart: true,
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Hello Hallie Wade!", result)
}

// Test_WaitForChild_ChildRegistered runs a test where the child workflow
// WaitingChildWorkflow is NOT mocked, and is run by registering the child
// workflow.
//
// This test will PASS as expected, demonstrating the expected behavior.
func Test_WaitForChild_ChildMocked_After2Hours(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.OnWorkflow(waitforchildstart.WaitingChildWorkflow, mock.Anything).
		Once().
		After(2 * time.Hour).
		Return(func(ctx workflow.Context) (string, error) {
			t.Log("Child returning")
			return "Greetings to You!", nil
		})

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("ProvideName", waitforchildstart.NameParams{Name: "Hallie Wade"})
	}, 1*time.Hour)

	env.ExecuteWorkflow(waitforchildstart.WaitingParentWorkflow, waitforchildstart.WaitingParentWorkflowParams{
		WaitForChildToStart: true,
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Greetings to You!", result)
}

func Test_WaitForChild_ChildMocked_NoAfter(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.OnWorkflow(waitforchildstart.WaitingChildWorkflow, mock.Anything).
		Once().
		// No After
		Return(func(ctx workflow.Context) (string, error) {
			t.Log("Child returning")
			return "Greetings to You!", nil
		})

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("ProvideName", waitforchildstart.NameParams{Name: "Hallie Wade"})
	}, 1*time.Hour)

	env.ExecuteWorkflow(waitforchildstart.WaitingParentWorkflow, waitforchildstart.WaitingParentWorkflowParams{
		WaitForChildToStart: true,
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Greetings to You!", result)
}

func Test_Impatient_ChildRegistered(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(waitforchildstart.WaitingChildWorkflow)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("ProvideName", waitforchildstart.NameParams{Name: "Hallie Wade"})
	}, 1*time.Hour)

	env.ExecuteWorkflow(waitforchildstart.WaitingParentWorkflow, waitforchildstart.WaitingParentWorkflowParams{
		WaitForChildToStart: false,
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Hello Hallie Wade!", result)
}

func Test_Impatient_ChildMocked(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.OnWorkflow(waitforchildstart.WaitingChildWorkflow, mock.Anything).
		Once().
		After(2 * time.Hour).
		Return(func(ctx workflow.Context) (string, error) {
			t.Log("Child returning")
			return "Greetings to You!", nil
		})

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("ProvideName", waitforchildstart.NameParams{Name: "Hallie Wade"})
	}, 1*time.Hour)

	env.ExecuteWorkflow(waitforchildstart.WaitingParentWorkflow, waitforchildstart.WaitingParentWorkflowParams{
		WaitForChildToStart: false,
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "Greetings to You!", result)
}
