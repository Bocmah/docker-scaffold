package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

type renderable struct {
	tmplPath   string
	outputPath string
}

func (r *renderable) render(conf *service.FullConfig) (*Rendered, error) {
	tmpl, parseErr := template.ParseFiles(r.tmplPath)

	if parseErr != nil {
		return nil, fmt.Errorf("parse template: %s", parseErr)
	}

	if ensureDirErr := r.ensureOutputDir(); ensureDirErr != nil {
		return nil, fmt.Errorf("ensure output dir: %s", ensureDirErr)
	}

	file, createFileErr := os.Create(r.outputPath)

	if createFileErr != nil {
		return nil, fmt.Errorf("create output file: %s", createFileErr)
	}

	defer file.Close()

	if executeErr := tmpl.Execute(file, conf); executeErr != nil {
		return nil, fmt.Errorf("execute template: %s", executeErr)
	}

	return &Rendered{Path: r.outputPath}, nil
}

func (r *renderable) ensureOutputDir() error {
	return os.MkdirAll(filepath.Dir(r.outputPath), 0755)
}

type renderableServices map[service.SupportedService]*renderable

func (r renderableServices) render(config *service.FullConfig) (RenderedServices, error) {
	renderedServ := RenderedServices{}

	for s, renderable := range r {
		if !config.Services.IsPresent(s) {
			continue
		}

		rendered, err := renderable.render(config)

		if err != nil {
			return nil, fmt.Errorf("render %s service: %s", s, err)
		}

		renderedServ[s] = rendered
	}

	return renderedServ, nil
}

func (r renderableServices) outputPaths() outputPaths {
	outputPaths := outputPaths{}

	for s, renderable := range r {
		outputPaths[s] = renderable.outputPath
	}

	return outputPaths
}

func newRenderableServices(rootTmplPath, rootOutputPath string) renderableServices {
	return renderableServices{
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
	}
}

type outputPaths map[service.SupportedService]string
