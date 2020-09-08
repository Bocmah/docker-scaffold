package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

type RenderableFile interface {
	GetTemplatePath() string
	GetOutputPath() string
}

func render(renderable RenderableFile, conf *service.FullConfig) (*Rendered, error) {
	rendered := &Rendered{}

	tmpl, parseErr := template.ParseFiles(renderable.GetTemplatePath())

	if parseErr != nil {
		return nil, fmt.Errorf("parse render: %s", parseErr)
	}

	absOutputPath, absErr := filepath.Abs(renderable.GetOutputPath())

	if absErr != nil {
		return nil, fmt.Errorf("absolute path to created output dir: %s", absErr)
	}

	outputDir := filepath.Dir(absOutputPath)

	if !pathExists(outputDir) {
		if mkdirErr := os.MkdirAll(outputDir, 0755); mkdirErr != nil {
			return nil, fmt.Errorf("MkdirAll: %s", mkdirErr)
		}

		rendered.CreatedDirs = append(rendered.CreatedDirs, outputDir)
	}

	file, createFileErr := os.Create(renderable.GetOutputPath())

	if createFileErr != nil {
		return nil, fmt.Errorf("create output file: %s", createFileErr)
	}

	defer file.Close()

	if executeErr := tmpl.Execute(file, conf); executeErr != nil {
		return nil, fmt.Errorf("execute render: %s", executeErr)
	}

	rendered.Path = absOutputPath

	return rendered, nil
}

func pathExists(path string) bool {
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return false
	}

	return true
}
