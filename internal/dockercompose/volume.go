package dockercompose

import (
	"fmt"
	"strings"
)

type VolumeDriver string

const (
	VolumeDriverLocal VolumeDriver = "local"
)

type ServiceVolume struct {
	Source string
	Target string
}

func (v *ServiceVolume) String() string {
	if v.Target == "" {
		return ""
	}

	return Mapping(v.Source, v.Target)
}

type ServiceVolumes []*ServiceVolume

func (v ServiceVolumes) Render() string {
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

func (v *NamedVolume) ToServiceVolume() *ServiceVolume {
	if v.Name == "" {
		return nil
	}

	return &ServiceVolume{Target: v.Name}
}

type NamedVolumes []*NamedVolume

func (v NamedVolumes) Render() string {
	var sb strings.Builder

	if len(v) != 0 {
		for _, vol := range v {
			sb.WriteString(vol.Render())
		}
	}

	return sb.String()
}

func (v NamedVolumes) IsEmpty() bool {
	return len(v) == 0
}

func (v NamedVolumes) ToServiceVolumes() ServiceVolumes {
	vols := ServiceVolumes{}

	for _, vol := range v {
		serviceVol := vol.ToServiceVolume()

		if serviceVol != nil {
			vols = append(vols, vol.ToServiceVolume())
		}
	}

	return vols
}
