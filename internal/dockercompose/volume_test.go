package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestNamedVolume_String(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.NamedVolume
		want  string
	}{
		"with name and local driver": {
			input: dockercompose.NamedVolume{Name: "test-data", Driver: "local"},
			want:  "test-data:",
		},
		"with name and non-local driver": {
			input: dockercompose.NamedVolume{Name: "test-data", Driver: "foo"},
			want: `test-data:
  driver: foo`,
		},
		"without name": {
			input: dockercompose.NamedVolume{Driver: "local"},
			want:  "",
		},
		"without driver": {
			input: dockercompose.NamedVolume{Name: "test-data"},
			want:  "",
		},
		"without driver and name": {
			input: dockercompose.NamedVolume{},
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

func TestVolume_String(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.Volume
		want  string
	}{
		"simple": {
			input: dockercompose.Volume{Source: "/home/test", Target: "/var/test"},
			want:  "/home/test:/var/test",
		},
		"no source": {
			input: dockercompose.Volume{Target: "/var/test"},
			want:  "/var/test",
		},
		"no target": {
			input: dockercompose.Volume{Source: "/home/test"},
			want:  "",
		},
		"no source and no target": {
			input: dockercompose.Volume{},
			want:  "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.String()
			if tc.want != got {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestVolumes_Render(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.ServiceVolumes
		want  string
	}{
		"simple": {
			input: dockercompose.ServiceVolumes{
				&dockercompose.Volume{Source: "/home/test", Target: "/var/test"},
				&dockercompose.Volume{Target: "/var/test"},
			},
			want: `volumes:
  - /home/test:/var/test
  - /var/test`},
		"empty": {
			input: dockercompose.ServiceVolumes{},
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
