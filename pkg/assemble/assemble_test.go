package assemble_test

import (
	"fmt"
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-scaffold/pkg/assemble"
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

func dummyConf() *service.FullConfig {
	return &service.FullConfig{
		AppName:     "Test App",
		ProjectRoot: "/home/test/app",
		OutputPath:  "/home/test/app/.docker",
		Services: &service.ServicesConfig{
			PHP: &service.PHPConfig{
				Version:    "7.4",
				Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd"},
			},
			NodeJS: &service.NodeJSConfig{
				Version: "latest",
			},
			Nginx: &service.NginxConfig{
				HttpPort:   80,
				HttpsPort:  443,
				ServerName: "test-server",
				FastCGI: service.FastCGI{
					PassPort:           9000,
					ReadTimeoutSeconds: 60,
				},
			},
			Database: &service.DatabaseConfig{
				Version: "8.0",
				System:  service.MySQL,
				Name:    "test-db",
				Port:    3306,
				Credentials: service.Credentials{
					RootPassword: "secret-root",
					Password:     "secret-user",
					Username:     "test-user",
				},
			},
		},
	}
}

func TestPhpAssemble(t *testing.T) {
	conf := dummyConf()

	tests := map[string]struct {
		opts []assemble.Option
		want *dockercompose.Service
	}{
		"no options": {
			want: &dockercompose.Service{
				Name: "test-app",
				Image: &dockercompose.Image{
					Name: "php",
					Tag:  fmt.Sprintf("%s-fpm", conf.Services.PHP.Version),
				},
				ContainerName: "test-app",
				WorkingDir:    "/var/www",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "/home/test/app", Target: "/var/www"},
				},
			},
		},
		"with options": {
			opts: []assemble.Option{
				assemble.WithDockerfilePath("/home/test/app/.docker/php/Dockerfile"),
				assemble.WithSharedNetwork(&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge}),
			},
			want: &dockercompose.Service{
				Name: "test-app",
				Build: &dockercompose.Build{
					Context:    "/home/test/app",
					Dockerfile: "/home/test/app/.docker/php/Dockerfile",
				},
				Image: &dockercompose.Image{
					Name: "test-app",
				},
				ContainerName: "test-app",
				WorkingDir:    "/var/www",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Networks: dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "/home/test/app", Target: "/var/www"},
				},
			},
		},
	}

	assembler := assemble.NewServiceAssembler(service.PHP)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := assembler(conf, tc.opts...)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("assembler mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
