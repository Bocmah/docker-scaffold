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

func TestParseFlagsError(t *testing.T) {
	var tests = []struct {
		args   []string
		errstr string
	}{
		{[]string{"-file"}, "flag needs an argument"},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			conf, output, err := parseFlags("test", tt.args)
			if conf != nil {
				t.Errorf("conf got %v, want nil", conf)
			}
			if err == nil {
				t.Fatalf("err got nil, want %q", tt.errstr)
			}
			if strings.Index(err.Error(), tt.errstr) < 0 {
				t.Errorf("err got %q, want to find %q", err.Error(), tt.errstr)
			}
			if strings.Index(output, "Usage of test") < 0 {
				t.Errorf("output got %q", output)
			}
		})
	}
}
