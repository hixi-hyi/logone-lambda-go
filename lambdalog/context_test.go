package lambdalog

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	t.Run("scenario", func(t *testing.T) {
		ctx := context.Background()
		lambdacontext.FunctionName = "dummyfunction"
		lambdacontext.FunctionVersion = "dummyversion"
		lambdacontext.MemoryLimitInMB = int(256)
		lc := &lambdacontext.LambdaContext{
			AwsRequestID:       "dummyid",
			InvokedFunctionArn: "dummyarn",
		}
		ctx = lambdacontext.NewContext(ctx, lc)
		r := NewLoggerDefault(Context())
		ctx = NewContextWithLogger(ctx, r)
		nr, ok := LoggerFromContext(ctx)
		assert.Exactly(t, true, ok)
		assert.Exactly(t, r, nr)
	})
}
