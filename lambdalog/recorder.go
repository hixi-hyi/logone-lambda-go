package lambdalog

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Logger struct {
	Config        *Config
	LogRequest    *LogRequest
	SeverityCount *SeverityCount
}

func NewLogger(ctx context.Context, c *Config) *Logger {
	lr := &LogRequest{
		Type:    c.Type,
		Context: NewLogContext(ctx),
		Runtime: &LogRuntime{},
		Config: &LogConfig{
			ElapsedUnit: c.ElapsedUnit.String(),
		},
	}

	r := &Logger{
		Config:        c,
		LogRequest:    lr,
		SeverityCount: &SeverityCount{},
	}
	return r
}

func NewLoggerDefault(ctx context.Context) *Logger {
	c := NewConfigDefault()
	return NewLogger(ctx, c)
}

func (r *Logger) Start() func() {
	lr := r.LogRequest.Runtime
	lr.StartTime = time.Now()
	return func() {
		r.Finish()
	}
}
func (r *Logger) Finish() {
	lr := r.LogRequest.Runtime
	lr.EndTime = time.Now()
	elapsed := lr.EndTime.Sub(lr.StartTime)
	lr.Elapsed = int64(elapsed / r.Config.ElapsedUnit)
	lr.Severity = r.SeverityCount.HighestSeverity().String()
	var logline []byte
	if r.Config.JsonIndent {
		logline, _ = json.MarshalIndent(r.LogRequest, "", "  ")
	} else {
		logline, _ = json.Marshal(r.LogRequest)
	}
	fmt.Println(string(logline))
}

func (r *Logger) Record(severity Severity, message string) *LogEntry {
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

func (r *Logger) Debug(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityDebug, fmt.Sprintf(f, v...))
}
func (r *Logger) Info(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityInfo, fmt.Sprintf(f, v...))
}
func (r *Logger) Warning(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityWarning, fmt.Sprintf(f, v...))
}
func (r *Logger) Error(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityError, fmt.Sprintf(f, v...))
}
func (r *Logger) Critical(f string, v ...interface{}) *LogEntry {
	return r.Record(SeverityCritical, fmt.Sprintf(f, v...))
}

func FileInfo(depth int) (string, string, int) {
	pc, filename, fileline, ok := runtime.Caller(depth)
	if !ok {
		return "???", "???", 0
	}
	return runtime.FuncForPC(pc).Name(), filename, fileline
}
