package service

import "fmt"

// NodeJSConfig is a user-defined config for Node.js
type NodeJSConfig struct {
	Version string
}

func (n *NodeJSConfig) FillDefaultsIfNotSet() {
	if n.Version == "" {
		n.Version = "latest"
	}
}

func (n *NodeJSConfig) Validate() error {
	errors := &ValidationErrors{}

	if n.Version == "" {
		errors.Add("Node.js version is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (n *NodeJSConfig) String() string {
	return fmt.Sprintf("NodeJSConfig{Version: %s}", n.Version)
}
