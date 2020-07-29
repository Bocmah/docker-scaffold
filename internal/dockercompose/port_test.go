package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestPortsMapping_Render(t *testing.T) {
	tests := map[string]struct {
		input *dockercompose.PortsMapping
		want  string
	}{
		"simple":                   {input: &dockercompose.PortsMapping{Host: 3306, Container: 33060}, want: `"3306:33060"`},
		"empty host":               {input: &dockercompose.PortsMapping{Container: 8000}, want: `"8000"`},
		"empty container":          {input: &dockercompose.PortsMapping{Host: 80}, want: ""},
		"empty host and container": {input: &dockercompose.PortsMapping{}, want: ""},
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

func TestPorts_Render(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.Ports
		want  string
	}{
		"one mapping": {
			input: dockercompose.Ports{dockercompose.PortsMapping{Host: 90, Container: 9000}},
			want: `ports:
  - "90:9000"`,
		},
		"two mappings": {
			input: dockercompose.Ports{
				dockercompose.PortsMapping{Host: 80, Container: 8080},
				dockercompose.PortsMapping{Host: 3000, Container: 3000},
			},
			want: `ports:
  - "80:8080"
  - "3000:3000"`},
		"no mappings": {
			input: dockercompose.Ports{},
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
