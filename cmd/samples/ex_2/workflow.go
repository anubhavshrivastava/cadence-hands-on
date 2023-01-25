package ex_2

import (
	"context"
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

func Workflow(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx, opts())
	logger := workflow.GetLogger(ctx)
	logger.Info("Test Workflow  workflow started")

	// waiting for signal from the user
	var signalValue string
	s := workflow.NewSelector(ctx)
	s.AddReceive(workflow.GetSignalChannel(ctx, "testSignalChannel"), func(c workflow.Channel, ok bool) {
		if ok {
			c.Receive(ctx, &signalValue)
			logger.Info("Received signal on testSignalChannel.", zap.String("Signal Value", signalValue))
		}
	})
	s.Select(ctx)

	logger.Info("Executing Activity")
	var resp string
	err := workflow.ExecuteActivity(ctx, SayHelloToSignalValue).Get(ctx, &resp)
	if err != nil {
		logger.Error("Activity  Failed.", zap.Error(err))
		return err
	}

	logger.Info("Activity Response " + resp)

	return nil
}

func SayHelloToSignalValue(ctx context.Context, signalValue string) (string, error) {
	return "Hello - Signal - " + signalValue, nil
}
