package core

import (
	"reflect"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []Command
	}{
		{
			name: "pipeline and output redirection",
			line: "echo hello | grep h > out.txt",
			want: []Command{
				{Name: "echo", Args: []string{"hello"}},
				{Name: "grep", Args: []string{"h"}, Output: "out.txt"},
			},
		},
		{
			name: "append redirections",
			line: "echo hello >> out.txt 2>> err.txt",
			want: []Command{
				{Name: "echo", Args: []string{"hello"}, Output: "out.txt", Append: true, Stderr: "err.txt", StderrAppend: true},
			},
		},
		{
			name: "quoted empty argument",
			line: `echo ""`,
			want: []Command{
				{Name: "echo", Args: []string{""}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommandLine(tt.line)
			if err != nil {
				t.Fatalf("ParseCommandLine() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommandLine() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
