package service

import "fmt"

type NodeJS struct {
	Version string
}

func (n *NodeJS) FillDefaultsIfNotSet() {
	if n.Version == "" {
		n.Version = "latest"
	}
}

func (n *NodeJS) Validate() error {
	errors := &ValidationErrors{}

	if n.Version == "" {
		errors.Add("Node.js version is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (n *NodeJS) String() string {
	return fmt.Sprintf("NodeJS{Version: %s}", n.Version)
}

