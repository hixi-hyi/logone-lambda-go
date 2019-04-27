package lambdalog

import (
	"context"
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

type LogEntry struct {
	Severity   string      `json:"severity"`
	Message    string      `json:"message"`
	Time       time.Time   `json:"time,omitempty"`
	Filename   string      `json:"filename"`
	Fileline   int         `json:"fileline"`
	Funcname   string      `json:"funcname"`
	Type       string      `json:"type,omitempty"`
	Elapsed    float64     `json:"elapsed,omitempty"`
	Attributes interface{} `json:"attributes,omitempty"`
}

func (lr *LogEntry) WithType(s string) *LogEntry {
	lr.Type = s
	return lr
}

func (lr *LogEntry) WithAttributes(i interface{}) *LogEntry {
	lr.Attributes = i
	return lr
}

type LogRequest struct {
	Type    string      `json:"type"`
	Context *LogContext `json:"context"`
	Runtime *LogRuntime `json:"runtime"`
	Config  *LogConfig  `json:"config"`
	//Extras  interface{}       `json:"extras,omitempty"`
}

type LogConfig struct {
	ElapsedUnit string `json:"elapsedUnit"`
}

type LogRuntime struct {
	Severity  string      `json:"severity"`
	StartTime time.Time   `json:"startTime"`
	EndTime   time.Time   `json:"endTime"`
	Elapsed   int64       `json:"elapsed"`
	Lines     []*LogEntry `json:"lines,omitempty"`
}

func (lr *LogRuntime) AppendLogEntry(l *LogEntry) {
	lr.Lines = append(lr.Lines, l)
}

type LogContext struct {
	FunctionName       string `json:"functionName"`
	FunctionVersion    string `json:"functionVersion"`
	MemoryLimitInMega  int    `json:"memoryLimitInMega"`
	InvokedFunctionArn string `json:"invokedFunctionArn"`
	AwsRequestId       string `json:"awsRequestId"`
	//Extras          interface{} `json:"extras,omitempty"`
}

func NewLogContext(ctx context.Context) *LogContext {
	lc := &LogContext{}

	lctx, _ := lambdacontext.FromContext(ctx)
	lc.FunctionName = lambdacontext.FunctionName
	lc.FunctionVersion = lambdacontext.FunctionVersion
	lc.MemoryLimitInMega = lambdacontext.MemoryLimitInMB
	lc.InvokedFunctionArn = lctx.InvokedFunctionArn
	lc.AwsRequestId = lctx.AwsRequestID

	return lc
}

type Severity int

const (
	SeverityUnknown Severity = iota
	SeverityDebug
	SeverityInfo
	SeverityWarning
	SeverityError
	SeverityCritical
)

var Severities = map[Severity]string{
	SeverityDebug:    "DEBUG",
	SeverityInfo:     "INFO",
	SeverityWarning:  "WARNING",
	SeverityError:    "ERROR",
	SeverityCritical: "CRITICAL",
}

func (s Severity) String() string {
	if v, ok := Severities[s]; ok {
		return v
	}
	return "UNKNOWN"
}

type SeverityCount map[Severity]int64

func (sc SeverityCount) CountUp(s Severity) {
	sc[s]++
}
func (sc SeverityCount) HighestSeverity() Severity {
	keys := []int{int(SeverityUnknown)}
	for k := range sc {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return Severity(keys[len(keys)-1])
}
