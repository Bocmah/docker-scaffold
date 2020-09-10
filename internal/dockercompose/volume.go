package dockercompose

import (
	"fmt"
	"strings"
)

// VolumeDriver is one of the volume drivers supported by docker
type VolumeDriver string

// All supported volume drivers
const (
	VolumeDriverLocal VolumeDriver = "local"
)

// ServiceVolume represents service-level volume mapping in docker-compose file
type ServiceVolume struct {
	Source string
	Target string
}

// String formats ServiceVolume as a mapping
func (v *ServiceVolume) String() string {
	if v.Target == "" {
		return ""
	}

	return mapping(v.Source, v.Target)
}

// ServiceVolumes represents service-level volumes directive
type ServiceVolumes []*ServiceVolume

// Render formats ServiceVolumes as YAML string
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

// NamedVolume represents top-level volume in docker-compose file
type NamedVolume struct {
	Name   string
	Driver VolumeDriver
}

// Render formats NamedVolume as YAML string
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

// ToServiceVolume transforms NamedVolume to ServiceVolume
func (v *NamedVolume) ToServiceVolume() *ServiceVolume {
	if v.Name == "" {
		return nil
	}

	return &ServiceVolume{Target: v.Name}
}

// NamedVolumes represents top-level 'volumes' directive in docker-compose file
type NamedVolumes []*NamedVolume

// Render formats NamedVolumes as YAML string
func (v NamedVolumes) Render() string {
	var sb strings.Builder

	if len(v) != 0 {
		for _, vol := range v {
			sb.WriteString(vol.Render())
		}
	}

	return sb.String()
}

// IsEmpty checks if NamedVolumes has zero volumes
func (v NamedVolumes) IsEmpty() bool {
	return len(v) == 0
}

// ToServiceVolumes transforms NamedVolumes to ServiceVolumes
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
