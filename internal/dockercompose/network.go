package dockercompose

import (
	"fmt"
	"strings"
)

type NetworkDriver string

const (
	NetworkDriverBridge  NetworkDriver = "bridge"
	NetworkDriverHost    NetworkDriver = "host"
	NetworkDriverOverlay NetworkDriver = "overlay"
	NetworkDriverMacvlan NetworkDriver = "macvlan"
	NetworkDriverNone    NetworkDriver = "none"
)

type Network struct {
	Name   string
	Driver NetworkDriver
}

func (n *Network) Render() string {
	if n.Name == "" || n.Driver == "" {
		return ""
	}

	return fmt.Sprintf("%s:\n  driver: %s", n.Name, n.Driver)
}

type ServiceNetworks []*Network

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

type Networks []*Network

func (n Networks) Render() string {
	var sb strings.Builder

	if len(n) != 0 {
		for _, vol := range n {
			sb.WriteString(vol.Render())
		}
	}

	return sb.String()
}

func (n Networks) IsEmpty() bool {
	return len(n) == 0
}
