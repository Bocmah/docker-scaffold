package render

import (
	"fmt"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

func Render(servicesWithFiles map[service.SupportedService][]RenderableFile, conf *service.FullConfig) (*RenderedServices, error) {
	renderedServices := RenderedServices{
		Services: map[service.SupportedService][]*Rendered{},
	}

	for serv, renderableFiles := range servicesWithFiles {
		for _, file := range renderableFiles {
			rendered, renderErr := render(file, conf)

			if renderErr != nil {
				return nil, fmt.Errorf("render services: %s", renderErr)
			}

			renderedServices.Services[serv] = append(renderedServices.Services[serv], rendered)
		}
	}

	return &renderedServices, nil
}
