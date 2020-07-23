package service_test

import (
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
	"reflect"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	got, err := service.LoadConfigFromFile("testdata/test.yaml")

	if err != nil {
		t.Errorf("Got error when loading correct config. Error - %v, Value - %v", err, got)
		return
	}

	want := &service.FullConfig{
		AppName:     "docker-scaffold",
		ProjectRoot: "/home/user/projects/test",
		OutputPath:  "/home/user/output",
		Services: &service.ServicesConfig{
			PHP: &service.PHP{
				Version:    "7.4",
				Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
			},
			Nginx: &service.Nginx{
				HttpPort:           80,
				HttpsPort:          443,
				ServerName:         "test-server",
				FastCGIPassPort:    9000,
				FastCGIReadTimeout: 60,
			},
			NodeJS: &service.NodeJS{
				Version: "10",
			},
			Database: &service.Database{
				System:  service.MySQL,
				Version: "5.7",
				Name:    "test-db",
				Port:    3306,
				Credentials: service.Credentials{
					Username:     "bocmah",
					Password:     "test",
					RootPassword: "testRoot",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Incorrectly loaded configuration. \nWant %v\nGot %v", want, got)
	}
}