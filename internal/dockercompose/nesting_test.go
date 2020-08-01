package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestNestingLevel_ApplyTo(t *testing.T) {
	tests := map[string]struct {
		input string
		level dockercompose.NestingLevel
		want  string
	}{
		"simple": {
			input: "test",
			level: dockercompose.NestingLevel(1),
			want:  "  test",
		},
		"multiline": {
			input: `line1
line2
line3`,
			level: dockercompose.NestingLevel(1),
			want: `  line1
  line2
  line3`,
		},
		"empty string": {
			input: "",
			level: dockercompose.NestingLevel(1),
			want:  "",
		},
		"nesting level above one": {
			input: "test",
			level: dockercompose.NestingLevel(2),
			want:  "    test",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.level.ApplyTo(tc.input)
			if tc.want != got {
				t.Fatalf("got: %v, want: %v", tc.want, got)
			}
		})
	}
}
