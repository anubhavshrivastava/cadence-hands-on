package hello_world_workflow

import (
	"context"
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func WorkflowFunc(ctx workflow.Context, name string) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: 1 * time.Minute,
		StartToCloseTimeout:    10 * time.Second,
		HeartbeatTimeout:       time.Minute * 2,
		TaskList:               "default_task_list",
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("Test Workflow  workflow started")

	var resp string
	logger.Info("Executing Activity")
	err := workflow.ExecuteActivity(ctx, ActivityOne, name).Get(ctx, &resp)
	if err != nil {
		logger.Error("Activity  Failed.", zap.Error(err))
		return err
	}

	logger.Info("Activity Response " + resp)

	return nil
}

func ActivityOne(ctx context.Context, name string) (string, error) {
	return "Hello Mr." + name, nil
}
