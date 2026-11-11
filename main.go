package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/Sunbr99/TerminalOS/config"
    "github.com/Sunbr99/TerminalOS/core"
)

func main() {
    cfg, err := config.Load(".")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
        os.Exit(1)
    }

    shell := core.NewShellEngine(cfg)
    fmt.Println("TerminalOS initialized.")
    fmt.Println("Type 'help' for commands or 'exit' to quit.")

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("terminalos> ")
        input, err := reader.ReadString('\n')
        if err == io.EOF {
            fmt.Println()
            return
        }
        if err != nil {
            fmt.Fprintf(os.Stderr, "read error: %v\n", err)
            return
        }

        input = strings.TrimSpace(input)
        if input == "" {
            continue
        }

        if err := shell.Execute(input); err != nil {
            fmt.Fprintf(os.Stderr, "terminalos: %v\n", err)
        }
    }
}
