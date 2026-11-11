package core

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

func openRedirection(target string, appendMode bool) (*os.File, error) {
    flags := os.O_CREATE | os.O_WRONLY
    if appendMode {
        flags |= os.O_APPEND
    } else {
        flags |= os.O_TRUNC
    }

    dir := filepath.Dir(target)
    if dir != "" && dir != "." {
        if err := os.MkdirAll(dir, 0o755); err != nil {
            return nil, err
        }
    }

    return os.OpenFile(target, flags, 0o644)
}

func (s *ShellEngine) runExternalCommand(cmdstr string) error {
    if cmdstr == "" {
        return nil
    }

    commands := []string{"sh", "-c", cmdstr}
    execCmd := exec.Command(commands[0], commands[1:]...)
    execCmd.Dir = s.Env.CurrentDir
    execCmd.Env = append(os.Environ(), s.envList()...)
    execCmd.Stdout = os.Stdout
    execCmd.Stderr = os.Stderr
    return execCmd.Run()
}

func (s *ShellEngine) runPipelineCommand(pipeline string) error {
    if pipeline == "" {
        return nil
    }

    execCmd := exec.Command("sh", "-c", pipeline)
    execCmd.Dir = s.Env.CurrentDir
    execCmd.Env = append(os.Environ(), s.envList()...)
    execCmd.Stdout = os.Stdout
    execCmd.Stderr = os.Stderr
    if err := execCmd.Run(); err != nil {
        return fmt.Errorf("pipeline execution failed: %w", err)
    }
    return nil
}
