package lambdalog

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

func TestLogContext(t *testing.T) {
	t.Run("NewLogContext", func(t *testing.T) {
		ctx := context.Background()
		lambdacontext.FunctionName = "dummyfunction"
		lambdacontext.FunctionVersion = "dummyversion"
		lambdacontext.MemoryLimitInMB = int(256)
		ctx = lambdacontext.NewContext(ctx, &lambdacontext.LambdaContext{
			AwsRequestID:       "dummyid",
			InvokedFunctionArn: "dummyarn",
		})
		lc := NewLogContext(ctx)

		assert.Exactly(t, lambdacontext.FunctionName, lc.FunctionName)
		assert.Exactly(t, lambdacontext.FunctionVersion, lc.FunctionVersion)
		assert.Exactly(t, lambdacontext.MemoryLimitInMB, lc.MemoryLimitInMega)
		assert.Exactly(t, "dummyid", lc.AwsRequestId)
		assert.Exactly(t, "dummyarn", lc.InvokedFunctionArn)
	})
}

func TestSeverityCount(t *testing.T) {
	t.Run("CountUp", func(t *testing.T) {
		sc := SeverityCount{}
		sc.CountUp(SeverityDebug)
		sc.CountUp(SeverityDebug)
		sc.CountUp(SeverityCritical)
		assert.Equal(t, int64(2), sc[SeverityDebug])
		assert.Equal(t, int64(1), sc[SeverityCritical])
		assert.Equal(t, SeverityCritical, sc.HighestSeverity())
	})
}
