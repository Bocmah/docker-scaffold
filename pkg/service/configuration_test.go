package service

import (
	"reflect"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	got, err := LoadConfigFromFile("testdata/test.yaml")

	if err != nil {
		t.Errorf("Got error when loading correct config. Error - %v, Value - %v", err, got)
		return
	}

	want := &Configuration{
		PHP: &PHP{
			Version: "7.4",
			Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
		},
		Nginx: &Nginx{
			Port: 80,
			ServerName: "scaffold",
			FastCGIPassPort: 9000,
			FastCGIReadTimeout: 60,
		},
		NodeJS: &NodeJS{
			Version: "10",
		},
		Database: &Database{
			System: MySQL,
			Version: "5.7",
			Name: "scaffold",
			Port: 3306,
			Credentials: Credentials{
				Username: "bocmah",
				Password: "test",
				RootPassword: "testRoot",
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Incorrectly loaded configuration. Want %+v, got %+v", want, got)
	}
}