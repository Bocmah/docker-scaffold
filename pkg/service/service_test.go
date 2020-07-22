package service_test

import (
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
	"testing"
)

func TestServicesConfig_IsPresent(t *testing.T) {
	conf := &service.ServicesConfig{}

	services := map[string]bool{
		"php":      false,
		"nodejs":   false,
		"nginx":    false,
		"database": false,
	}

	for s, expectedPresent := range services {
		if conf.IsPresent(s) != expectedPresent {
			t.Errorf("Service %s is present in empty configuration", s)
		}
	}

	conf = &service.ServicesConfig{
		PHP: &service.PHP{
			Version:    "7.4",
			Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
		},
		Nginx: &service.Nginx{
			Port:               80,
			ServerName:         "docker-scaffold",
			FastCGIPassPort:    9000,
			FastCGIReadTimeout: 60,
		},
		NodeJS: &service.NodeJS{
			Version: "10",
		},
		Database: &service.Database{
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

	services = map[string]bool{
		"php":      true,
		"nodejs":   true,
		"nginx":    true,
		"database": true,
	}

	for s, expectedPresent := range services {
		if conf.IsPresent(s) != expectedPresent {
			t.Errorf("Service %s is expected to be present in configuration %v", s, *conf)
		}
	}
}
