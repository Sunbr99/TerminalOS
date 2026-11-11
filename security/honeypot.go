package security

import "strings"

type Firewall struct {
	DefaultPolicy string
	AllowedIPs    []string
}

func NewFirewall(defaultPolicy string, allowedIPs []string) *Firewall {
	return &Firewall{DefaultPolicy: defaultPolicy, AllowedIPs: allowedIPs}
}

func (f *Firewall) Allow(ip string) bool {
	if f == nil {
		return true
	}
	if f.DefaultPolicy == "allow" {
		return true
	}
	if len(f.AllowedIPs) == 0 {
		return false
	}
	for _, allowed := range f.AllowedIPs {
		if strings.TrimSpace(allowed) == ip {
			return true
		}
	}
	return false
}
