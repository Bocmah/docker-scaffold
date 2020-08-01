package dockercompose_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
)

func TestConfig_Render(t *testing.T) {
	network := dockercompose.Network{Name: "test-network", Driver: dockercompose.NetworkDriverBridge}
	networks := dockercompose.ServiceNetworks{&network}
	rootMount := dockercompose.Volume{Source: "/home/test/app", Target: "/var/www"}
	namedVol := dockercompose.NamedVolume{Name: "test-data", Driver: dockercompose.VolumeDriverLocal}

	php := dockercompose.Service{
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
		Networks: networks,
		Volumes:  dockercompose.ServiceVolumes{&rootMount},
	}

	server := dockercompose.Service{
		Name: "webserver",
		Image: &dockercompose.Image{
			Name: "nginx",
			Tag:  "alpine",
		},
		ContainerName: "webserver",
		Restart:       dockercompose.RestartPolicyUnlessStopped,
		Ports: dockercompose.Ports{
			&dockercompose.PortsMapping{Container: 80, Host: 80},
			&dockercompose.PortsMapping{Container: 443, Host: 443},
		},
		Networks: networks,
		Volumes:  dockercompose.ServiceVolumes{&rootMount, &dockercompose.Volume{Source: "./nginx/conf.d/", Target: "/etc/nginx/conf.d/"}},
	}

	db := dockercompose.Service{
		Name: "db",
		Image: &dockercompose.Image{
			Name: "mysql",
			Tag:  "5.7.22",
		},
		ContainerName: "db",
		Restart:       dockercompose.RestartPolicyUnlessStopped,
		Ports: dockercompose.Ports{
			&dockercompose.PortsMapping{Container: 3306, Host: 3306},
		},
		Environment: dockercompose.Environment{
			"MYSQL_ROOT_PASSWORD": "secret",
		},
		Networks: networks,
		Volumes: dockercompose.ServiceVolumes{
			&dockercompose.Volume{Source: namedVol.Name, Target: "/var/lib/mysql"},
		},
	}

	conf := dockercompose.Config{
		Version:  "3.8",
		Services: []*dockercompose.Service{&php, &server, &db},
		Networks: []*dockercompose.Network{&network},
		Volumes:  []*dockercompose.NamedVolume{&namedVol},
	}

	want := `version: "3.8"
services:
  php:
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
      - /home/test/app:/var/www
  webserver:
    container_name: webserver
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    networks:
      - test-network
    volumes:
      - /home/test/app:/var/www
      - ./nginx/conf.d/:/etc/nginx/conf.d/
  db:
    container_name: db
    image: mysql:5.7.22
    restart: unless-stopped
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: secret
    networks:
      - test-network
    volumes:
      - test-data:/var/lib/mysql
networks:
  test-network:
    driver: bridge
volumes:
  test-data:`

	got := conf.Render()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("conf.Render() mismatch (-want +got):\n%s", diff)
	}
}
