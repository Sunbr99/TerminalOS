package commands

import (
	"errors"
	"fmt"
	"strings"
)

var ErrExit = errors.New("exit")

type EnvStore interface {
	Get(string) string
	Set(string, string)
	Dir() string
	Chdir(string) error
	Expand(string) string
	Remove(string)
	VarsMap() map[string]string
}

func Help(args []string, env EnvStore) error {
	commands := []string{"pwd", "cd", "ls", "echo", "cat", "grep", "wc", "env", "export", "unset", "help", "exit", "history"}
	fmt.Println("Available TerminalOS commands:")
	for _, name := range commands {
		fmt.Printf("  %s\n", name)
	}
	return nil
}

func Exit(args []string, env EnvStore) error {
	return ErrExit
}

func History(args []string, history []string) error {
	for _, entry := range history {
		fmt.Println(entry)
	}
	return nil
}

func _printCommandList() {
	fmt.Println(strings.Join([]string{"pwd", "cd", "ls", "echo", "cat", "grep", "wc", "env", "export", "unset", "help", "exit"}, ", "))
}
