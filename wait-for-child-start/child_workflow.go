package waitforchildstart

import (
	"go.temporal.io/sdk/workflow"
)

type GreetingParams struct {
	Message string
	Name    string
}

func WaitingChildWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Debug("Started")
	var greetingParams GreetingParams
	workflow.GetSignalChannel(ctx, "SendGreeting").Receive(ctx, &greetingParams)
	greeting := greetingParams.Message + " " + greetingParams.Name + "!"
	logger.Info("Child workflow execution: " + greeting)
	return greeting, nil
}
