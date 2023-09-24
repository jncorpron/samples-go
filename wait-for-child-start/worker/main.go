package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	waitforchildstart "github.com/temporalio/samples-go/wait-for-child-start"
)

func main() {
	// The client is a heavyweight object that should be created only once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "waitforchildstart", worker.Options{})

	w.RegisterWorkflow(waitforchildstart.WaitingParentWorkflow)
	w.RegisterWorkflow(waitforchildstart.WaitingChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
