package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestDoubleQuotted(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"simple": {
			input: "test",
			want:  `"test"`,
		},
		"empty string": {
			input: "",
			want:  "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := dockercompose.DoubleQuotted(tc.input)
			if tc.want != got {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestMapping(t *testing.T) {
	tests := map[string]struct {
		str1 string
		str2 string
		want string
	}{
		"simple": {
			str1: "test",
			str2: "mapping",
			want: "test:mapping",
		},
		"empty str1": {
			str1: "",
			str2: "mapping",
			want: "mapping",
		},
		"empty str2": {
			str1: "test",
			str2: "",
			want: "test",
		},
		"empty str1 and str2": {
			str1: "",
			str2: "",
			want: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := dockercompose.Mapping(tc.str1, tc.str2)
			if tc.want != got {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
