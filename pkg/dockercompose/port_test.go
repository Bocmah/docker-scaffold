package dockercompose

import (
	"testing"
)

func TestPortsMapping_String(t *testing.T) {
	tests := map[string]struct {
		input *PortsMapping
		want  string
	}{
		"simple":                   {input: &PortsMapping{Host: "3306", Container: "33060"}, want: `"3306:33060"`},
		"empty host":               {input: &PortsMapping{Container: "8000"}, want: `"8000"`},
		"empty container":          {input: &PortsMapping{Host: "80"}, want: ""},
		"empty host and container": {input: &PortsMapping{}, want: ""},
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

func TestPorts_String(t *testing.T) {
	tests := map[string]struct {
		input Ports
		want  string
	}{
		"one mapping": {
			input: Ports{PortsMapping{Host: "90", Container: "9000"}},
			want: `ports:
  - "90:9000"`,
		},
		"two mappings": {
			input: Ports{PortsMapping{Host: "80", Container: "8080"}, PortsMapping{Host: "3000", Container: "3000"}},
			want: `ports:
  - "80:8080"
  - "3000:3000"`},
		"no mappings": {
			input: Ports{},
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
