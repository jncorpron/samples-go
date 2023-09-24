package waitforchildstart

import (
	"go.temporal.io/sdk/workflow"
)

type NameParams struct {
	Name string
}

type WaitingParentWorkflowParams struct {
	WaitForChildToStart bool
}

func WaitingParentWorkflow(ctx workflow.Context, params WaitingParentWorkflowParams) (string, error) {
	logger := workflow.GetLogger(ctx)

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "ABC-SIMPLE-CHILD-WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	childFuture := workflow.ExecuteChildWorkflow(ctx, WaitingChildWorkflow)

	if params.WaitForChildToStart {
		logger.Debug("Waiting for child to start")
		// wait for the child to start
		err := childFuture.GetChildWorkflowExecution().Get(ctx, nil)
		if err != nil {
			return "", err
		}
	} else {
		logger.Debug("Not waiting for child to start")
	}

	// wait for a signal that provides the name for the greeting
	var nameParams NameParams
	_ = workflow.GetSignalChannel(ctx, "ProvideName").Receive(ctx, &nameParams)

	logger.Debug("Got ProvideName signal")

	// signal the child to provide the name so that it returns the greeting message.
	err := workflow.SignalExternalWorkflow(
		ctx,
		"ABC-SIMPLE-CHILD-WORKFLOW-ID",
		"",
		"SendGreeting",
		GreetingParams{Message: "Hello", Name: nameParams.Name},
	).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	logger.Debug("Signalled child")

	var result string
	err = childFuture.Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return "", err
	}

	logger.Info("Parent execution completed.", "Result", result)
	return result, nil
}
