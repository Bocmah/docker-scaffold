package dockercompose_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestNetwork_String(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.Network
		want  string
	}{
		"simple": {
			input: dockercompose.Network{Name: "test-network", Driver: dockercompose.NetworkDriverBridge},
			want: `test-network:
  driver: bridge`},
		"no driver": {
			input: dockercompose.Network{Name: "service-network"},
			want:  "",
		},
		"no name": {
			input: dockercompose.Network{Driver: dockercompose.NetworkDriverHost},
			want:  "",
		},
		"no name and no driver": {
			input: dockercompose.Network{},
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

func TestNetworks_Render(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.ServiceNetworks
		want  string
	}{
		"simple": {
			input: dockercompose.ServiceNetworks{
				&dockercompose.Network{Name: "test-data", Driver: dockercompose.NetworkDriverBridge},
				&dockercompose.Network{Name: "test-data-1", Driver: dockercompose.NetworkDriverHost},
			},
			want: `networks:
  - test-data
  - test-data-1`},
		"empty": {
			input: dockercompose.ServiceNetworks{},
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

func TestNetworks_ToServiceNetworks(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.Networks
		want  dockercompose.ServiceNetworks
	}{
		"simple": {
			input: dockercompose.Networks{
				&dockercompose.Network{Driver: dockercompose.NetworkDriverBridge, Name: "test-network"},
				&dockercompose.Network{Driver: dockercompose.NetworkDriverBridge, Name: "test-network2"},
			},
			want: dockercompose.ServiceNetworks{
				&dockercompose.Network{Driver: dockercompose.NetworkDriverBridge, Name: "test-network"},
				&dockercompose.Network{Driver: dockercompose.NetworkDriverBridge, Name: "test-network2"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.ToServiceNetworks()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NamedVolume.ToServiceNetworks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
