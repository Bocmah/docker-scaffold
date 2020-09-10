package service_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func TestFile_GetTemplatePath(t *testing.T) {
	file := service.File{
		PathOnHost:   "/host/Dockerfile",
		Type:         service.Dockerfile,
		TemplatePath: "/path/to/template/test.tmpl",
	}

	want := "/path/to/template/test.tmpl"
	got := file.GetTemplatePath()

	if want != got {
		t.Fatalf("Incorrect template path: want %s got %s", want, got)
	}
}

func TestFile_GetOutputPath(t *testing.T) {
	file := service.File{
		PathOnHost:   "/host/Dockerfile",
		Type:         service.Dockerfile,
		TemplatePath: "/path/to/template/test.tmpl",
	}

	want := "/host/Dockerfile"
	got := file.GetOutputPath()

	if want != got {
		t.Fatalf("Incorrect output path: want %s got %s", want, got)
	}
}

func TestFile_IsMountable(t *testing.T) {
	tests := map[string]struct {
		input service.File
		want  bool
	}{
		"with path on host and path in container": {
			input: service.File{
				Type:            service.ConfigFile,
				PathOnHost:      "/home/test/conf.txt",
				PathInContainer: "/container/conf.txt",
				TemplatePath:    "/path/to/template/conf.tmpl",
			},
			want: true,
		},
		"with path on host, without path in container": {
			input: service.File{
				Type:         service.Dockerfile,
				PathOnHost:   "/host/Dockerfile",
				TemplatePath: "/path/to/template/test.tmpl",
			},
			want: false,
		},
		"with path in container, but without path on host": {
			input: service.File{
				Type:            service.ConfigFile,
				PathInContainer: "/container/conf.txt",
			},
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.IsMountable()
			if tc.want != got {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
