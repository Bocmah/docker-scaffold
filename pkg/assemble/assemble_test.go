package assemble_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-gen/pkg/assemble"
	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"
)

func TestDockerCompose(t *testing.T) {
	conf := dummyConf()

	network := &dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge}
	want := &dockercompose.Config{
		Version: "3.8",
		Services: []*dockercompose.Service{
			{
				Name:          "php-fpm",
				Build:         &dockercompose.Build{Context: "/home/test/app", Dockerfile: "/home/test/app/.docker/php/Dockerfile"},
				Image:         &dockercompose.Image{Name: "test-app"},
				ContainerName: "test-app",
				WorkingDir:    "/var/www",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Volumes: dockercompose.ServiceVolumes{
					{Source: "/home/test/app", Target: "/var/www"},
				},
				Networks: dockercompose.ServiceNetworks{network},
			},
			{
				Name:          "webserver",
				Image:         &dockercompose.Image{Name: "nginx", Tag: "alpine"},
				ContainerName: "webserver",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					{Host: 80, Container: 80},
					{Host: 443, Container: 443},
				},
				Networks: dockercompose.ServiceNetworks{network},
				Volumes: dockercompose.ServiceVolumes{
					{Source: "/home/test/app", Target: "/var/www"},
					{Source: "/home/test/app/.docker/nginx/conf.d/app.conf", Target: "/etc/nginx/conf.d/app.conf"},
				},
			},
			{
				Name:          "db",
				Image:         &dockercompose.Image{Name: "mysql", Tag: "8.0"},
				ContainerName: "db",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					{Host: 3306, Container: 3306},
				},
				Environment: dockercompose.Environment{
					"MYSQL_DATABASE":      "test-db",
					"MYSQL_ROOT_PASSWORD": "secret-root",
					"MYSQL_USER":          "test-user",
					"MYSQL_PASSWORD":      "secret-password",
				},
				Volumes: dockercompose.ServiceVolumes{
					{Source: "test-app-data", Target: "/var/lib/mysql"},
				},
				Networks: dockercompose.ServiceNetworks{network},
			},
			{
				Name:          "nodejs",
				Build:         &dockercompose.Build{Context: "/home/test/app", Dockerfile: "/home/test/app/.docker/nodejs/Dockerfile"},
				ContainerName: "nodejs",
				Networks:      dockercompose.ServiceNetworks{network},
				Volumes: dockercompose.ServiceVolumes{
					{Source: "/home/test/app", Target: "/opt"},
				},
				WorkingDir: "/opt",
			},
		},
		Networks: dockercompose.Networks{
			&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
		},
		Volumes: dockercompose.NamedVolumes{
			&dockercompose.NamedVolume{Name: "test-app-data", Driver: dockercompose.VolumeDriverLocal},
		},
	}

	got := assemble.DockerCompose(conf)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("DockerCompose mismatch (-want +got):\n%s", diff)
	}
}
