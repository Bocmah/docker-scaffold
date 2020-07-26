package dockercompose

import (
	"fmt"
	"strings"
)

type Ports []PortsMapping

func (p Ports) String() string {
	length := len(p)

	if length == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("ports:\n")

	for i, m := range p {
		if i+1 == length {
			_, _ = fmt.Fprintf(&sb, "  - %s", m)
		} else {
			_, _ = fmt.Fprintf(&sb, "  - %s\n", m)
		}
	}

	return sb.String()
}

type PortsMapping struct {
	Host      string
	Container string
}

func (m PortsMapping) String() string {
	if m.Container == "" {
		return ""
	}

	if m.Host == "" {
		return DoubleQuotted(m.Container)
	}

	return DoubleQuotted(Mapping(m.Host, m.Container))
}
