package dockercompose

import "testing"

func TestNetwork_String(t *testing.T) {
	tests := map[string]struct {
		input Network
		want  string
	}{
		"simple": {
			input: Network{Name: "test-network", Driver: Bridge},
			want: `test-network:
  driver: bridge`},
		"no driver": {
			input: Network{Name: "service-network"},
			want:  "",
		},
		"no name": {
			input: Network{Driver: Host},
			want:  "",
		},
		"no name and no driver": {
			input: Network{},
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
