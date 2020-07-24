package dockercompose

import (
	"fmt"
)

type NetworkDriver string

const (
	Bridge  NetworkDriver = "bridge"
	Host    NetworkDriver = "host"
	Overlay NetworkDriver = "overlay"
	Macvlan NetworkDriver = "macvlan"
	None    NetworkDriver = "none"
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
