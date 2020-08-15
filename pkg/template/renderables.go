package template

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

type renderableServices struct {
	services       map[service.SupportedService]*renderable
	rootTmplPath   string
	rootOutputPath string
}

func (r renderableServices) rootOutputPathExists() bool {
	if _, statErr := os.Stat(r.rootOutputPath); os.IsNotExist(statErr) {
		return false
	}

	return true
}

func (r renderableServices) makeRootOutputPath() (absPath string, err error) {
	if mkdirErr := os.MkdirAll(r.rootOutputPath, 0755); mkdirErr != nil {
		return "", fmt.Errorf("MkdirAll: %s", mkdirErr)
	}

	absPath, absErr := filepath.Abs(r.rootOutputPath)

	if absErr != nil {
		return "", fmt.Errorf("absolute path: %s", absErr)
	}

	return absPath, nil
}

func (r renderableServices) render(config *service.FullConfig) (*RenderedServices, error) {
	var createdDirs []string

	if !r.rootOutputPathExists() {
		absPath, makeRootErr := r.makeRootOutputPath()

		if makeRootErr != nil {
			return nil, fmt.Errorf("make root output path: %s", makeRootErr)
		}

		createdDirs = append(createdDirs, absPath)
	}

	rendered, renderErr := r.renderServices(config)

	if renderErr != nil {
		return nil, fmt.Errorf("render services: %s", renderErr)
	}

	rendered.CreatedDirs = append(rendered.CreatedDirs, createdDirs...)

	return rendered, nil
}

func (r renderableServices) renderServices(config *service.FullConfig) (*RenderedServices, error) {
	renderedServices := RenderedServices{
		Services: map[service.SupportedService]*Rendered{},
	}

	for s, renderable := range r.services {
		if !config.Services.IsPresent(s) {
			continue
		}

		rendered, err := renderable.render(config)

		if err != nil {
			return nil, fmt.Errorf("render %s service: %s", s, err)
		}

		renderedServices.Services[s] = rendered
	}

	return &renderedServices, nil
}

func (r renderableServices) outputPaths() outputPaths {
	outputPaths := outputPaths{}

	for s, renderable := range r.services {
		outputPaths[s] = renderable.outputPath
	}

	return outputPaths
}

func newRenderableServices(rootTmplPath, rootOutputPath string) renderableServices {
	return renderableServices{
		services: map[service.SupportedService]*renderable{
			service.PHP: {
				tmplPath:   fullPath(rootTmplPath, "php/php.dockerfile.gotmpl"),
				outputPath: fullPath(rootOutputPath, "php/Dockerfile"),
			},
			service.Nginx: {
				tmplPath:   fullPath(rootTmplPath, "nginx/conf.gotmpl"),
				outputPath: fullPath(rootOutputPath, "nginx/conf.d/app.conf"),
			},
			service.NodeJS: {
				tmplPath:   fullPath(rootTmplPath, "nodejs/nodejs.dockerfile.gotmpl"),
				outputPath: fullPath(rootOutputPath, "nodejs/Dockerfile"),
			},
		},
		rootTmplPath:   rootTmplPath,
		rootOutputPath: rootOutputPath,
	}
}

type outputPaths map[service.SupportedService]string
