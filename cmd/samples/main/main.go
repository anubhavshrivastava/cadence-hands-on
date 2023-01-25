package main

import (
	"flag"
	"time"

	"github.com/anubhavshrivastava/cadence-hands-on/cmd/samples/common"
	"github.com/anubhavshrivastava/cadence-hands-on/cmd/samples/ex_1"
	"github.com/anubhavshrivastava/cadence-hands-on/cmd/samples/ex_2"
	"github.com/anubhavshrivastava/cadence-hands-on/cmd/samples/hello_world_workflow"
	"github.com/google/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
)

var helper common.SampleHelper

func main() {

	var command, workflowName, workflowId, signalInput string
	flag.StringVar(&command, "cmd", "start-workflow", "Command is start-workflow, start-worker, signal-workflow")
	flag.StringVar(&workflowName, "wf", "hello-world", "Workflow Name to start")
	flag.StringVar(&workflowId, "wfId", "", "Workflow Id to Signal")
	flag.StringVar(&signalInput, "sigData", "", "Signal Data")
	flag.Parse()

	helper.SetupServiceConfig()

	switch command {
	case "start-workflow":
		startWorkflow(workflowName, "anubhav")
	case "start-worker":
		startWorker(workflowName)
		select {}
	case "signal-workflow":
		signalWorkflow(workflowId, signalInput)
		select {}
	}

}
func signalWorkflow(workflowId, signalInput string) {
	helper.SignalWorkflow(workflowId, "testSignalChannel", signalInput)
}

func startWorker(workflowName string) {
	switch workflowName {
	case "hello-world":
		helper.RegisterWorkflow(hello_world_workflow.WorkflowFunc)
		helper.RegisterActivity(hello_world_workflow.ActivityOne)

		workerOptions := worker.Options{
			MetricsScope: helper.WorkerMetricScope,
			Logger:       helper.Logger,
			FeatureFlags: client.FeatureFlags{
				WorkflowExecutionAlreadyCompletedErrorEnabled: true,
			},
		}
		helper.StartWorkers(helper.Config.DomainName, "default_task_list", workerOptions)
	case "ex-1":
		helper.RegisterWorkflow(ex_1.Workflow)
		helper.RegisterActivity(ex_1.SayWorld)
		helper.RegisterActivity(ex_1.SayHello)

		workerOptions := worker.Options{
			MetricsScope: helper.WorkerMetricScope,
			Logger:       helper.Logger,
			FeatureFlags: client.FeatureFlags{
				WorkflowExecutionAlreadyCompletedErrorEnabled: true,
			},
		}
		helper.StartWorkers(helper.Config.DomainName, "default_task_list", workerOptions)
	case "ex-2":
		helper.RegisterWorkflow(ex_2.Workflow)
		helper.RegisterActivity(ex_2.SayHelloToSignalValue)

		workerOptions := worker.Options{
			MetricsScope: helper.WorkerMetricScope,
			Logger:       helper.Logger,
			FeatureFlags: client.FeatureFlags{
				WorkflowExecutionAlreadyCompletedErrorEnabled: true,
			},
		}
		helper.StartWorkers(helper.Config.DomainName, "default_task_list", workerOptions)
	}

}

func startWorkflow(workflowName string, inputs ...interface{}) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                              workflowName + "_" + uuid.New().String(),
		TaskList:                        "default_task_list",
		ExecutionStartToCloseTimeout:    1000 * time.Minute,
		DecisionTaskStartToCloseTimeout: 100 * time.Minute,
		RetryPolicy: &workflow.RetryPolicy{
			InitialInterval:          1 * time.Second,
			BackoffCoefficient:       1,
			MaximumInterval:          0,
			ExpirationInterval:       0,
			MaximumAttempts:          10,
			NonRetriableErrorReasons: nil,
		},
	}

	switch workflowName {
	case "hello-world":
		{
			helper.RegisterWorkflow(hello_world_workflow.WorkflowFunc)
			helper.RegisterActivity(hello_world_workflow.ActivityOne)
			helper.StartWorkflow(workflowOptions, hello_world_workflow.WorkflowFunc, inputs...)
		}

	case "ex-1":
		{
			helper.RegisterWorkflow(ex_1.Workflow)
			helper.RegisterActivity(ex_1.SayHello)
			helper.RegisterActivity(ex_1.SayWorld)
			helper.StartWorkflow(workflowOptions, ex_1.Workflow, inputs...)
		}
	case "ex-2":
		{
			helper.RegisterWorkflow(ex_2.Workflow)
			helper.RegisterActivity(ex_2.SayHelloToSignalValue)
			helper.StartWorkflow(workflowOptions, ex_2.Workflow)
		}
	}
}
