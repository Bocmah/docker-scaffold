package dockercompose

import (
	"fmt"
	"strconv"
	"strings"
)

// Ports represents 'ports' directive in docker-compose file
type Ports []*PortsMapping

// Render formats Ports as YAML string
func (p Ports) Render() string {
	length := len(p)

	if length == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("ports:\n")

	for i, m := range p {
		if i+1 == length {
			_, _ = fmt.Fprintf(&sb, "  - %s", m.Render())
		} else {
			_, _ = fmt.Fprintf(&sb, "  - %s\n", m.Render())
		}
	}

	return sb.String()
}

// PortsMapping represents a single mapping of host port to container port
type PortsMapping struct {
	Host      int
	Container int
}

// Render formats PortsMapping as YAML string
func (m *PortsMapping) Render() string {
	if m.Container == 0 {
		return ""
	}

	if m.Host == 0 {
		return doubleQuotted(strconv.Itoa(m.Container))
	}

	return doubleQuotted(mapping(strconv.Itoa(m.Host), strconv.Itoa(m.Container)))
}
