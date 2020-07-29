package dockercompose_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestNestingLevel_ApplyTo(t *testing.T) {
	tests := map[string]struct {
		input string
		level dockercompose.NestingLevel
		want  string
	}{
		"simple": {
			input: "test",
			level: dockercompose.NestingLevel(1),
			want:  "  test",
		},
		"multiline": {
			input: `line1
line2
line3`,
			level: dockercompose.NestingLevel(1),
			want: `  line1
  line2
  line3`,
		},
		"empty string": {
			input: "",
			level: dockercompose.NestingLevel(1),
			want:  "",
		},
		"nesting level above one": {
			input: "test",
			level: dockercompose.NestingLevel(2),
			want:  "    test",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.level.ApplyTo(tc.input)
			if tc.want != got {
				t.Fatalf("got: %v, want: %v", tc.want, got)
			}
		})
	}
}

func TestService_Render(t *testing.T) {
	service := dockercompose.Service{
		Name: "php",
		Build: dockercompose.Build{
			Context:    "/home/test",
			Dockerfile: "Dockerfile.test",
		},
		Image: dockercompose.Image{
			Name: "php",
			Tag:  "7.4",
		},
		ContainerName: "app",
		WorkingDir:    "/var/www",
		Restart:       dockercompose.RestartPolicyUnlessStopped,
		Environment: dockercompose.Environment{
			"SERVICE_NAME": "test-service",
		},
		Networks: dockercompose.Networks{
			dockercompose.Network{Name: "test-network", Driver: dockercompose.NetworkDriverBridge},
		},
		Volumes: dockercompose.Volumes{
			dockercompose.Volume{Source: "/home/test/app", Target: "/var/www"},
		},
	}

	want := `php:
  container_name: app
  working_dir: /var/www
  build:
    context: /home/test
    dockerfile: Dockerfile.test
  image: php:7.4
  restart: unless-stopped
  environment:
    SERVICE_NAME: test-service
  networks:
    - test-network
  volumes:
    - /home/test/app:/var/www`

	got := service.Render()

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
