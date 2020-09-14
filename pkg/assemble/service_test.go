package assemble_test

import (
	"fmt"
	"testing"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"
	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/pkg/assemble"
	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func assertAssemblerReturnsNilIfServiceIsNotPresent(t *testing.T, assembler assemble.ServiceAssembler, serv service.SupportedService) {
	t.Helper()

	conf := &service.FullConfig{
		AppName:     "Test App",
		ProjectRoot: "/home/test/app",
		OutputPath:  "/home/test/app/.docker",
		Services:    &service.ServicesConfig{},
	}

	got := assembler(conf)

	if got != nil {
		t.Errorf("assembler did not return nil for conf without %s. Instead it returned %v", serv, got)
	}
}

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
				HTTPPort:   80,
				HTTPSPort:  443,
				ServerName: "test-server",
				FastCGI: &service.FastCGI{
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
					Password:     "secret-password",
					Username:     "test-user",
				},
			},
		},
	}
}

func TestPhpAssemble(t *testing.T) {
	assembler := assemble.NewServiceAssembler(service.PHP)

	assertAssemblerReturnsNilIfServiceIsNotPresent(t, assembler, service.PHP)

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
				assemble.WithNetworks(dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				}),
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

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := assembler(conf, tc.opts...)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("PHP assembler mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNginxAssemble(t *testing.T) {
	assembler := assemble.NewServiceAssembler(service.Nginx)

	assertAssemblerReturnsNilIfServiceIsNotPresent(t, assembler, service.Nginx)

	conf := dummyConf()
	confWithoutPorts := dummyConf()

	confWithoutPorts.Services.Nginx.HTTPPort = 0
	confWithoutPorts.Services.Nginx.HTTPSPort = 0

	tests := map[string]struct {
		conf *service.FullConfig
		opts []assemble.Option
		want *dockercompose.Service
	}{
		"no options": {
			conf: conf,
			want: &dockercompose.Service{
				Name: "webserver",
				Image: &dockercompose.Image{
					Name: "nginx",
					Tag:  "alpine",
				},
				ContainerName: "webserver",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					&dockercompose.PortsMapping{Host: 80, Container: 80},
					&dockercompose.PortsMapping{Host: 443, Container: 443},
				},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "/home/test/app", Target: "/var/www"},
				},
			},
		},
		"with options": {
			conf: conf,
			opts: []assemble.Option{
				assemble.WithNetworks(dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				}),
				assemble.WithVolumes(dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "./nginx/conf.d/", Target: "/etc/nginx/conf.d/"},
				}),
			},
			want: &dockercompose.Service{
				Name: "webserver",
				Image: &dockercompose.Image{
					Name: "nginx",
					Tag:  "alpine",
				},
				ContainerName: "webserver",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					&dockercompose.PortsMapping{Host: 80, Container: 80},
					&dockercompose.PortsMapping{Host: 443, Container: 443},
				},
				Networks: dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "/home/test/app", Target: "/var/www"},
					&dockercompose.ServiceVolume{Source: "./nginx/conf.d/", Target: "/etc/nginx/conf.d/"},
				},
			},
		},
		"without ports": {
			conf: confWithoutPorts,
			want: &dockercompose.Service{
				Name: "webserver",
				Image: &dockercompose.Image{
					Name: "nginx",
					Tag:  "alpine",
				},
				ContainerName: "webserver",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					&dockercompose.PortsMapping{Host: 80, Container: 80},
					&dockercompose.PortsMapping{Host: 443, Container: 443},
				},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "/home/test/app", Target: "/var/www"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := assembler(tc.conf, tc.opts...)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("nginx assembler mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDatabaseAssemble(t *testing.T) {
	assembler := assemble.NewServiceAssembler(service.Database)

	assertAssemblerReturnsNilIfServiceIsNotPresent(t, assembler, service.Database)

	conf := dummyConf()

	tests := map[string]struct {
		opts []assemble.Option
		want *dockercompose.Service
	}{
		"no options": {
			want: &dockercompose.Service{
				Name: "db",
				Image: &dockercompose.Image{
					Name: string(service.MySQL),
					Tag:  "8.0",
				},
				ContainerName: "db",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					&dockercompose.PortsMapping{Host: 3306, Container: 3306},
				},
			},
		},
		"with options": {
			opts: []assemble.Option{
				assemble.WithNetworks(dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				}),
				assemble.WithVolumes(dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "test-data", Target: "/var/lib/mysql"},
				}),
				assemble.WithEnvironment(dockercompose.Environment{
					"MYSQL_DATABASE":      "test-db",
					"MYSQL_ROOT_PASSWORD": "secret-root",
				}),
			},
			want: &dockercompose.Service{
				Name: "db",
				Image: &dockercompose.Image{
					Name: string(service.MySQL),
					Tag:  "8.0",
				},
				ContainerName: "db",
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					&dockercompose.PortsMapping{Host: 3306, Container: 3306},
				},
				Environment: dockercompose.Environment{
					"MYSQL_DATABASE":      "test-db",
					"MYSQL_ROOT_PASSWORD": "secret-root",
				},
				Networks: dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: "test-data", Target: "/var/lib/mysql"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := assembler(conf, tc.opts...)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("Database assembler mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNodeJSAssemble(t *testing.T) {
	assembler := assemble.NewServiceAssembler(service.NodeJS)

	assertAssemblerReturnsNilIfServiceIsNotPresent(t, assembler, service.NodeJS)

	conf := dummyConf()

	tests := map[string]struct {
		opts []assemble.Option
		want *dockercompose.Service
	}{
		"no options": {
			want: &dockercompose.Service{
				Name: "nodejs",
				Image: &dockercompose.Image{
					Name: "node",
					Tag:  "alpine",
				},
				ContainerName: "nodejs",
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: conf.ProjectRoot, Target: "/opt"},
				},
				WorkingDir: "/opt",
			},
		},
		"with options": {
			opts: []assemble.Option{
				assemble.WithDockerfilePath("/home/test/app/.docker/node/Dockerfile"),
				assemble.WithNetworks(dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				}),
			},
			want: &dockercompose.Service{
				Name: "nodejs",
				Build: &dockercompose.Build{
					Context:    conf.ProjectRoot,
					Dockerfile: "/home/test/app/.docker/node/Dockerfile",
				},
				ContainerName: "nodejs",
				Networks: dockercompose.ServiceNetworks{
					&dockercompose.Network{Name: "test-app-network", Driver: dockercompose.NetworkDriverBridge},
				},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: conf.ProjectRoot, Target: "/opt"},
				},
				WorkingDir: "/opt",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := assembler(conf, tc.opts...)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("Database assembler mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUnknownAssemble(t *testing.T) {
	assembler := assemble.NewServiceAssembler(service.SupportedService(100))

	got := assembler(dummyConf())

	if got != nil {
		t.Fatalf("Unknown assembler did not return nil. Returned %v instead", got)
	}
}
