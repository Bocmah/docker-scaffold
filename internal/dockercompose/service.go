package dockercompose

import (
	"fmt"
	"strings"
)

type renderableItem interface {
	Render() string
}

type Service struct {
	Name          string
	Build         *Build
	Image         *Image
	ContainerName string
	WorkingDir    string
	Restart       RestartPolicy
	Ports         Ports
	Environment   Environment
	Networks      ServiceNetworks
	Volumes       ServiceVolumes
}

func (s *Service) Render() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:", s.Name))

	nesting := NestingLevel(1)

	if s.ContainerName != "" {
		sb.WriteString("\n")
		sb.WriteString(nesting.ApplyTo(fmt.Sprintf("container_name: %s", s.ContainerName)))
	}

	if s.WorkingDir != "" {
		sb.WriteString("\n")
		sb.WriteString(nesting.ApplyTo(fmt.Sprintf("working_dir: %s", s.WorkingDir)))
	}

	var renderables []renderableItem

	if s.Build != nil {
		renderables = append(renderables, s.Build)
	}

	if s.Image != nil {
		renderables = append(renderables, s.Image)
	}

	renderables = append(renderables, s.Restart, s.Ports, s.Environment, s.Networks, s.Volumes)

	for _, r := range renderables {
		rendered := r.Render()

		if rendered != "" {
			rendered = nesting.ApplyTo(rendered)

			sb.WriteString("\n")
			sb.WriteString(rendered)
		}
	}

	return sb.String()
}
