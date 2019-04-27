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

	l := &Logger{
		Config:        c,
		LogRequest:    lr,
		SeverityCount: &SeverityCount{},
	}
	return l
}

func NewLoggerDefault(ctx context.Context) *Logger {
	c := NewConfigDefault()
	return NewLogger(ctx, c)
}

func (l *Logger) Start() func() {
	lr := l.LogRequest.Runtime
	lr.StartTime = time.Now()
	return func() {
		l.Finish()
	}
}
func (l *Logger) Finish() {
	lr := l.LogRequest.Runtime
	lr.EndTime = time.Now()
	elapsed := lr.EndTime.Sub(lr.StartTime)
	lr.Elapsed = int64(elapsed / l.Config.ElapsedUnit)
	lr.Severity = l.SeverityCount.HighestSeverity().String()
	var logline []byte
	if l.Config.JsonIndent {
		logline, _ = json.MarshalIndent(l.LogRequest, "", "  ")
	} else {
		logline, _ = json.Marshal(l.LogRequest)
	}
	fmt.Println(string(logline))
}

func (l *Logger) Record(severity Severity, message string) *LogEntry {
	funcname, filename, fileline := FileInfo(3)
	e := &LogEntry{
		Severity: severity.String(),
		Message:  message,
		Time:     time.Now(),
		Filename: filename,
		Fileline: fileline,
		Funcname: funcname,
	}
	l.LogRequest.Runtime.AppendLogEntry(e)
	l.SeverityCount.CountUp(severity)
	return e
}

func (l *Logger) Debug(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityDebug, fmt.Sprintf(f, v...))
}
func (l *Logger) Info(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityInfo, fmt.Sprintf(f, v...))
}
func (l *Logger) Warning(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityWarning, fmt.Sprintf(f, v...))
}
func (l *Logger) Error(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityError, fmt.Sprintf(f, v...))
}
func (l *Logger) Critical(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityCritical, fmt.Sprintf(f, v...))
}

func FileInfo(depth int) (string, string, int) {
	pc, filename, fileline, ok := runtime.Caller(depth)
	if !ok {
		return "???", "???", 0
	}
	return runtime.FuncForPC(pc).Name(), filename, fileline
}
