package ex_1

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func opts() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		ScheduleToStartTimeout: 1 * time.Minute,
		StartToCloseTimeout:    10 * time.Second,
		HeartbeatTimeout:       time.Minute * 2,
		TaskList:               "default_task_list",
	}
}

func Workflow(ctx workflow.Context, name string) error {
	ctx = workflow.WithActivityOptions(ctx, opts())
	logger := workflow.GetLogger(ctx)
	logger.Info("Test Workflow  workflow started")

	var helloStr string
	logger.Info("Executing Activity")
	err := workflow.ExecuteActivity(ctx, SayHello).Get(ctx, &helloStr)
	if err != nil {
		logger.Error("Activity  Failed.", zap.Error(err))
		return err
	}

	var worldStr string
	logger.Info("Executing Activity")
	err = workflow.ExecuteActivity(ctx, SayWorld).Get(ctx, &worldStr)
	if err != nil {
		logger.Error("Activity  Failed.", zap.Error(err))
		return err
	}

	finalStr := fmt.Sprintf("%s %s Mr. %s ..... This is Ex 1", helloStr, worldStr, name)

	logger.Info("Activity Response " + finalStr)

	return nil
}

func SayHello(ctx context.Context) (string, error) {
	return "Hello", nil
}

func SayWorld(ctx context.Context) (string, error) {
	return "World", nil
}
