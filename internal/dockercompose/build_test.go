package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"
)

func TestBuild_Render(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.Build
		want  string
	}{
		"with context and dockerfile": {
			input: dockercompose.Build{Context: "/home/user/context", Dockerfile: "Dockerfile.debug"},
			want: `build:
  context: /home/user/context
  dockerfile: Dockerfile.debug`,
		},
		"no dockerfile": {
			input: dockercompose.Build{Context: "/home/user/context"},
			want:  "build: /home/user/context",
		},
		"no context": {
			input: dockercompose.Build{Dockerfile: "Dockerfile.debug"},
			want:  "",
		},
		"no context and no dockerfile": {
			input: dockercompose.Build{},
			want:  "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.Render()
			if tc.want != got {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
