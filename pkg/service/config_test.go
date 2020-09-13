package service_test

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func TestLoadConfigFromFile(t *testing.T) {
	got, err := service.LoadConfigFromFile("testdata/test.yaml")

	if err != nil {
		t.Errorf("Got error when loading correct config. Error - %v, Value - %v", err, got)
		return
	}

	want := &service.FullConfig{
		AppName:     "phpdocker-gen",
		ProjectRoot: "/home/user/projects/test",
		OutputPath:  "/home/user/output",
		Services: &service.ServicesConfig{
			PHP: &service.PHPConfig{
				Version:    "7.4",
				Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
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
			NodeJS: &service.NodeJSConfig{
				Version: "10",
			},
			Database: &service.DatabaseConfig{
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

func TestFullConfigValid_Validate(t *testing.T) {
	conf := &service.FullConfig{
		AppName:     "phpdocker-gen",
		ProjectRoot: "/home/user/projects/test",
		OutputPath:  "/home/user/output",
		Services: &service.ServicesConfig{
			PHP: &service.PHPConfig{
				Version:    "7.4",
				Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
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
			NodeJS: &service.NodeJSConfig{
				Version: "10",
			},
			Database: &service.DatabaseConfig{
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

	validationErr := conf.Validate()

	if validationErr != nil {
		t.Fatalf("Encountered non-nil validation error on valid config: %s", validationErr)
	}
}

func TestFullConfigInvalid_Validate(t *testing.T) {
	conf := &service.FullConfig{
		OutputPath: "/home/user/output",
		Services: &service.ServicesConfig{
			PHP: &service.PHPConfig{
				Version:    "7.4",
				Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
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
			NodeJS: &service.NodeJSConfig{
				Version: "10",
			},
			Database: &service.DatabaseConfig{
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

	validationErr := conf.Validate()

	if validationErr == nil {
		t.Fatalf("Encountered nil validation error on invalid config: %s", validationErr)
	}

	stringErr := validationErr.Error()

	expectedErrs := []string{"App name is required", "Project root is required"}

	for _, expectedErr := range expectedErrs {
		if !strings.Contains(stringErr, expectedErr) {
			t.Fatalf("validation err %s does not contain expected err %s", stringErr, expectedErr)
		}
	}
}

func TestFullConfig_GetServiceFiles(t *testing.T) {
	outputPath := "/home/user/output"

	conf := &service.FullConfig{
		AppName:     "phpdocker-gen",
		ProjectRoot: "/home/user/projects/test",
		OutputPath:  outputPath,
		Services: &service.ServicesConfig{
			PHP: &service.PHPConfig{
				Version:    "7.4",
				Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
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
			NodeJS: &service.NodeJSConfig{
				Version: "10",
			},
			Database: &service.DatabaseConfig{
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

	want := service.Files{
		service.PHP: []*service.File{
			{
				Type:         service.Dockerfile,
				PathOnHost:   filepath.Join(outputPath, "php/Dockerfile"),
				TemplatePath: filepath.Join("../../tmpl", "php/php.dockerfile.gotmpl"),
			},
		},
		service.Nginx: []*service.File{
			{
				Type:            service.ConfigFile,
				PathOnHost:      filepath.Join(outputPath, "nginx/conf.d/app.conf"),
				PathInContainer: "/etc/nginx/conf.d/app.conf",
				TemplatePath:    filepath.Join("../../tmpl", "nginx/conf.gotmpl"),
			},
		},
		service.NodeJS: []*service.File{
			{
				Type:         service.Dockerfile,
				PathOnHost:   filepath.Join(outputPath, "nodejs/Dockerfile"),
				TemplatePath: filepath.Join("../../tmpl", "nodejs/nodejs.dockerfile.gotmpl"),
			},
		},
	}

	got := conf.GetServiceFiles()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("conf.GetServiceFiles() mismatch (-want +got):\n%s", diff)
	}
}

func TestFullConfig_GetEnvironment(t *testing.T) {
	conf := &service.FullConfig{
		AppName:     "phpdocker-gen",
		ProjectRoot: "/home/user/projects/test",
		OutputPath:  "/home/user/output",
		Services: &service.ServicesConfig{
			PHP: &service.PHPConfig{
				Version:    "7.4",
				Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd", "pdo_mysql"},
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
			NodeJS: &service.NodeJSConfig{
				Version: "10",
			},
		},
	}

	env := conf.GetEnvironment()

	if env != nil {
		t.Errorf("encountered non-nil environment when nil is expected: %s", env)
	}

	conf.Services.Database = &service.DatabaseConfig{
		System:  service.MySQL,
		Version: "5.7",
		Name:    "test-db",
		Port:    3306,
		Credentials: service.Credentials{
			Username:     "bocmah",
			Password:     "test",
			RootPassword: "testRoot",
		},
	}

	env = conf.GetEnvironment()

	wantEnv := service.Environment{
		service.Database: {
			"MYSQL_ROOT_PASSWORD": "testRoot",
			"MYSQL_DATABASE":      "test-db",
			"MYSQL_USER":          "bocmah",
			"MYSQL_PASSWORD":      "test",
		},
	}

	if diff := cmp.Diff(wantEnv, env); diff != "" {
		t.Errorf("conf.GetEnvironment() mismatch (-want +got):\n%s", diff)
	}
}
