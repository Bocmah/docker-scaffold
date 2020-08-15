package template

import "github.com/Bocmah/phpdocker-scaffold/pkg/service"

type Rendered struct {
	Path        string
	CreatedDirs []string
}

type RenderedServices struct {
	Services    map[service.SupportedService]*Rendered
	CreatedDirs []string
}
