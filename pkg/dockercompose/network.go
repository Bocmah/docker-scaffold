package dockercompose

import (
	"fmt"
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

func (n Network) String() string {
	if n.Name == "" || n.Driver == "" {
		return ""
	}

	return fmt.Sprintf("%s:\n  driver: %s", n.Name, n.Driver)
}
