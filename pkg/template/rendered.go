package template

import "github.com/Bocmah/phpdocker-scaffold/pkg/service"

type Rendered struct {
	Path string
}

type RenderedServices map[service.SupportedService]*Rendered
