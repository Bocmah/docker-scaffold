package dockercompose_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestService_Render(t *testing.T) {
	service := dockercompose.Service{
		Name: "php",
		Build: &dockercompose.Build{
			Context:    "/home/test",
			Dockerfile: "Dockerfile.test",
		},
		Image: &dockercompose.Image{
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
			&dockercompose.Network{Name: "test-network", Driver: dockercompose.NetworkDriverBridge},
		},
		Volumes: dockercompose.Volumes{
			&dockercompose.Volume{Source: "/home/test/app", Target: "/var/www"},
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

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("service.Render() mismatch (-want +got):\n%s", diff)
	}
}
