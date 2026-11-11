package library

import (
	"fmt"
)

type Plugin interface {
	Name() string
	Run(args []string) error
}

type Registry struct {
	plugins map[string]Plugin
}

func NewRegistry() *Registry {
	return &Registry{plugins: map[string]Plugin{}}
}

func (r *Registry) Register(plugin Plugin) {
	if r.plugins == nil {
		r.plugins = map[string]Plugin{}
	}
	r.plugins[plugin.Name()] = plugin
}

func (r *Registry) List() []string {
	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}
	return names
}

func (r *Registry) Run(name string, args []string) error {
	plugin, ok := r.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %q not found", name)
	}
	return plugin.Run(args)
}
