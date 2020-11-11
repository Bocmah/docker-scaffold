package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Bocmah/phpdocker-gen/internal/box"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

// RenderableFile is a file which can be transformed from template to resulting file
type RenderableFile interface {
	// GetTemplatePath returns path to template from which file will be rendered
	GetTemplatePath() string
	// GetOutputPath returns a path to which resulting file should be rendered
	GetOutputPath() string
}

func render(renderable RenderableFile, conf *service.FullConfig) (*Rendered, error) {
	rendered := &Rendered{}

	tmpl := string(box.Get(renderable.GetTemplatePath()))

	parsedTmpl := template.Must(template.New("").Parse(tmpl))

	outputDir := filepath.Dir(renderable.GetOutputPath())

	if !pathExists(outputDir) {
		if mkdirErr := AppFs.MkdirAll(outputDir, 0755); mkdirErr != nil {
			return nil, fmt.Errorf("MkdirAll: %s", mkdirErr)
		}

		rendered.CreatedDirs = append(rendered.CreatedDirs, outputDir)
	}

	file, createFileErr := AppFs.Create(renderable.GetOutputPath())

	if createFileErr != nil {
		return nil, fmt.Errorf("create output file: %s", createFileErr)
	}

	defer file.Close()

	if executeErr := parsedTmpl.Execute(file, conf); executeErr != nil {
		return nil, fmt.Errorf("execute render: %s", executeErr)
	}

	rendered.Path = renderable.GetOutputPath()

	return rendered, nil
}

func pathExists(path string) bool {
	if _, statErr := AppFs.Stat(path); os.IsNotExist(statErr) {
		return false
	}

	return true
}
