package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/pkg/render"

	"gopkg.in/yaml.v2"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
	"github.com/spf13/afero"
)

func createFileWithTestConf(t *testing.T, fs afero.Fs) afero.File {
	t.Helper()

	testConf := map[string]interface{}{
		"appName":     "test-app",
		"projectRoot": "/home/user/projects/test",
		"outputPath":  "/home/user/output",
		"services": map[interface{}]interface{}{
			"php": map[interface{}]interface{}{
				"version": "7.4",
				"extensions": []interface{}{
					"mbstring",
					"zip",
					"exif",
					"pcntl",
					"gd",
				},
			},
			"nginx": map[interface{}]interface{}{
				"httpPort":   80,
				"serverName": "test-server",
				"fastCGI": map[interface{}]interface{}{
					"passPort":           9000,
					"readTimeoutSeconds": 50,
				},
			},
			"nodejs": map[interface{}]interface{}{
				"version": "10",
			},
			"database": map[interface{}]interface{}{
				"system":       "mysql",
				"version":      "5.7",
				"name":         "test-db",
				"port":         3306,
				"username":     "bocmah",
				"password":     "test",
				"rootPassword": "testRoot",
			},
		},
	}

	yamlTestConf, marshalErr := yaml.Marshal(testConf)

	if marshalErr != nil {
		t.Fatalf("failed to marshal: %s", marshalErr)
	}

	tmpfile, tempFileErr := afero.TempFile(fs, "", "*.yaml")

	if tempFileErr != nil {
		t.Fatalf("failed to create tempfile: %s", tempFileErr)
	}

	if _, err := tmpfile.Write(yamlTestConf); err != nil {
		t.Fatalf("failed to write to tempfile: %s", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close tempfile: %s", err)
	}

	return tmpfile
}

func createTestDockerComposeConf() map[string]interface{} {
	const projectRoot = "/home/user/projects/test"
	const outputPath = "/home/user/output"
	const networkName = "test-app-network"

	return map[string]interface{}{
		"version": "3.8",
		"services": map[interface{}]interface{}{
			"php-fpm": map[interface{}]interface{}{
				"container_name": "test-app",
				"working_dir":    "/var/www",
				"build": map[interface{}]interface{}{
					"context":    projectRoot,
					"dockerfile": filepath.Join(outputPath, "/php/Dockerfile"),
				},
				"image":   "test-app",
				"restart": string(dockercompose.RestartPolicyUnlessStopped),
				"networks": []interface{}{
					networkName,
				},
				"volumes": []interface{}{
					projectRoot + ":/var/www",
				},
			},
			"webserver": map[interface{}]interface{}{
				"container_name": "webserver",
				"image":          "nginx:alpine",
				"ports": []interface{}{
					"80:80",
					"443:443",
				},
				"networks": []interface{}{
					networkName,
				},
				"volumes": []interface{}{
					projectRoot + ":/var/www",
					filepath.Join(outputPath, "/nginx/conf.d/app.conf") + ":/etc/nginx/conf.d/app.conf",
				},
				"restart": string(dockercompose.RestartPolicyUnlessStopped),
			},
			"db": map[interface{}]interface{}{
				"container_name": "db",
				"image":          "mysql:5.7",
				"restart":        string(dockercompose.RestartPolicyUnlessStopped),
				"ports": []interface{}{
					"3306:3306",
				},
				"environment": map[interface{}]interface{}{
					"MYSQL_DATABASE":      "test-db",
					"MYSQL_ROOT_PASSWORD": "testRoot",
					"MYSQL_USER":          "bocmah",
					"MYSQL_PASSWORD":      "test",
				},
				"networks": []interface{}{
					networkName,
				},
				"volumes": []interface{}{
					"test-app-data:/var/lib/mysql",
				},
			},
			"nodejs": map[interface{}]interface{}{
				"container_name": "nodejs",
				"working_dir":    "/opt",
				"build": map[interface{}]interface{}{
					"context":    projectRoot,
					"dockerfile": filepath.Join(outputPath, "/nodejs/Dockerfile"),
				},
				"networks": []interface{}{
					networkName,
				},
				"volumes": []interface{}{
					projectRoot + ":/opt",
				},
			},
		},
		"networks": map[interface{}]interface{}{
			networkName: map[interface{}]interface{}{
				"driver": string(dockercompose.NetworkDriverBridge),
			},
		},
		"volumes": map[interface{}]interface{}{
			"test-app-data": nil,
		},
	}
}

func TestApp(t *testing.T) {
	fs := afero.NewMemMapFs()

	testConf := createFileWithTestConf(t, fs)

	AppFs = fs
	service.AppFs = fs
	render.AppFs = fs

	const programName = "phpdocker-gen"
	const fileFlag = "-file"

	os.Args = []string{programName, fileFlag, testConf.Name()}

	const testFilesRoot = "testdata/template_render/.docker"
	const actualFilesRoot = "/home/user/output"
	files := map[service.SupportedService]struct {
		pathToTest   string
		pathToActual string
	}{
		service.PHP: {
			pathToTest:   filepath.Join(testFilesRoot, "php/Dockerfile"),
			pathToActual: filepath.Join(actualFilesRoot, "php/Dockerfile"),
		},
		service.Nginx: {
			pathToTest:   filepath.Join(testFilesRoot, "nginx/conf.d/app.conf"),
			pathToActual: filepath.Join(actualFilesRoot, "nginx/conf.d/app.conf"),
		},
		service.NodeJS: {
			pathToTest:   filepath.Join(testFilesRoot, "nodejs/Dockerfile"),
			pathToActual: filepath.Join(actualFilesRoot, "nodejs/Dockerfile"),
		},
	}

	main()

	for _, paths := range files {
		if _, statErr := fs.Stat(paths.pathToActual); os.IsNotExist(statErr) {
			t.Fatalf("File %s was not created", paths.pathToActual)
		}

		if diff := compareTestFileWithActual(paths.pathToTest, paths.pathToActual, fs); diff != "" {
			t.Fatalf(diff)
		}
	}

	pathToDockerCompose := filepath.Join(actualFilesRoot, "docker-compose.yml")

	if _, statErr := fs.Stat(pathToDockerCompose); os.IsNotExist(statErr) {
		t.Fatalf("File %s was not created", pathToDockerCompose)
	}

	dockerComposeContent, readErr := afero.ReadFile(fs, pathToDockerCompose)

	if readErr != nil {
		t.Fatalf("Failed to read file %s: %s", pathToDockerCompose, readErr)
	}

	dockerComposeYaml := map[string]interface{}{}

	unmarshallErr := yaml.Unmarshal(dockerComposeContent, dockerComposeYaml)

	if unmarshallErr != nil {
		t.Fatalf("Failed to unmarshall docker compose: %s", unmarshallErr)
	}

	if diff := cmp.Diff(createTestDockerComposeConf(), dockerComposeYaml); diff != "" {
		t.Fatalf("docker-compose mismatch (-want +got):\n%s", diff)
	}
}

func compareTestFileWithActual(pathToTest, pathToActual string, fs afero.Fs) string {
	testFile, readTestErr := ioutil.ReadFile(pathToTest)

	if readTestErr != nil {
		return fmt.Sprintf("Could not read test file %s. Reason: %s", pathToTest, readTestErr)
	}

	actualFile, readActualErr := afero.ReadFile(fs, pathToActual)

	if readActualErr != nil {
		return fmt.Sprintf("Could not read actual file %s. Reason: %s", actualFile, readActualErr)
	}

	if diff := cmp.Diff(testFile, actualFile); diff != "" {
		return fmt.Sprintf("Test and actual file mismatch (-want +got):\n%s", diff)
	}

	return ""
}
