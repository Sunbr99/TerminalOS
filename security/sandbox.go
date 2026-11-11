package security

import (
    "fmt"
    "os"
    "time"
)

type Logger struct {
    Enabled bool
}

func NewLogger(enabled bool) *Logger {
    return &Logger{Enabled: enabled}
}

func (l *Logger) Log(message string) {
    if l == nil || !l.Enabled {
        return
    }
    now := time.Now().UTC().Format(time.RFC3339)
    fmt.Fprintf(os.Stdout, "[%s] %s\n", now, message)
}
