package library

import (
	"fmt"
	"strings"
)

type HelloPlugin struct{}

func (p *HelloPlugin) Name() string { return "hello" }

func (p *HelloPlugin) Run(args []string) error {
	if len(args) == 0 {
		fmt.Println("Hello from TerminalOS")
		return nil
	}
	fmt.Printf("Hello %s\n", strings.Join(args, " "))
	return nil
}
