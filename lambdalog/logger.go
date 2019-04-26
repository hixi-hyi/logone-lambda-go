package lambdalog

import "context"

type ManagerConfig struct {
	Type            string
	DefaultSeverity Severity
	OutputSeverity  Severity
}

type Manager struct {
	Config *ManagerConfig
}

func NewManager(mc *ManagerConfig) *Manager {
	return &Manager{mc}
}

func NewManagerDefault() *Manager {
	mc := &ManagerConfig{
		Type:            "request",
		DefaultSeverity: SeverityDebug,
		OutputSeverity:  SeverityDebug,
	}
	return NewManager(mc)
}

func (m *Manager) Recording(ctx context.Context) (*Recorder, func()) {
	r := NewRecorder(ctx, &RecorderConfig{m.Config})
	return r, func() {
		r.Finish()
	}
}
