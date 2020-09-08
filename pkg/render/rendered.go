package render

import (
	"fmt"
	"os"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

type Rendered struct {
	Path        string
	CreatedDirs []string
}

func (r *Rendered) deleteCreatedDirs() error {
	for _, created := range r.CreatedDirs {
		err := os.RemoveAll(created)

		if err != nil {
			return fmt.Errorf("remove all: %s", err)
		}
	}

	return nil
}

type RenderedServices struct {
	Services map[service.SupportedService][]*Rendered
}

func (rs *RenderedServices) DeleteAllCreatedFiles() error {
	for _, rendered := range rs.Services {
		for _, r := range rendered {
			err := r.deleteCreatedDirs()

			if err != nil {
				return fmt.Errorf("delete created dirs: %s", err)
			}
		}
	}

	return nil
}
