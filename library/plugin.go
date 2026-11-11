package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func Echo(args []string, env EnvStore) error {
	expanded := make([]string, 0, len(args))
	for _, arg := range args {
		expanded = append(expanded, env.Expand(arg))
	}
	fmt.Println(strings.Join(expanded, " "))
	return nil
}

func Cat(args []string, env EnvStore) error {
	if len(args) == 0 {
		_, err := io.Copy(os.Stdout, os.Stdin)
		return err
	}

	for _, path := range args {
		filename := env.Expand(path)
		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	}
	return nil
}

func Grep(args []string, env EnvStore) error {
	if len(args) < 2 {
		return fmt.Errorf("grep requires a pattern and file(s)")
	}

	pattern := env.Expand(args[0])
	for _, filepath := range args[1:] {
		file := env.Expand(filepath)
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, pattern) {
				fmt.Println(line)
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

func Wc(args []string, env EnvStore) error {
	if len(args) == 0 {
		return fmt.Errorf("wc requires at least one file")
	}

	for _, filepath := range args {
		file := env.Expand(filepath)
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		text := string(data)
		newlineCount := strings.Count(text, "\n")
		wordCount := len(strings.Fields(text))
		byteCount := len(data)
		fmt.Printf("%d %d %d %s\n", newlineCount, wordCount, byteCount, file)
	}
	return nil
}

func Env(args []string, env EnvStore) error {
	items := env.VarsMap()
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s=%s\n", key, items[key])
	}
	return nil
}

func Export(args []string, env EnvStore) error {
	if len(args) == 0 {
		return nil
	}
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			env.Set(parts[0], env.Expand(parts[1]))
			continue
		}
		if value := env.Get(arg); value != "" {
			env.Set(arg, value)
		}
	}
	return nil
}

func Unset(args []string, env EnvStore) error {
	for _, arg := range args {
		env.Remove(arg)
	}
	return nil
}
