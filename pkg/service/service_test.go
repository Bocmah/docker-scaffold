package service_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func dummyConfigWithAllServices() *service.ServicesConfig {
	return &service.ServicesConfig{
		PHP: &service.PHPConfig{
			Version:    "7.4",
			Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
		},
		Nginx: &service.NginxConfig{
			HTTPPort:   80,
			ServerName: "phpdocker-gen",
			FastCGI: &service.FastCGI{
				PassPort:           9000,
				ReadTimeoutSeconds: 60,
			},
		},
		NodeJS: &service.NodeJSConfig{
			Version: "10",
		},
		Database: &service.DatabaseConfig{
			System:  service.MySQL,
			Version: "5.7",
			Name:    "phpdocker-gen",
			Port:    3306,
			Credentials: service.Credentials{
				Username:     "bocmah",
				Password:     "test",
				RootPassword: "testRoot",
			},
		},
	}
}

func TestServicesConfig_IsPresent(t *testing.T) {
	conf := &service.ServicesConfig{}

	services := map[service.SupportedService]bool{
		service.PHP:      false,
		service.NodeJS:   false,
		service.Nginx:    false,
		service.Database: false,
	}

	for s, expectedPresent := range services {
		if conf.IsPresent(s) != expectedPresent {
			t.Errorf("Service %s is present in empty configuration", s)
		}
	}

	conf = dummyConfigWithAllServices()

	services = map[service.SupportedService]bool{
		service.PHP:      true,
		service.NodeJS:   true,
		service.Nginx:    true,
		service.Database: true,
	}

	for s, expectedPresent := range services {
		if conf.IsPresent(s) != expectedPresent {
			t.Errorf("Service %s is expected to be present in configuration %v", s, *conf)
		}
	}
}

func TestServicesConfig_PresentServicesCount(t *testing.T) {
	tests := map[string]struct {
		input *service.ServicesConfig
		want  int
	}{
		"all services": {
			input: dummyConfigWithAllServices(),
			want:  4,
		},
		"one service": {
			input: &service.ServicesConfig{
				PHP: &service.PHPConfig{
					Version:    "7.4",
					Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
				},
			},
			want: 1,
		},
		"no services": {
			input: &service.ServicesConfig{},
			want:  0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.PresentServicesCount()
			if tc.want != got {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
