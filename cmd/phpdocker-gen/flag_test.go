package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseFlagsCorrect(t *testing.T) {
	var tests = []struct {
		args []string
		conf Config
	}{
		{
			[]string{},
			Config{file: "", args: []string{}},
		},
		{
			[]string{"something"},
			Config{file: "", args: []string{"something"}},
		},
		{
			[]string{"-file", "path/to/file"},
			Config{file: "path/to/file", args: []string{}}},

		{
			[]string{"-file", "path/to/file", "another/path/to/file"},
			Config{file: "path/to/file", args: []string{"another/path/to/file"}},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			conf, output, err := parseFlags("test", tt.args)
			if err != nil {
				t.Errorf("err got %v, want nil", err)
			}
			if output != "" {
				t.Errorf("output got %q, want empty", output)
			}
			if !reflect.DeepEqual(*conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v", *conf, tt.conf)
			}
		})
	}
}
