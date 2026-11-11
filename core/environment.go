package core

import (
	"fmt"
	"strings"
	"unicode"
)

type Command struct {
	Name         string
	Args         []string
	Output       string
	Append       bool
	Stderr       string
	StderrAppend bool
}

func (c Command) String() string {
	parts := make([]string, 0, len(c.Args)+1)
	parts = append(parts, shellQuote(c.Name))
	for _, arg := range c.Args {
		parts = append(parts, shellQuote(arg))
	}
	cmd := strings.Join(parts, " ")
	if c.Output != "" {
		cmd += " > " + shellQuote(c.Output)
	}
	if c.Stderr != "" {
		cmd += " 2> " + shellQuote(c.Stderr)
	}
	return cmd
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	if !strings.ContainsAny(value, " \t\r\n\"'|><$&;(){}[]*?!~`) ") {
		return value
	}
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}

func ParseCommandLine(line string) ([]Command, error) {
	tokens := tokenize(line)
	if len(tokens) == 0 {
		return nil, nil
	}

	commands := make([]Command, 0, 4)
	current := Command{}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case "|":
			if current.Name == "" {
				return nil, fmt.Errorf("empty pipeline segment")
			}
			commands = append(commands, current)
			current = Command{}
		case ">":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing redirection target")
			}
			current.Output = tokens[i+1]
			current.Append = false
			i++
		case ">>":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing redirection target")
			}
			current.Output = tokens[i+1]
			current.Append = true
			i++
		case "2>":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing redirection target")
			}
			current.Stderr = tokens[i+1]
			current.StderrAppend = false
			i++
		case "2>>":
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("missing redirection target")
			}
			current.Stderr = tokens[i+1]
			current.StderrAppend = true
			i++
		default:
			if current.Name == "" {
				current.Name = token
			} else {
				current.Args = append(current.Args, token)
			}
		}
	}

	if current.Name != "" || len(current.Args) > 0 || len(commands) > 0 {
		commands = append(commands, current)
	}
	return commands, nil
}

func tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	var quote rune
	escaped := false
	tokenStarted := false

	flush := func() {
		if tokenStarted {
			tokens = append(tokens, current.String())
			current.Reset()
			tokenStarted = false
		}
	}

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if escaped {
			current.WriteRune(r)
			tokenStarted = true
			escaped = false
			continue
		}

		if r == '\\' {
			if quote == '\'' {
				current.WriteRune(r)
				tokenStarted = true
				continue
			}
			escaped = true
			continue
		}

		if quote != 0 {
			if r == quote {
				quote = 0
				continue
			}
			current.WriteRune(r)
			tokenStarted = true
			continue
		}

		switch {
		case unicode.IsSpace(r):
			flush()
		case r == '\'', r == '"':
			quote = r
			tokenStarted = true
		case r == '|':
			flush()
			tokens = append(tokens, "|")
		case r == '2' && i+1 < len(runes) && runes[i+1] == '>':
			flush()
			i++
			operator := "2>"
			if i+1 < len(runes) && runes[i+1] == '>' {
				operator = "2>>"
				i++
			}
			tokens = append(tokens, operator)
		case r == '>':
			flush()
			operator := ">"
			if i+1 < len(runes) && runes[i+1] == '>' {
				operator = ">>"
				i++
			}
			tokens = append(tokens, operator)
		default:
			current.WriteRune(r)
			tokenStarted = true
		}
	}

	if escaped {
		current.WriteRune('\\')
		tokenStarted = true
	}
	flush()
	return tokens
}
