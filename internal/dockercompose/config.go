package dockercompose

import (
	"fmt"
	"strings"
)

// Config represents docker-compose file as a struct
type Config struct {
	Version  string
	Services []*Service
	Networks Networks
	Volumes  NamedVolumes
}

// Render formats Config as YAML string
func (c *Config) Render() string {
	if c.Version == "" || len(c.Services) == 0 {
		return ""
	}

	nesting := NestingLevel(1)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("version: %s", doubleQuotted(c.Version)))

	sb.WriteString("\n")
	sb.WriteString("services:")

	for _, s := range c.Services {
		sb.WriteString("\n")
		sb.WriteString(nesting.ApplyTo(s.Render()))
	}

	if !c.Networks.IsEmpty() {
		sb.WriteString("\nnetworks:\n")
		sb.WriteString(nesting.ApplyTo(c.Networks.Render()))
	}

	if !c.Volumes.IsEmpty() {
		sb.WriteString("\nvolumes:\n")
		sb.WriteString(nesting.ApplyTo(c.Volumes.Render()))
	}

	return sb.String()
}
