package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Sunbr99/TerminalOS/commands"
	"github.com/Sunbr99/TerminalOS/config"
	"github.com/Sunbr99/TerminalOS/library"
)

type ShellEngine struct {
	Config   config.Config
	Env      *Environment
	Builtins map[string]BuiltinFunc
	History  []string
	Registry *library.Registry
}

type BuiltinFunc func(args []string, env *Environment, shell *ShellEngine) error

func NewShellEngine(cfg config.Config) *ShellEngine {
	env := NewEnvironment()
	shell := &ShellEngine{
		Config:   cfg,
		Env:      env,
		Builtins: map[string]BuiltinFunc{},
		Registry: library.NewRegistry(),
	}

	shell.RegisterBuiltin("pwd", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Pwd(args, env)
	})
	shell.RegisterBuiltin("cd", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Cd(args, env)
	})
	shell.RegisterBuiltin("ls", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Ls(args, env)
	})
	shell.RegisterBuiltin("echo", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Echo(args, env)
	})
	shell.RegisterBuiltin("cat", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Cat(args, env)
	})
	shell.RegisterBuiltin("grep", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Grep(args, env)
	})
	shell.RegisterBuiltin("wc", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Wc(args, env)
	})
	shell.RegisterBuiltin("env", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Env(args, env)
	})
	shell.RegisterBuiltin("export", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Export(args, env)
	})
	shell.RegisterBuiltin("unset", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Unset(args, env)
	})
	shell.RegisterBuiltin("help", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Help(args, env)
	})
	shell.RegisterBuiltin("exit", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.Exit(args, env)
	})
	shell.RegisterBuiltin("history", func(args []string, env *Environment, shell *ShellEngine) error {
		return commands.History(args, shell)
	})

	shell.Registry.Register(&library.HelloPlugin{})
	return shell
}

func (s *ShellEngine) RegisterBuiltin(name string, fn BuiltinFunc) {
	s.Builtins[name] = fn
}

func (s *ShellEngine) Execute(line string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	s.History = append(s.History, line)
	parsed, err := ParseCommandLine(line)
	if err != nil {
		return err
	}
	if len(parsed) == 0 {
		return nil
	}

	if len(parsed) == 1 {
		cmd := parsed[0]
		if fn, ok := s.Builtins[cmd.Name]; ok {
			err = fn(cmd.Args, s.Env, s)
			if errors.Is(err, commands.ErrExit) {
				fmt.Println("Goodbye.")
				return err
			}
			return err
		}
		return s.runExternal(cmd)
	}

	return s.runPipeline(parsed)
}

func (s *ShellEngine) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("terminalos> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println()
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if err := s.Execute(line); err != nil {
			if errors.Is(err, commands.ErrExit) {
				return
			}
			fmt.Fprintf(os.Stderr, "terminalos: %v\n", err)
		}
	}
}

func (s *ShellEngine) runExternal(cmd Command) error {
	if cmd.Name == "" {
		return nil
	}

	argv := make([]string, 0, len(cmd.Args)+1)
	argv = append(argv, cmd.Name)
	argv = append(argv, cmd.Args...)

	execCmd := execCommand(cmd.Name, cmd.Args...)
	execCmd.Dir = s.Env.CurrentDir
	execCmd.Env = append(os.Environ(), s.envList()...)

	if cmd.Output != "" {
		target := s.Env.Expand(cmd.Output)
		file, err := openRedirection(target, cmd.Append)
		if err != nil {
			return err
		}
		defer file.Close()
		execCmd.Stdout = file
	} else {
		execCmd.Stdout = os.Stdout
	}

	if cmd.Stderr != "" {
		target := s.Env.Expand(cmd.Stderr)
		file, err := openRedirection(target, cmd.StderrAppend)
		if err != nil {
			return err
		}
		defer file.Close()
		execCmd.Stderr = file
	} else {
		execCmd.Stderr = os.Stderr
	}

	return execCmd.Run()
}

func (s *ShellEngine) runPipeline(commands []Command) error {
	segments := make([]string, 0, len(commands))
	for _, cmd := range commands {
		if _, ok := s.Builtins[cmd.Name]; ok {
			return fmt.Errorf("builtin %q cannot be used in a pipeline", cmd.Name)
		}
		segments = append(segments, cmd.String())
	}

	pipeline := strings.Join(segments, " | ")
	execCmd := execCommand("sh", "-c", pipeline)
	execCmd.Dir = s.Env.CurrentDir
	execCmd.Env = append(os.Environ(), s.envList()...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}

func (s *ShellEngine) envList() []string {
	vars := make([]string, 0, len(s.Env.Vars))
	keys := make([]string, 0, len(s.Env.Vars))
	for k := range s.Env.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vars = append(vars, k+"="+s.Env.Vars[k])
	}
	return vars
}

func execCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
