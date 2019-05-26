package lambdalog

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
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

func TestLogger(t *testing.T) {
	t.Run("scenario", func(t *testing.T) {
		l := NewLoggerDefault(Context())
		finish := l.Start()
		attribute := map[string]string{"test": "1"}
		l.Debug("%s", "debug").WithTags("aws-sdk-error", "notify")
		l.Critical("%s", "critical").WithTags("aws-sdk-error", "error").WithAttributes(attribute)
		finish()

		assert.Exactly(t, int64(2), l.LogRequest.Runtime.Tags["aws-sdk-error"])
		assert.Exactly(t, int64(1), l.LogRequest.Runtime.Tags["notify"])
		assert.Exactly(t, int64(1), l.LogRequest.Runtime.Tags["error"])
		assert.Exactly(t, SeverityCritical, l.LogRequest.Runtime.Severity)
		// TODO test to output line
	})

}
