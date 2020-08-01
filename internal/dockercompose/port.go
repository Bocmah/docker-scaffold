package dockercompose

import (
	"fmt"
	"strconv"
	"strings"
)

type Ports []*PortsMapping

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

type PortsMapping struct {
	Host      int
	Container int
}

func (m *PortsMapping) Render() string {
	if m.Container == 0 {
		return ""
	}

	if m.Host == 0 {
		return DoubleQuotted(strconv.Itoa(m.Container))
	}

	return DoubleQuotted(Mapping(strconv.Itoa(m.Host), strconv.Itoa(m.Container)))
}
