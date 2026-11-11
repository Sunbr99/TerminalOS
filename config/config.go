package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FirewallConfig struct {
	Enabled       bool     `json:"enabled"`
	DefaultPolicy string   `json:"default_policy"`
	AllowedIPs    []string `json:"allowed_ips"`
}

type HoneypotConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

type LoggerConfig struct {
	Enabled bool `json:"enabled"`
}

type SecurityConfig struct {
	Mode     string         `json:"mode"`
	Firewall FirewallConfig `json:"firewall"`
	Honeypot HoneypotConfig `json:"honeypot"`
	Logger   LoggerConfig   `json:"logger"`
}

type Config struct {
	Security SecurityConfig `json:"security"`
}

func DefaultConfig() Config {
	return Config{
		Security: SecurityConfig{
			Mode: "local_only",
			Firewall: FirewallConfig{
				Enabled:       true,
				DefaultPolicy: "deny",
				AllowedIPs:    []string{"127.0.0.1", "::1"},
			},
			Honeypot: HoneypotConfig{
				Enabled: true,
				Port:    2222,
			},
			Logger: LoggerConfig{Enabled: true},
		},
	}
}

func Load(basePath string) (Config, error) {
	configPath := filepath.Join(basePath, "config", "default_config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	if cfg.Security.Mode == "" {
		cfg.Security.Mode = "local_only"
	}
	if !cfg.Security.Firewall.Enabled {
		cfg.Security.Firewall.Enabled = true
	}
	if cfg.Security.Firewall.DefaultPolicy == "" {
		cfg.Security.Firewall.DefaultPolicy = "deny"
	}
	if cfg.Security.Honeypot.Port == 0 {
		cfg.Security.Honeypot.Port = 2222
	}

	return cfg, nil
}
