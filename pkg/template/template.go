package template

import (
	"fmt"
	"os"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

func RenderTemplatesFromConfiguration(conf *service.FullConfig) (*RenderedServices, error) {
	outputPath := getOutputPath(conf)
	renderables := newRenderableServices("../../tmpl", outputPath)

	rendered, err := renderables.render(conf)

	if err != nil {
		return nil, fmt.Errorf("render renderable services: %s", err)
	}

	return rendered, nil
}

func getOutputPath(config *service.FullConfig) string {
	if config.OutputPath == "" {
		return config.ProjectRoot
	}

	return config.OutputPath
}

func fullPath(rootPath, fromRootPath string) string {
	return rootPath + string(os.PathSeparator) + fromRootPath
}
