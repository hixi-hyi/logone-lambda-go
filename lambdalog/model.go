package lambdalog

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

type LogEntry struct {
	Severity   Severity    `json:"severity"`
	Message    string      `json:"message"`
	Time       time.Time   `json:"time,omitempty"`
	Filename   string      `json:"filename"`
	Fileline   int         `json:"fileline"`
	Funcname   string      `json:"funcname"`
	Tags       []string    `json:"tags,omitempty"`
	Elapsed    float64     `json:"elapsed,omitempty"`
	Attributes interface{} `json:"attributes,omitempty"`
}

func (lr *LogEntry) WithTags(tags ...string) *LogEntry {
	lr.Tags = tags
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
	Severity   Severity      `json:"severity"`
	StartTime  time.Time     `json:"startTime"`
	EndTime    time.Time     `json:"endTime"`
	Elapsed    int64         `json:"elapsed"`
	Lines      []*LogEntry   `json:"lines,omitempty"`
	Tags       LogTags       `json:"tags,omitempty"`
	Severities SeverityCount `json:"-"`
}

func NewLogRuntime() *LogRuntime {
	lr := &LogRuntime{}
	lr.Tags = LogTags{}
	lr.Severities = SeverityCount{}
	return lr
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

type LogTags map[string]int64

func (lt LogTags) CountUp(tags ...string) {
	for _, t := range tags {
		lt[t] += 1
	}
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

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
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
