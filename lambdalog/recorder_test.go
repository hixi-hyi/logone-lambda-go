package lambdalog

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func Context() context.Context {
	ctx := context.Background()
	lambdacontext.FunctionName = "dummyfunction"
	lambdacontext.FunctionVersion = "dummyversion"
	lambdacontext.MemoryLimitInMB = int(256)
	lc := &lambdacontext.LambdaContext{
		AwsRequestID:       "dummyid",
		InvokedFunctionArn: "dummyarn",
	}
	ctx = lambdacontext.NewContext(ctx, lc)
	return ctx
}

func TestRecorder(t *testing.T) {
	t.Run("scenario", func(t *testing.T) {
		r := NewRecorderDefault(Context())
		finish := r.Start()
		r.Debug("%s", "a")
		finish()
	})

}
