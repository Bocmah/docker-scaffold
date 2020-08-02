package service_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

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

	conf = &service.ServicesConfig{
		PHP: &service.PHPConfig{
			Version:    "7.4",
			Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
		},
		Nginx: &service.NginxConfig{
			HTTPPort:   80,
			ServerName: "docker-scaffold",
			FastCGI: service.FastCGI{
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
			Name:    "docker-scaffold",
			Port:    3306,
			Credentials: service.Credentials{
				Username:     "bocmah",
				Password:     "test",
				RootPassword: "testRoot",
			},
		},
	}

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
