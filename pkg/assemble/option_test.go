package assemble

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/pkg/service"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"
)

func newTestDockerComposeConfig() *dockercompose.Config {
	network := &dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge}

	return &dockercompose.Config{
		Version: "3.8",
		Services: []*dockercompose.Service{
			{
				Name:          "test-app",
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
}

func newTestServiceFiles() service.Files {
	return service.Files{
		service.PHP: []*service.File{
			{
				Type:         service.Dockerfile,
				PathOnHost:   "/home/test/app/.docker/php/Dockerfile",
				TemplatePath: "../../tmpl/php/php.dockerfile.gotmpl",
			},
		},
		service.Nginx: []*service.File{
			{
				Type:            service.ConfigFile,
				PathOnHost:      "/home/test/app/.docker/nginx/conf.d/app.conf",
				PathInContainer: "/etc/nginx/conf.d/app.conf",
				TemplatePath:    "../../tmpl/nginx/conf.gotmpl",
			},
		},
		service.NodeJS: []*service.File{
			{
				Type:         service.Dockerfile,
				PathOnHost:   "/home/test/app/.docker/nodejs/Dockerfile",
				TemplatePath: "../../tmpl/nodejs/nodejs.dockerfile.gotmpl",
			},
		},
	}
}

func newTestServiceEnv() service.Environment {
	return service.Environment{
		service.Database: {
			"MYSQL_DATABASE":      "test-db",
			"MYSQL_ROOT_PASSWORD": "secret-root",
			"MYSQL_USER":          "test-user",
			"MYSQL_PASSWORD":      "secret-password",
		},
	}
}

func TestOptionsAssembler(t *testing.T) {
	compose := newTestDockerComposeConfig()
	serviceFiles := newTestServiceFiles()
	serviceEnv := newTestServiceEnv()

	tests := map[string]struct {
		input service.SupportedService
		want  []Option
	}{
		"php": {
			input: service.PHP,
			want: []Option{
				networksOption{
					Networks: dockercompose.ServiceNetworks{&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge}},
				},
				dockerfilePathOption("/home/test/app/.docker/php/Dockerfile"),
			},
		},
		"nginx": {
			input: service.Nginx,
			want: []Option{
				networksOption{
					Networks: dockercompose.ServiceNetworks{&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge}},
				},
				volumesOption{
					Volumes: dockercompose.ServiceVolumes{&dockercompose.ServiceVolume{Source: "/home/test/app/.docker/nginx/conf.d/app.conf", Target: "/etc/nginx/conf.d/app.conf"}},
				},
			},
		},
		"nodejs": {
			input: service.NodeJS,
			want: []Option{
				networksOption{
					Networks: dockercompose.ServiceNetworks{&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge}},
				},
				dockerfilePathOption("/home/test/app/.docker/nodejs/Dockerfile"),
			},
		},
		"database": {
			input: service.Database,
			want: []Option{
				volumesOption{
					Volumes: dockercompose.ServiceVolumes{&dockercompose.ServiceVolume{Source: "test-app-data", Target: "/var/lib/mysql"}},
				},
				networksOption{
					Networks: dockercompose.ServiceNetworks{&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge}},
				},
				environmentOption{
					Environment: dockercompose.Environment{
						"MYSQL_DATABASE":      "test-db",
						"MYSQL_ROOT_PASSWORD": "secret-root",
						"MYSQL_USER":          "test-user",
						"MYSQL_PASSWORD":      "secret-password",
					},
				},
			},
		},
	}

	optsAssembler := &optionsAssembler{compose: compose, serviceFiles: serviceFiles, serviceEnv: serviceEnv, databaseSystemInUse: service.MySQL}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := optsAssembler.assembleForService(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("optsAssembler.assembleForService() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
