package security

type Honeypot struct {
	Enabled bool
	Port    int
}

func NewHoneypot(enabled bool, port int) *Honeypot {
	return &Honeypot{Enabled: enabled, Port: port}
}

func (h *Honeypot) IsEnabled() bool {
	return h != nil && h.Enabled
}
