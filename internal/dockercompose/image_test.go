package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestImage_Render(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.Image
		want  string
	}{
		"simple": {
			input: dockercompose.Image{Name: "nginx", Tag: "alpine"},
			want:  "image: nginx:alpine",
		},
		"no tag": {
			input: dockercompose.Image{Name: "nginx"},
			want:  "image: nginx",
		},
		"no tag and no name": {
			input: dockercompose.Image{},
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
