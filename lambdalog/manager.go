package lambdalog

import (
	"context"
	"time"
)

type Config struct {
	Type            string
	DefaultSeverity Severity
	OutputSeverity  Severity
	JsonIndent      bool
	ElapsedUnit     time.Duration
}

type Manager struct {
	Config *Config
}

func NewConfigDefault() *Config {
	return &Config{
		Type:            "request",
		DefaultSeverity: SeverityDebug,
		OutputSeverity:  SeverityDebug,
		ElapsedUnit:     time.Millisecond,
	}
}

func NewManager(mc *Config) *Manager {
	return &Manager{mc}
}

func NewManagerDefault() *Manager {
	return NewManager(NewConfigDefault())
}

func (m *Manager) Recording(ctx context.Context) (*Logger, func()) {
	r := NewLogger(ctx, m.Config)
	return r, func() {
		r.Finish()
	}
}

func (m *Manager) RecordingInContext(ctx context.Context) (context.Context, func()) {
	l, write := m.Recording(ctx)
	nctx := NewContextWithLogger(ctx, l)
	return nctx, write
}
