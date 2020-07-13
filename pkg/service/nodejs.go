package service

type NodeJS struct {
	Version string
}

func (n *NodeJS) FillDefaultsIfNotSet() {
	if n.Version == "" {
		n.Version = "latest"
	}
}

func (n *NodeJS) Validate() *ValidationErrors {
	errors := &ValidationErrors{}

	if n.Version == "" {
		errors.Add("Node.js version is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

