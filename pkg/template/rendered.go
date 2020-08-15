package template

import (
	"fmt"
	"os"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

type Rendered struct {
	Path        string
	CreatedDirs []string
}

type RenderedServices struct {
	Services    map[service.SupportedService]*Rendered
	CreatedDirs []string
}

func (rs *RenderedServices) DeleteAllCreatedFiles() error {
	for _, created := range rs.CreatedDirs {
		err := os.RemoveAll(created)

		if err != nil {
			return fmt.Errorf("remove all: %s", err)
		}
	}

	rs.Services = map[service.SupportedService]*Rendered{}
	rs.CreatedDirs = []string{}

	return nil
}
