# logone-lambda-go
The library is support to structured logging using JSON format for aws-lambda.
CloudWatch Logs is support to JSON format, you will be easier to investigate a log.

## Caution
It's experimental version.
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
    logger, finish := manager.Recording(ctx)
    defer finish()
    log.Info("%s", "test")
    return nil, nil
}
```
If there is a possibility that `panic()` will occur, call `finish()` using the `recover ()` statement.

### Using nested function
```
func handler(ctx context.Context, req events.CloudWatchEvent) (interface{}, error) {
    logger, finish := manager.Recording(ctx)
    defer finish()
    ctx = lambdalog.NewContextWithLogger(ctx, logger)
    func2(ctx)
    return nil, nil
}
func func2(ctx context.Context) {
    logger, _ := lambdalog.LoggerFromContext(ctx)
    logger.Info("nested function").WithTags("call-function")
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

### With Tags
You can add annotations to investigate the log.
If you want to investigate the tag, you can get logs using `{ $.runtime.tags.aws-sdk-error >= 1 }`
```
res, err := sns.Publish(params)
if err != nil {
    log.Error("error occured in sns.Publish: %s", err).WithTags("aws-sdk-error")
    return err
}
```
### With Attributes
You can add an annotation to the log to get more details.
The `Attributes` value must be defined as a type that can be JSON.Marshal. If you want to output value that cannot be JSON.Marshal, you use fmt.Sprintf or primitive type. (e.g. error is not struct or primitive type, you must use fmt.Sprintf("%s", err) or err.Error())
```
log.Info("publish successfully").WithArrtibutes(res)
```
### With Error
You can add an annotation to the log to get more details.
```
if err != nil {
    log.Critical("error occured").WithError(err)
}
```

## Outputs
Below is an example of output from [./example/lambda].
```
{
    "type": "request",
    "context": {
        "functionName": "example-logone-lambda-go",
        "functionVersion": "$LATEST",
        "memoryLimitInMega": 128,
        "invokedFunctionArn": "arn:aws:lambda:ap-northeast-1:123456789012:function:example-logone-lambda-go",
        "awsRequestId": "9871c5f2-3920-45ef-a54d-5c8a32caf5fc"
    },
    "runtime": {
        "severity": "INFO",
        "startTime": "2019-05-26T02:59:02.168056547Z",
        "endTime": "2019-05-26T02:59:02.168068974Z",
        "elapsed": 0,
        "lines": [
            {
                "severity": "DEBUG",
                "message": "lambda invoked",
                "time": "2019-05-26T02:59:02.168065633Z",
                "fileline": 23,
                "funcname": "main.handler",
                "tags": [
                    "request"
                ],
                "attributes": {
                    "event": "example"
                }
            },
            {
                "severity": "INFO",
                "message": "nested function",
                "time": "2019-05-26T02:59:02.168068487Z",
                "fileline": 13,
                "funcname": "main.(*Function).Func",
                "tags": [
                    "call-function"
                ]
            }
        ],
        "tags": {
            "call-function": 1,
            "request": 1
        }
    },
    "config": {
        "elapsedUnit": "1ms"
    }
}
```

## ToDo
* godoc
* Support to OutputFunc in lambdalog.Config. It is fmt.Println() now.
* Support to OutputSeverity in lambdalog.Config. It is print all logs now.
* Support to DefaultSeverity in lambdalog.Config. It is "UNKNOWN" now, if you don't write a any log.
* Support to OutputColumns in lambdalog.Config. It is `RFILELINE | RFUNCNAME | CELAPSED_UNIT` now.
* Do not want to consider about Attributes limitation. Attributes is support to only type can be JSON.Marshal or primitive now.
