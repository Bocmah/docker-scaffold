package dockercompose

import (
	"fmt"
	"strings"
)

// NetworkDriver is one of the network drivers supported by docker
type NetworkDriver string

const (
	NetworkDriverBridge  NetworkDriver = "bridge"
	NetworkDriverHost    NetworkDriver = "host"
	NetworkDriverOverlay NetworkDriver = "overlay"
	NetworkDriverMacvlan NetworkDriver = "macvlan"
	NetworkDriverNone    NetworkDriver = "none"
)

// Network is a top-level network in docker-compose file
type Network struct {
	Name   string
	Driver NetworkDriver
}

// Render formats Network as YAML string
func (n *Network) Render() string {
	if n.Name == "" || n.Driver == "" {
		return ""
	}

	return fmt.Sprintf("%s:\n  driver: %s", n.Name, n.Driver)
}

// ServiceNetworks is service-level networks
type ServiceNetworks []*Network

// Render formats ServiceNetworks as YAML string
func (n ServiceNetworks) Render() string {
	length := len(n)

	if length == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("networks:\n")

	for i, network := range n {
		sb.WriteString(fmt.Sprintf("  - %s", network.Name))

		if i+1 != length {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// Networks is a top-level networks directive
type Networks []*Network

// Render formats Networks as YAML string
func (n Networks) Render() string {
	var sb strings.Builder

	if len(n) != 0 {
		for _, vol := range n {
			sb.WriteString(vol.Render())
		}
	}

	return sb.String()
}

// IsEmpty checks if Networks has zero networks
func (n Networks) IsEmpty() bool {
	return len(n) == 0
}

// ToServiceNetworks transforms top-level Networks to service-level ServiceNetworks
func (n Networks) ToServiceNetworks() ServiceNetworks {
	return ServiceNetworks(n)
}
