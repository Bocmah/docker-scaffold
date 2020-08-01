package dockercompose

import (
	"fmt"
	"strings"
)

type Config struct {
	Version  string
	Services []*Service
	Networks []*Network
	Volumes  []*NamedVolume
}

func (c *Config) Render() string {
	if c.Version == "" || len(c.Services) == 0 {
		return ""
	}

	nesting := NestingLevel(1)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("version: %s", DoubleQuotted(c.Version)))

	sb.WriteString("\n")
	sb.WriteString("services:")

	for _, s := range c.Services {
		sb.WriteString("\n")
		sb.WriteString(nesting.ApplyTo(s.Render()))
	}

	if len(c.Networks) != 0 {
		sb.WriteString("\nnetworks:")

		for _, n := range c.Networks {
			sb.WriteString("\n")
			sb.WriteString(nesting.ApplyTo(n.Render()))
		}
	}

	if len(c.Volumes) != 0 {
		sb.WriteString("\nvolumes:")

		for _, v := range c.Volumes {
			sb.WriteString("\n")
			sb.WriteString(nesting.ApplyTo(v.Render()))
		}
	}

	return sb.String()
}
