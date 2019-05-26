package main

import (
	"context"

	"github.com/hixi-hyi/logone-lambda-go/lambdalog"
)

type Function struct{}

func (f *Function) Func(ctx context.Context) {
	logger, _ := lambdalog.LoggerFromContext(ctx)
	logger.Info("nested function").WithTags("call-function")
}
