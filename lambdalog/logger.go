package lambdalog

import "context"

type Config struct {
	Type            string
	DefaultSeverity Severity
	OutputSeverity  Severity
	JsonIndent      bool
}

type Manager struct {
	Config *Config
}

func NewManager(mc *Config) *Manager {
	return &Manager{mc}
}

func NewManagerDefault() *Manager {
	mc := &Config{
		Type:            "request",
		DefaultSeverity: SeverityDebug,
		OutputSeverity:  SeverityDebug,
	}
	return NewManager(mc)
}

func (m *Manager) Recording(ctx context.Context) (*Recorder, func()) {
	r := NewRecorder(ctx, m.Config)
	return r, func() {
		r.Finish()
	}
}
