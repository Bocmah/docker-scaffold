package render

import (
	"fmt"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

// Rendered is a rendered file
type Rendered struct {
	Path        string
	CreatedDirs []string
}

func (r *Rendered) deleteCreatedDirs() error {
	for _, created := range r.CreatedDirs {
		err := AppFs.RemoveAll(created)

		if err != nil {
			return fmt.Errorf("remove all: %s", err)
		}
	}

	return nil
}

// RenderedServices is per-service rendered files
type RenderedServices struct {
	Services map[service.SupportedService][]*Rendered
}

// DeleteAllCreatedFiles deletes all created files for all services
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
