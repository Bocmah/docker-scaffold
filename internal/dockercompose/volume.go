package dockercompose

import (
	"fmt"
	"strings"
)

type VolumeDriver string

const (
	VolumeDriverLocal VolumeDriver = "local"
)

type Volume struct {
	Source string
	Target string
}

func (v *Volume) String() string {
	if v.Target == "" {
		return ""
	}

	return Mapping(v.Source, v.Target)
}

type Volumes []*Volume

func (v Volumes) Render() string {
	length := len(v)

	if length == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("volumes:\n")

	for i, volume := range v {
		sb.WriteString(fmt.Sprintf("  - %s", volume))

		if i+1 != length {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

type NamedVolume struct {
	Name   string
	Driver VolumeDriver
}

func (v *NamedVolume) Render() string {
	if v.Name == "" || v.Driver == "" {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:", v.Name))

	if v.Driver == VolumeDriverLocal {
		return sb.String()
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  driver: %s", v.Driver))

	return sb.String()
}
