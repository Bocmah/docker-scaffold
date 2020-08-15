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
	rendered := &Rendered{}

	tmpl, parseErr := template.ParseFiles(r.tmplPath)

	if parseErr != nil {
		return nil, fmt.Errorf("parse template: %s", parseErr)
	}

	if !r.outputPathExists() {
		absPath, makeErr := r.makeOutputPath()

		if makeErr != nil {
			return nil, fmt.Errorf("make root output path: %s", makeErr)
		}

		rendered.CreatedDirs = append(rendered.CreatedDirs, absPath)
	}

	file, createFileErr := os.Create(r.outputPath)

	if createFileErr != nil {
		return nil, fmt.Errorf("create output file: %s", createFileErr)
	}

	defer file.Close()

	if executeErr := tmpl.Execute(file, conf); executeErr != nil {
		return nil, fmt.Errorf("execute template: %s", executeErr)
	}

	rendered.Path = r.outputPath

	return rendered, nil
}

func (r *renderable) outputPathExists() bool {
	if _, statErr := os.Stat(filepath.Dir(r.outputPath)); os.IsNotExist(statErr) {
		return false
	}

	return true
}

func (r *renderable) makeOutputPath() (absPath string, err error) {
	outputDir := filepath.Dir(r.outputPath)

	if mkdirErr := os.MkdirAll(outputDir, 0755); mkdirErr != nil {
		return "", fmt.Errorf("MkdirAll: %s", mkdirErr)
	}

	absPath, absErr := filepath.Abs(outputDir)

	if absErr != nil {
		return "", fmt.Errorf("absolute path: %s", absErr)
	}

	return absPath, nil
}
