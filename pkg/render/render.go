package render

import (
	"fmt"
	"os"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

// RenderServices renders files for all services from service.FullConfig
func RenderServices(conf *service.FullConfig) (*RenderedServices, error) {
	renderedServices := RenderedServices{
		Services: map[service.SupportedService][]*Rendered{},
	}

	for serv, renderableFiles := range conf.GetServiceFiles() {
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

// RenderDockerCompose renders docker-compose.yml file
func RenderDockerCompose(conf *dockercompose.Config, outputPath string) error {
	file, createErr := os.Create(outputPath)

	if createErr != nil {
		return fmt.Errorf("create output file: %s", createErr)
	}

	defer file.Close()

	_, writeErr := file.WriteString(conf.Render())

	if writeErr != nil {
		return fmt.Errorf("write to output file: %s", writeErr)
	}

	return nil
}
