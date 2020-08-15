package template_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-scaffold/pkg/template"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

func TestRenderTemplatesFromConfiguration(t *testing.T) {
	dir, err := ioutil.TempDir("testdata", "output")

	defer os.RemoveAll(dir)

	outputPath := dir + "/.docker"

	conf := &service.FullConfig{
		AppName:     "awesome-app",
		ProjectRoot: "/home/test/app",
		OutputPath:  outputPath,
		Services: &service.ServicesConfig{
			PHP: &service.PHPConfig{
				Version:    "7.4",
				Extensions: []string{"mbstring", "exif", "pdo_mysql"},
			},
			Nginx: &service.NginxConfig{
				HTTPPort:   80,
				HTTPSPort:  443,
				ServerName: "awesomeapp",
				FastCGI: &service.FastCGI{
					PassPort:           9000,
					ReadTimeoutSeconds: 60,
				},
			},
			Database: &service.DatabaseConfig{
				System:  service.MySQL,
				Version: "8.0",
				Name:    "awesome-db",
				Port:    3306,
				Credentials: service.Credentials{
					Username:     "test-user",
					Password:     "test-password",
					RootPassword: "test-root-password",
				},
			},
			NodeJS: &service.NodeJSConfig{
				Version: "10",
			},
		},
	}

	rendered, err := template.RenderTemplatesFromConfiguration(conf)

	if err != nil {
		t.Fatalf("Encountered non-nil error in correct test case: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output path %v still doesn't exist after calling rendering function", outputPath)
	}

	absPath, absErr := filepath.Abs(outputPath + string(os.PathSeparator))

	if absErr != nil {
		t.Fatal(absErr)
	}

	wantRendered := template.RenderedServices{
		Services: map[service.SupportedService]*template.Rendered{
			service.PHP: {
				Path:        outputPath + string(os.PathSeparator) + "php/Dockerfile",
				CreatedDirs: []string{absPath + string(os.PathSeparator) + "php"},
			},
			service.Nginx: {
				Path:        outputPath + string(os.PathSeparator) + "nginx/conf.d/app.conf",
				CreatedDirs: []string{absPath + string(os.PathSeparator) + "nginx/conf.d"},
			},
			service.NodeJS: {
				Path:        outputPath + string(os.PathSeparator) + "nodejs/Dockerfile",
				CreatedDirs: []string{absPath + string(os.PathSeparator) + "nodejs"},
			},
		},
		CreatedDirs: []string{absPath},
	}

	if diff := cmp.Diff(wantRendered, rendered); diff != "" {
		t.Fatalf("RenderedServices mismatch (-want +got):\n%s", diff)
	}

	testFilesRoot := "testdata/template_render/.docker"
	testFiles := map[service.SupportedService]string{
		service.PHP:    testFilesRoot + string(os.PathSeparator) + "php/Dockerfile",
		service.Nginx:  testFilesRoot + string(os.PathSeparator) + "nginx/conf.d/app.conf",
		service.NodeJS: testFilesRoot + string(os.PathSeparator) + "nodejs/Dockerfile",
	}

	if diff := compareRenderedWithExpected(rendered, testFiles); diff != "" {
		t.Fatalf(diff)
	}
}

func compareRenderedWithExpected(renderedServices template.RenderedServices, testFiles map[service.SupportedService]string) (diff string) {
	for serv, expected := range testFiles {
		renderedService, ok := renderedServices.Services[serv]

		if !ok {
			return fmt.Sprintf("Service %s was not rendered", serv)
		}

		expectedFile, readExpectedErr := ioutil.ReadFile(expected)

		if readExpectedErr != nil {
			return fmt.Sprintf("Could not read expected file for service %s at path %s. Reason: %s", serv, expected, readExpectedErr)
		}

		renderedFile, readRenderedErr := ioutil.ReadFile(renderedService.Path)

		if readRenderedErr != nil {
			return fmt.Sprintf("Could not read expected file for service %s at path %s. Reason: %s", serv, renderedService.Path, readRenderedErr)
		}

		if diff := cmp.Diff(expectedFile, renderedFile); diff != "" {
			return fmt.Sprintf("Expected and rendered file for service %s mismatch (-want +got):\n%s", serv, diff)
		}
	}

	return ""
}
