package dockercompose_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"
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
		input dockercompose.ServiceVolume
		want  string
	}{
		"simple": {
			input: dockercompose.ServiceVolume{Source: "/home/test", Target: "/var/test"},
			want:  "/home/test:/var/test",
		},
		"no source": {
			input: dockercompose.ServiceVolume{Target: "/var/test"},
			want:  "/var/test",
		},
		"no target": {
			input: dockercompose.ServiceVolume{Source: "/home/test"},
			want:  "",
		},
		"no source and no target": {
			input: dockercompose.ServiceVolume{},
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
				&dockercompose.ServiceVolume{Source: "/home/test", Target: "/var/test"},
				&dockercompose.ServiceVolume{Target: "/var/test"},
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

func TestNamedVolume_ToServiceVolume(t *testing.T) {
	tests := map[string]struct {
		input *dockercompose.NamedVolume
		want  *dockercompose.ServiceVolume
	}{
		"with name and driver": {
			input: &dockercompose.NamedVolume{Name: "test-data", Driver: "local"},
			want:  &dockercompose.ServiceVolume{Target: "test-data"},
		},
		"without name": {
			input: &dockercompose.NamedVolume{Driver: "local"},
			want:  nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.ToServiceVolume()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NamedVolume.ToServiceVolume() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamedVolumes_ToServiceVolumes(t *testing.T) {
	tests := map[string]struct {
		input dockercompose.NamedVolumes
		want  dockercompose.ServiceVolumes
	}{
		"simple": {
			input: dockercompose.NamedVolumes{
				&dockercompose.NamedVolume{Driver: dockercompose.VolumeDriverLocal, Name: "test-data"},
				&dockercompose.NamedVolume{Driver: dockercompose.VolumeDriverLocal, Name: "vol2"},
			},
			want: dockercompose.ServiceVolumes{
				&dockercompose.ServiceVolume{Target: "test-data"},
				&dockercompose.ServiceVolume{Target: "vol2"},
			},
		},
		"with unnamed volumes": {
			input: dockercompose.NamedVolumes{
				&dockercompose.NamedVolume{Driver: dockercompose.VolumeDriverLocal, Name: "test-data"},
				&dockercompose.NamedVolume{Driver: dockercompose.VolumeDriverLocal},
				&dockercompose.NamedVolume{Driver: dockercompose.VolumeDriverLocal, Name: "vol3"},
			},
			want: dockercompose.ServiceVolumes{
				&dockercompose.ServiceVolume{Target: "test-data"},
				&dockercompose.ServiceVolume{Target: "vol3"},
			},
		},
		"all unnamed volumes": {
			input: dockercompose.NamedVolumes{
				&dockercompose.NamedVolume{Driver: dockercompose.VolumeDriverLocal},
				&dockercompose.NamedVolume{Driver: dockercompose.VolumeDriverLocal},
			},
			want: dockercompose.ServiceVolumes{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.ToServiceVolumes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NamedVolume.ToServiceVolumes() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
