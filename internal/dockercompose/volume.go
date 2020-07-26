package dockercompose

import (
	"fmt"
	"strings"
)

type Volume struct {
	Source string
	Target string
}

func (v Volume) String() string {
	if v.Target == "" {
		return ""
	}

	return mapping(v.Source, v.Target)
}

type NamedVolume struct {
	Name   string
	Driver string
}

func (v NamedVolume) String() string {
	if v.Name == "" || v.Driver == "" {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:", v.Name))

	if v.Driver == "local" {
		return sb.String()
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  driver: %s", v.Driver))

	return sb.String()
}
