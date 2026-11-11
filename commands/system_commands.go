package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

func Pwd(args []string, env EnvStore) error {
    fmt.Println(env.Dir())
    return nil
}

func Cd(args []string, env EnvStore) error {
    target := "."
    if len(args) > 0 {
        target = args[0]
    }
    target = env.Expand(target)
    if target == "~" || strings.HasPrefix(target, "~/") {
        home, err := os.UserHomeDir()
        if err != nil {
            return err
        }
        if target == "~" {
            target = home
        } else {
            target = filepath.Join(home, strings.TrimPrefix(target, "~/"))
        }
    }
    return env.Chdir(target)
}

func Ls(args []string, env EnvStore) error {
    target := env.Dir()
    if len(args) > 0 {
        target = env.Expand(args[0])
    }

    entries, err := os.ReadDir(target)
    if err != nil {
        return err
    }

    names := make([]string, 0, len(entries))
    for _, entry := range entries {
        names = append(names, entry.Name())
    }
    sort.Strings(names)
    for _, name := range names {
        fmt.Println(name)
    }
    return nil
}
