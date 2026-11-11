package security

type Sandbox struct {
	Enabled bool
	Mode    string
}

func NewSandbox(enabled bool, mode string) *Sandbox {
	return &Sandbox{Enabled: enabled, Mode: mode}
}

func (s *Sandbox) Allowed() bool {
	return s != nil && s.Enabled
}

func (s *Sandbox) ModeName() string {
	if s == nil || s.Mode == "" {
		return "local_only"
	}
	return s.Mode
}
