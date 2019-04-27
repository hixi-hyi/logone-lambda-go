# logone-lambda-go
The library is support to structured logging using JSON format for aws-lambda.
CloudWatch Logs is also support to JSON format, so you will be easier to investigate a log.

## Caution
Log messages are temporarily stored in memory.
Be careful in below point on a huge system.
* Memory usage.
* The log is not written out until the function ends.

## Usage
### Simple Example
```
import github.com/hixi-hyi/logone-lambda-go/lambdalog

var manager *lambdalog.Manager
func init() {
    manager = lambdalog.NewManagerDefault()
}
func main() {
    lambda.Start(handler)
}

func handler(ctx context.Context, req events.CloudWatchEvent) (interface{}, error) {
    log, flush := manager.Recording(ctx)
    defer flush()
    log.Info("%s", "test")
    return nil, nil
}
// {"type":"request","context":{"functionName":"test","functionVersion":"$LATEST","memoryLimitInMega":128,"invokedFunctionArn":"arn:aws:lambda:ap-northeast-1:123456789:function:test","awsRequestId":"b609c7aa-8b33-181f-70a8-b0d4d44b31f1"},"runtime":{"severity":"INFO","startTime":"2019-04-27T20:12:20.5258606Z","endTime":"2019-04-27T20:12:20.5370457Z","elapsed":11,"lines":[{"severity":"INFO","message":"test","time":"2019-04-27T20:12:20.5370358Z","filename":"/yourcode/main.go","fileline":80,"funcname":"main.handler"}]},"config":{"elapsedUnit":"1ms"}}
```

### Using nested function
```
func handler(ctx context.Context, req events.CloudWatchEvent) (interface{}, error) {
    log, flush := manager.Recording(ctx)
    defer flush()
    ctx = lambdalog.NewContextWithLogger(ctx)
    func2(ctx)
    return nil, nil
}
func func2(ctx context.Context) {
    log := lambdalog.LoggerFromContext(ctx)
}
```

### With Config
You can change the some parameters.
```
manager = lambdalog.NewManager(&lambdalog.Config{
	Type:            "request",
	DefaultSeverity: lambdalog.SeverityDebug,
	ElapsedUnit:     time.Millisecond,
	JsonIndent:      false
})
```

### With Attributes And Type
You can add the anotation to log.
The `Attributes` value must be defined as a type that can be JSON.Marshal. If you want to output value that cannot be JSON.Marshal, you use fmt.Sprintf or primitive type. (e.g. error is not struct or primitive type, you must use fmt.Sprintf("%s", err) or err.Error())
```
res, err := sns.Publish(params)
if err != nil {
    log.Error("error occured in sns.Publish: %s", err).WithType("aws-sdk-error")
    // {"severity":"ERROR","message":"error occured in sns.Publish: xxxx","time":"2019-04-27T20:12:20.5370358Z","filename":"/yourcode/main.go","fileline":80,"funcname":"main.handler"}
    return err
}
log.Info("publish successfully).WithArrtibutes(res)
// {"severity":"INFO","message":"publish successfully","time":"2019-04-27T20:28:05.8390507Z","filename":"/yourcode/main.go", "fileline":81,"funcname":"main.handler","attributes":{"MessageId":"xxxxxx"}}

```

## ToDo
* godoc
* Support to OutputFunc in lambdalog.Config. It is fmt.Println() now.
* Support to OutputSeverity in lambdalog.Config. It is print all logs now.
* Support to DefaultSeverity in lambdalog.Config. It is "UNKNOWN" now, if you don't write a any log.
* Support to OutputColumns in lambdalog.Config. Now is `RFILENAME | RFILELINE | RFUNCNAME | CELAPSED_UNIT` now.
* Do not want to consider about Attributes limitation. Now Attributes is support to only type can be JSON.Marshal or primitive.
