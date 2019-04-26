package lambdalog

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Recorder struct {
	Config        *Config
	LogRequest    *LogRequest
	SeverityCount *SeverityCount
}

func NewRecorder(ctx context.Context, c *Config) *Recorder {
	lr := &LogRequest{
		Type:    c.Type,
		Context: NewLogContext(ctx),
		Runtime: &LogRuntime{},
	}

	r := &Recorder{
		Config:        c,
		LogRequest:    lr,
		SeverityCount: &SeverityCount{},
	}
	return r
}

func NewRecorderDefault(ctx context.Context) *Recorder {
	c := &Config{
		Type:            "request",
		DefaultSeverity: SeverityDebug,
		OutputSeverity:  SeverityDebug,
	}
	return NewRecorder(ctx, c)
}

func (r *Recorder) Start() func() {
	r.LogRequest.Runtime.Start()
	return func() {
		r.Finish()
	}
}
func (r *Recorder) Finish() {
	r.LogRequest.Runtime.Finish(r.SeverityCount.HighestSeverity())
	logline, _ := json.Marshal(r.LogRequest)
	fmt.Println(string(logline))
}

func (r *Recorder) Record(severity Severity, message string) *LogEntry {
	funcname, filename, fileline := FileInfo(3)
	e := &LogEntry{
		Severity: severity.String(),
		Message:  message,
		Time:     time.Now(),
		Filename: filename,
		Fileline: fileline,
		Funcname: funcname,
	}
	r.LogRequest.Runtime.AppendLogEntry(e)
	r.SeverityCount.CountUp(severity)
	return e
}

func (r *Recorder) Debug(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityDebug, fmt.Sprintf(f, v...))
}
func (r *Recorder) Info(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityInfo, fmt.Sprintf(f, v...))
}
func (r *Recorder) Warning(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityWarning, fmt.Sprintf(f, v...))
}
func (r *Recorder) Error(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityError, fmt.Sprintf(f, v...))
}
func (r *Recorder) Critical(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityCritical, fmt.Sprintf(f, v...))
}

func FileInfo(depth int) (string, string, int) {
	pc, filename, fileline, ok := runtime.Caller(depth)
	if !ok {
		return "???", "???", 0
	}
	return runtime.FuncForPC(pc).Name(), filename, fileline
}
