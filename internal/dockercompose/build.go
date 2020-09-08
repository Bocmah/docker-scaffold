package dockercompose

import (
	"fmt"
	"strings"
)

// Build represents 'build' directive in docker-compose file
type Build struct {
	Context    string
	Dockerfile string
}

// Render formats Build as YAML string
func (b *Build) Render() string {
	if b.Context == "" {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("build:")

	if b.Dockerfile == "" {
		sb.WriteString(fmt.Sprintf(" %s", b.Context))

		return sb.String()
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  context: %s\n", b.Context))
	sb.WriteString(fmt.Sprintf("  dockerfile: %s", b.Dockerfile))

	return sb.String()
}
