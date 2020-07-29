package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestEnvironment_Render(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.Environment
		want  string
	}{
		"simple": {
			input: dockercompose.Environment{
				"SOME_VAR": "foo",
			},
			want: `environment:
  SOME_VAR: foo`,
		},
		"with empty value": {
			input: dockercompose.Environment{
				"SOME_VAR": "",
			},
			want: `environment:
  SOME_VAR:`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.Render()
			if tc.want != got {
				t.Fatalf("expected:\n %v\n got:\n %v", tc.want, got)
			}
		})
	}
}
