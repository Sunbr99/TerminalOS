package tests

import (
    "testing"

    "github.com/Sunbr99/TerminalOS/config"
    "github.com/Sunbr99/TerminalOS/core"
)

func TestParseCommandLine(t *testing.T) {
    commands, err := core.ParseCommandLine("echo hello | grep h > out.txt")
    if err != nil {
        t.Fatalf("parse failed: %v", err)
    }
    if len(commands) != 2 {
        t.Fatalf("expected 2 commands, got %d", len(commands))
    }
    if commands[0].Name != "echo" {
        t.Fatalf("first command name mismatch: %s", commands[0].Name)
    }
    if commands[1].Name != "grep" {
        t.Fatalf("second command name mismatch: %s", commands[1].Name)
    }
}

func TestEnvironmentExpansion(t *testing.T) {
    env := core.NewEnvironment()
    env.Set("NAME", "TerminalOS")
    value := env.Expand("hello $NAME")
    if value != "hello TerminalOS" {
        t.Fatalf("expected expansion, got %q", value)
    }

    cfg, err := config.Load("..")
    if err != nil {
        t.Fatalf("config load failed: %v", err)
    }
    if cfg.Security.Mode != "local_only" {
        t.Fatalf("unexpected config mode: %s", cfg.Security.Mode)
    }
}

func TestShellEngineInit(t *testing.T) {
    shell := core.NewShellEngine(config.DefaultConfig())
    if shell == nil {
        t.Fatal("shell engine was nil")
    }
    if _, ok := shell.Builtins["pwd"]; !ok {
        t.Fatal("pwd builtin missing")
    }
}
