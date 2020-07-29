package dockercompose

import (
	"bufio"
	"fmt"
	"strings"
)

type NestingLevel int

func (n NestingLevel) ApplyTo(str string) string {
	if n <= 0 {
		return str
	}

	spaces := strings.Repeat(" ", int(n)*2)
	scanner := bufio.NewScanner(strings.NewReader(str))

	var lines []string

	for scanner.Scan() {
		lines = append(lines, spaces+scanner.Text())
	}

	return strings.Join(lines, "\n")
}

type renderableItem interface {
	Render() string
}

type Service struct {
	Name          string
	Build         Build
	Image         Image
	ContainerName string
	WorkingDir    string
	Restart       RestartPolicy
	Ports         Ports
	Environment   Environment
	Networks      Networks
	Volumes       Volumes
}

func (s Service) Render() string {
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

	renderables := []renderableItem{s.Build, s.Image, s.Restart, s.Ports, s.Environment, s.Networks, s.Volumes}

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
