package render_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/pkg/render"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func TestRenderTemplatesFromConfiguration(t *testing.T) {
	absPathToOutput, err := filepath.Abs("testdata")

	if err != nil {
		t.Fatalf("failed to find absolute path to testdata dir: %s", err)
	}

	dir, tempDirErr := ioutil.TempDir(absPathToOutput, "output")

	if tempDirErr != nil {
		t.Fatalf("failed to create temp dir: %s", tempDirErr)
	}

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

	rendered, renderErr := render.RenderServices(conf)

	if renderErr != nil {
		t.Fatalf("Encountered non-nil error in correct test case: %v", renderErr)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output path %v still doesn't exist after calling rendering function", outputPath)
	}

	absPath, absErr := filepath.Abs(outputPath + string(os.PathSeparator))

	if absErr != nil {
		t.Fatal(absErr)
	}

	wantRendered := &render.RenderedServices{
		Services: map[service.SupportedService][]*render.Rendered{
			service.PHP: {
				{
					Path:        filepath.Join(outputPath, "php/Dockerfile"),
					CreatedDirs: []string{filepath.Join(absPath, "php")},
				},
			},
			service.Nginx: {
				{
					Path:        filepath.Join(outputPath, "nginx/conf.d/app.conf"),
					CreatedDirs: []string{filepath.Join(absPath, "nginx/conf.d")},
				},
			},
			service.NodeJS: {
				{
					Path:        filepath.Join(outputPath, "nodejs/Dockerfile"),
					CreatedDirs: []string{filepath.Join(absPath, "nodejs")},
				},
			},
		},
	}

	if diff := cmp.Diff(wantRendered, rendered); diff != "" {
		t.Fatalf("RenderedServices mismatch (-want +got):\n%s", diff)
	}

	testFilesRoot := "testdata/template_render/.docker"
	testFiles := map[service.SupportedService][]string{
		service.PHP:    {filepath.Join(testFilesRoot, "php/Dockerfile")},
		service.Nginx:  {filepath.Join(testFilesRoot, "nginx/conf.d/app.conf")},
		service.NodeJS: {filepath.Join(testFilesRoot, "nodejs/Dockerfile")},
	}

	if diff := compareRenderedWithExpected(rendered, testFiles); diff != "" {
		t.Fatalf(diff)
	}
}

func TestRenderDockerCompose(t *testing.T) {
	absPathToOutput, absErr := filepath.Abs("testdata")

	if absErr != nil {
		t.Fatalf("failed to find absolute path to testdata dir: %s", absErr)
	}

	dir, tempDirErr := ioutil.TempDir(absPathToOutput, "output")

	if tempDirErr != nil {
		t.Fatalf("failed to create temp dir: %s", tempDirErr)
	}

	defer os.RemoveAll(dir)

	network := &dockercompose.Network{Name: "awesome-app-network", Driver: dockercompose.NetworkDriverBridge}
	namedVolume := &dockercompose.NamedVolume{Name: "awesome-app-data", Driver: dockercompose.VolumeDriverLocal}
	projectRoot := "/home/test/app"
	workDir := "/var/www"

	conf := &dockercompose.Config{
		Version:  "3.8",
		Networks: dockercompose.Networks{network},
		Volumes:  dockercompose.NamedVolumes{namedVolume},
		Services: []*dockercompose.Service{
			{
				Name:          "php",
				ContainerName: "php",
				Build: &dockercompose.Build{
					Context:    projectRoot,
					Dockerfile: filepath.Join(projectRoot, ".docker/php/Dockerfile"),
				},
				Image: &dockercompose.Image{
					Name: "awesome-app",
				},
				WorkingDir: workDir,
				Restart:    dockercompose.RestartPolicyUnlessStopped,
				Networks:   dockercompose.ServiceNetworks{network},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: projectRoot, Target: workDir},
				},
			},
			{
				Name:          "webserver",
				ContainerName: "webserver",
				Image:         &dockercompose.Image{Name: "nginx", Tag: "alpine"},
				Ports: dockercompose.Ports{
					&dockercompose.PortsMapping{Host: 80, Container: 80},
					&dockercompose.PortsMapping{Host: 443, Container: 443},
				},
				Networks: dockercompose.ServiceNetworks{network},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: projectRoot, Target: workDir},
					&dockercompose.ServiceVolume{Source: "./nginx/conf.d/", Target: "/etc/nginx/conf.d/"},
				},
			},
			{
				Name:          "db",
				ContainerName: "db",
				Image:         &dockercompose.Image{Name: "mysql", Tag: "8.0"},
				Restart:       dockercompose.RestartPolicyUnlessStopped,
				Ports: dockercompose.Ports{
					&dockercompose.PortsMapping{Host: 3306, Container: 3306},
				},
				Environment: dockercompose.Environment{
					"MYSQL_DATABASE":      "test-db",
					"MYSQL_ROOT_PASSWORD": "secret-root",
					"MYSQL_USER":          "test-user",
					"MYSQL_PASSWORD":      "secret-password",
				},
				Networks: dockercompose.ServiceNetworks{network},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: namedVolume.Name, Target: "/var/lib/mysql"},
				},
			},
			{
				Name:          "nodejs",
				ContainerName: "nodejs",
				Build: &dockercompose.Build{
					Context:    projectRoot,
					Dockerfile: filepath.Join(projectRoot, ".docker/nodejs/Dockerfile"),
				},
				WorkingDir: "/opt",
				Networks:   dockercompose.ServiceNetworks{network},
				Volumes: dockercompose.ServiceVolumes{
					&dockercompose.ServiceVolume{Source: projectRoot, Target: "/opt"},
				},
			},
		},
	}

	outputPath := filepath.Join(dir, "docker-compose.yml")

	renderErr := render.RenderDockerCompose(conf, outputPath)

	if renderErr != nil {
		t.Fatalf("encountered non nil err with correct configuration: %s", renderErr)
	}

	data, readErr := ioutil.ReadFile(outputPath)

	if readErr != nil {
		t.Fatalf("failed to read resulting docker-compose.yml: %s", readErr)
	}

	want := map[string]interface{}{
		"version": "3.8",
		"services": map[interface{}]interface{}{
			"php": map[interface{}]interface{}{
				"container_name": "php",
				"working_dir":    workDir,
				"build": map[interface{}]interface{}{
					"context":    projectRoot,
					"dockerfile": filepath.Join(projectRoot, ".docker/php/Dockerfile"),
				},
				"image":   "awesome-app",
				"restart": string(dockercompose.RestartPolicyUnlessStopped),
				"networks": []interface{}{
					network.Name,
				},
				"volumes": []interface{}{
					projectRoot + ":" + workDir,
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
					network.Name,
				},
				"volumes": []interface{}{
					projectRoot + ":" + workDir,
					"./nginx/conf.d/:/etc/nginx/conf.d/",
				},
			},
			"db": map[interface{}]interface{}{
				"container_name": "db",
				"image":          "mysql:8.0",
				"restart":        string(dockercompose.RestartPolicyUnlessStopped),
				"ports": []interface{}{
					"3306:3306",
				},
				"environment": map[interface{}]interface{}{
					"MYSQL_DATABASE":      "test-db",
					"MYSQL_ROOT_PASSWORD": "secret-root",
					"MYSQL_USER":          "test-user",
					"MYSQL_PASSWORD":      "secret-password",
				},
				"networks": []interface{}{
					network.Name,
				},
				"volumes": []interface{}{
					namedVolume.Name + ":" + "/var/lib/mysql",
				},
			},
			"nodejs": map[interface{}]interface{}{
				"container_name": "nodejs",
				"working_dir":    "/opt",
				"build": map[interface{}]interface{}{
					"context":    projectRoot,
					"dockerfile": filepath.Join(projectRoot, ".docker/nodejs/Dockerfile"),
				},
				"networks": []interface{}{
					network.Name,
				},
				"volumes": []interface{}{
					projectRoot + ":" + "/opt",
				},
			},
		},
		"networks": map[interface{}]interface{}{
			network.Name: map[interface{}]interface{}{
				"driver": string(network.Driver),
			},
		},
		"volumes": map[interface{}]interface{}{
			namedVolume.Name: nil,
		},
	}

	got := map[string]interface{}{}

	if unmarshallErr := yaml.Unmarshal(data, got); unmarshallErr != nil {
		t.Fatalf("failed to unmarshall docker-compose.yml: %s", unmarshallErr)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Expected and rendered dockercompose mismatch (-want +got):\n%s", diff)
	}
}

func compareRenderedWithExpected(renderedServices *render.RenderedServices, testFiles map[service.SupportedService][]string) (diff string) {
	for serv, files := range testFiles {
		renderedService, ok := renderedServices.Services[serv]

		if !ok {
			return fmt.Sprintf("Service %s was not rendered", serv)
		}

		for index, file := range files {
			expectedFile, readExpectedErr := ioutil.ReadFile(file)

			if readExpectedErr != nil {
				return fmt.Sprintf("Could not read expected file for service %s at path %s. Reason: %s", serv, file, readExpectedErr)
			}

			renderedFile, readRenderedErr := ioutil.ReadFile(renderedService[index].Path)

			if readRenderedErr != nil {
				return fmt.Sprintf("Could not read expected file for service %s at path %s. Reason: %s", serv, renderedService[index].Path, readRenderedErr)
			}

			if diff := cmp.Diff(expectedFile, renderedFile); diff != "" {
				return fmt.Sprintf("Expected and rendered file for service %s mismatch (-want +got):\n%s", serv, diff)
			}
		}
	}

	return ""
}
