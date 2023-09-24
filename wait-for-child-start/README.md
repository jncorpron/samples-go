# A demo of the GetChildWorkflowExecution issue

When a child workflow is mocked using OnWorkflow, GetChildWorkflowExecution's
does not work as expected. The returned future, when Get is called appears to
block until the mocked child workflow is completely resolved to a return value,
ending the completing the child workflow's execution. This completion of the
mocked child workflow includes the testsuite skipping time forward if the mocked
child workflow set any wait duration with After or AfterFn. The completed
child workflow then cannot be found when signalled, even thought it is expected
to be running after the GetChildWorkflowExecution future resolves.

## Expected Behavior

GetChildWorkflowExecution() is documented as "returns a future that will be
ready when child workflow execution started". Our expectation is that even if
OnWorkflow mock is used, that the wait duration does not skip forward when the
GetChildWorkflowExecution() returned future is resolved in the parent.

## Reproduction Tests

parent_workflow.go contains tests demonstrating the various behaviors, of
mocking the child or not, using After or not, and finally calling
GetChildWorkflowExecution or not. All tests are expected to nominally PASS if
but instead the tests Test_WaitForChild_ChildMocked_After2Hours and
Test_WaitForChild_ChildMocked_NoAfter do in fact currently FAIL due to the
unexpected behavior of GetChildWorkflowExecution for mocked child workflows.

The following is an example test run demonstrating how the test failures
manifest.

```text
Running tool: /bin/go test -timeout 600s -run ^(Test_WaitForChild_ChildRegistered|Test_WaitForChild_ChildMocked_After2Hours|Test_WaitForChild_ChildMocked_NoAfter|Test_Impatient_ChildRegistered|Test_Impatient_ChildMocked)$ github.com/temporalio/samples-go/wait-for-child-start -race -v

=== RUN   Test_WaitForChild_ChildRegistered
2023/09/23 23:59:59 INFO  ExecuteChildWorkflow WorkflowType WaitingChildWorkflow
2023/09/23 23:59:59 DEBUG Waiting for child to start
2023/09/23 23:59:59 DEBUG Started
2023/09/23 23:59:59 DEBUG Auto fire timer TimerID 0 TimerDuration 1h0m0s TimeSkipped 1h0m0s
2023/09/23 23:59:59 DEBUG Got ProvideName signal
2023/09/23 23:59:59 DEBUG Signalled child
2023/09/23 23:59:59 INFO  Child workflow execution: Hello Hallie Wade!
2023/09/23 23:59:59 INFO  Parent execution completed. Result Hello Hallie Wade!
--- PASS: Test_WaitForChild_ChildRegistered (0.00s)
=== RUN   Test_WaitForChild_ChildMocked_After2Hours
2023/09/23 23:59:59 INFO  ExecuteChildWorkflow WorkflowType WaitingChildWorkflow
2023/09/23 23:59:59 DEBUG Waiting for child to start
2023/09/23 23:59:59 DEBUG Auto fire timer TimerID 0 TimerDuration 1h0m0s TimeSkipped 1h0m0s
2023/09/23 23:59:59 DEBUG Auto fire timer TimerID 2 TimerDuration 2h0m0s TimeSkipped 1h0m0s
    go-samples/wait-for-child-start/parent_workflow_test.go:52: Child returning
2023/09/23 23:59:59 DEBUG Got ProvideName signal
    go-samples/wait-for-child-start/parent_workflow_test.go:65:
          Error Trace:  go-samples/wait-for-child-start/parent_workflow_test.go:65
          Error:        Received unexpected error:
                        workflow execution error (type: WaitingParentWorkflow, workflowID: default-test-workflow-id, runID: default-test-run-id): unknown external workflow execution (type: UnknownExternalWorkflowExecutionError, retryable: true)
          Test:         Test_WaitForChild_ChildMocked_After2Hours
--- FAIL: Test_WaitForChild_ChildMocked_After2Hours (0.00s)
=== RUN   Test_WaitForChild_ChildMocked_NoAfter
2023/09/23 23:59:59 INFO  ExecuteChildWorkflow WorkflowType WaitingChildWorkflow
2023/09/23 23:59:59 DEBUG Waiting for child to start
    go-samples/wait-for-child-start/parent_workflow_test.go:78: Child returning
2023/09/23 23:59:59 DEBUG Auto fire timer TimerID 0 TimerDuration 1h0m0s TimeSkipped 1h0m0s
2023/09/23 23:59:59 DEBUG Got ProvideName signal
    go-samples/wait-for-child-start/parent_workflow_test.go:91:
          Error Trace:  go-samples/wait-for-child-start/parent_workflow_test.go:91
          Error:        Received unexpected error:
                        workflow execution error (type: WaitingParentWorkflow, workflowID: default-test-workflow-id, runID: default-test-run-id): unknown external workflow execution (type: UnknownExternalWorkflowExecutionError, retryable: true)
          Test:         Test_WaitForChild_ChildMocked_NoAfter
--- FAIL: Test_WaitForChild_ChildMocked_NoAfter (0.00s)
=== RUN   Test_Impatient_ChildRegistered
2023/09/23 23:59:59 INFO  ExecuteChildWorkflow WorkflowType WaitingChildWorkflow
2023/09/23 23:59:59 DEBUG Not waiting for child to start
2023/09/23 23:59:59 DEBUG Started
2023/09/23 23:59:59 DEBUG Auto fire timer TimerID 0 TimerDuration 1h0m0s TimeSkipped 1h0m0s
2023/09/23 23:59:59 DEBUG Got ProvideName signal
2023/09/23 23:59:59 DEBUG Signalled child
2023/09/23 23:59:59 INFO  Child workflow execution: Hello Hallie Wade!
2023/09/23 23:59:59 INFO  Parent execution completed. Result Hello Hallie Wade!
--- PASS: Test_Impatient_ChildRegistered (0.00s)
=== RUN   Test_Impatient_ChildMocked
2023/09/23 23:59:59 INFO  ExecuteChildWorkflow WorkflowType WaitingChildWorkflow
2023/09/23 23:59:59 DEBUG Not waiting for child to start
2023/09/23 23:59:59 DEBUG Auto fire timer TimerID 0 TimerDuration 1h0m0s TimeSkipped 1h0m0s
2023/09/23 23:59:59 DEBUG Got ProvideName signal
2023/09/23 23:59:59 DEBUG Signalled child
2023/09/23 23:59:59 DEBUG Auto fire timer TimerID 2 TimerDuration 2h0m0s TimeSkipped 1h0m0s
    go-samples/wait-for-child-start/parent_workflow_test.go:124: Child returning
2023/09/23 23:59:59 INFO  Workflow has unhandled signals SignalNames [SendGreeting]
2023/09/23 23:59:59 INFO  Parent execution completed. Result Greetings to You!
--- PASS: Test_Impatient_ChildMocked (0.00s)
FAIL
FAIL  github.com/temporalio/samples-go/wait-for-child-start 0.880s
```

## Workaround

We currently utilize the following test helper which disables calling
GetChildWorkflowExecution, and instead does a sleep (allows the mocked child
execution to start in the testsuite environment).

```golang
var disableChildWorkflowExecutionCheck = false

// DisableChildWorkflowExecutionCheck causes WaitForChildToStart to not use GetChildWorkflowExecution
// for the duration of a test run.
func DisableChildWorkflowExecutionCheck(tb testing.TB) {
  tb.Helper()
  disableChildWorkflowExecutionCheck = true
  tb.Cleanup(func() {
    disableChildWorkflowExecutionCheck = false
  })
}

// WaitForChildToStart waits until a child workflow has started. This wait can be disabled if needed
func WaitForChildToStart(ctx workflow.Context, childFuture workflow.ChildWorkflowFuture) error {
  if disableChildWorkflowExecutionCheck {
    // instead we sleep 1 second to allow the test environment to start the child workflow.
    return workflow.Sleep(ctx, 1*time.Second)
  }
  return childFuture.GetChildWorkflowExecution().Get(ctx, nil)
}
```
