package dockercompose

import "testing"

func TestNamedVolume_String(t *testing.T) {
	tests := map[string]struct {
		input NamedVolume
		want  string
	}{
		"with name and local driver": {
			input: NamedVolume{Name: "test-data", Driver: "local"},
			want:  "test-data:",
		},
		"with name and non-local driver": {
			input: NamedVolume{Name: "test-data", Driver: "foo"},
			want: `test-data:
  driver: foo`,
		},
		"without name": {
			input: NamedVolume{Driver: "local"},
			want:  "",
		},
		"without driver": {
			input: NamedVolume{Name: "test-data"},
			want:  "",
		},
		"without driver and name": {
			input: NamedVolume{},
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
