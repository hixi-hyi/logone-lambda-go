package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hixi-hyi/logone-lambda-go/lambdalog"
)

func main() {
	lambda.Start(handler)
}

type Event struct {
	Event string `json:"event"`
}

func handler(ctx context.Context, request Event) (string, error) {
	manager := lambdalog.NewManagerDefault()
	logger, finish := manager.Recording(ctx)
	defer finish()
	ctx = lambdalog.NewContextWithLogger(ctx, logger)
	logger.Debug("lambda invoked").WithTags("request").WithAttributes(request)

	f := &Function{}
	f.Func(ctx)

	return "success", nil
}
