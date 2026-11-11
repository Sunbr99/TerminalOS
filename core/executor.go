package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment struct {
	Vars       map[string]string
	CurrentDir string
}

func NewEnvironment() *Environment {
	wd, err := os.Getwd()
	if err != nil {
		wd = "/"
	}

	env := &Environment{Vars: map[string]string{}}
	env.CurrentDir = wd
	env.Set("HOME", os.Getenv("HOME"))
	env.Set("PATH", os.Getenv("PATH"))
	env.Set("PWD", wd)
	env.Set("SHELL", "/bin/bash")
	return env
}

func (e *Environment) Get(key string) string {
	if key == "PWD" {
		return e.CurrentDir
	}
	if value, ok := e.Vars[key]; ok {
		return value
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

func (e *Environment) Set(key, value string) {
	if key == "" {
		return
	}
	e.Vars[key] = value
	if key == "PWD" {
		e.CurrentDir = value
	}
}

func (e *Environment) Dir() string {
	return e.CurrentDir
}

func (e *Environment) Chdir(dir string) error {
	if dir == "" {
		dir = "."
	}
	cleaned := e.Expand(dir)
	if cleaned == "~" {
		cleaned = e.Get("HOME")
	}
	if cleaned == "" {
		cleaned = "."
	}

	if _, err := os.Stat(cleaned); err != nil {
		return fmt.Errorf("cd: %s: %w", cleaned, err)
	}

	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return err
	}

	e.CurrentDir = abs
	e.Set("PWD", abs)
	return nil
}

func (e *Environment) Expand(value string) string {
	if value == "" {
		return ""
	}

	if strings.HasPrefix(value, "~") {
		home := e.Get("HOME")
		if home != "" && (value == "~" || strings.HasPrefix(value, "~/")) {
			return strings.Replace(value, "~", home, 1)
		}
	}

	var builder strings.Builder
	for i := 0; i < len(value); i++ {
		if value[i] != '$' {
			builder.WriteByte(value[i])
			continue
		}

		if i+1 >= len(value) {
			builder.WriteByte(value[i])
			continue
		}

		if value[i+1] == '{' {
			end := strings.IndexByte(value[i+2:], '}')
			if end == -1 {
				builder.WriteString(value[i:])
				break
			}
			key := value[i+2 : i+2+end]
			builder.WriteString(e.Get(key))
			i = i + 2 + end
			continue
		}

		j := i + 1
		for j < len(value) && (value[j] == '_' || value[j] == '$' || (value[j] >= 'a' && value[j] <= 'z') || (value[j] >= 'A' && value[j] <= 'Z') || (value[j] >= '0' && value[j] <= '9')) {
			j++
		}
		if j == i+1 {
			builder.WriteByte(value[i])
			continue
		}
		key := value[i+1 : j]
		builder.WriteString(e.Get(key))
		i = j - 1
	}

	return builder.String()
}

func (e *Environment) VarsMap() map[string]string {
	return e.Vars
}

func (e *Environment) Remove(key string) {
	delete(e.Vars, key)
}
