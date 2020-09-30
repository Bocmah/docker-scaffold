package render_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/Bocmah/phpdocker-gen/pkg/render"
	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func createDir(t *testing.T, fs afero.Fs, path string) {
	t.Helper()

	mkdirAllErr := fs.MkdirAll(path, 0755)

	if mkdirAllErr != nil {
		t.Fatalf("Failed to create dir %s, err %s", path, mkdirAllErr)
	}
}

func createFile(t *testing.T, fs afero.Fs, path string) {
	t.Helper()

	writeFileErr := afero.WriteFile(fs, path, []byte("test"), 0644)

	if writeFileErr != nil {
		t.Fatalf("Failed to write file %s, err %s", path, writeFileErr)
	}
}

func checkFileDoesntExist(t *testing.T, fs afero.Fs, path string) {
	t.Helper()

	if _, err := fs.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("File %s exists when it shouldn't", path)
	}
}

func TestRenderedServices_DeleteAllCreatedFiles(t *testing.T) {
	render.AppFs = afero.NewMemMapFs()

	testFiles := map[service.SupportedService]struct {
		dirs     []string
		finalDir string
		file     string
	}{
		service.PHP: {
			dirs:     []string{"output", "output/.docker", "output/.docker/php"},
			finalDir: "output/.docker/php",
			file:     "output/.docker/php/Dockerfile",
		},
		service.Nginx: {
			dirs:     []string{"output", "output/.docker", "output/.docker/nginx", "output/.docker/nginx/conf.d"},
			finalDir: "output/.docker/nginx/conf.d",
			file:     "output/.docker/nginx/conf.d/app.conf",
		},
		service.NodeJS: {
			dirs:     []string{"another_output", "another_output/.docker", "another_output/.docker/nginx", "another_output/.docker/nginx/conf.d"},
			finalDir: "another_output/.docker/nginx/conf.d",
			file:     "another_output/.docker/nginx/conf.d/app.conf",
		},
	}

	for _, files := range testFiles {
		createDir(t, render.AppFs, files.dirs[len(files.dirs)-1])
		createFile(t, render.AppFs, files.file)
	}

	rendered := render.RenderedServices{
		Services: map[service.SupportedService][]*render.Rendered{
			service.PHP: {
				&render.Rendered{
					Path:        filepath.Join(testFiles[service.PHP].finalDir, testFiles[service.PHP].file),
					CreatedDirs: testFiles[service.PHP].dirs,
				},
			},
			service.Nginx: {
				&render.Rendered{
					Path:        filepath.Join(testFiles[service.Nginx].finalDir, testFiles[service.Nginx].file),
					CreatedDirs: testFiles[service.Nginx].dirs,
				},
			},
			service.NodeJS: {
				&render.Rendered{
					Path:        filepath.Join(testFiles[service.NodeJS].finalDir, testFiles[service.NodeJS].file),
					CreatedDirs: testFiles[service.NodeJS].dirs,
				},
			},
		},
	}

	deleteErr := rendered.DeleteAllCreatedFiles()

	if deleteErr != nil {
		t.Fatalf("Encountered err while deleting files: %s", deleteErr)
	}

	for _, files := range testFiles {
		for _, dir := range files.dirs {
			checkFileDoesntExist(t, render.AppFs, dir)

		}

		checkFileDoesntExist(t, render.AppFs, files.file)
	}
}
