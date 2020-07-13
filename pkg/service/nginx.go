package service

type Nginx struct {
	Port               int
	ServerName         string `yaml:"serverName"`
	FastCGIPassPort    int    `yaml:"fastCGIPassPort"`
	FastCGIReadTimeout int    `yaml:"fastCGIReadTimeout"`
}

func (n *Nginx) FillDefaultsIfNotSet() {
	if n.Port == 0 {
		n.Port = 80
	}

	if n.FastCGIPassPort == 0 {
		n.FastCGIPassPort = 9000
	}

	if n.FastCGIReadTimeout == 0 {
		n.FastCGIReadTimeout = 60
	}
}

func (n *Nginx) Validate() *ValidationErrors {
	errors := &ValidationErrors{}

	if n.Port == 0 {
		errors.Add("nginx port is required")
	}

	if n.ServerName == "" {
		errors.Add("nginx server name is required")
	}

	if n.FastCGIPassPort == 0 {
		errors.Add("nginx FastCGI pass port is required")
	}

	if n.FastCGIReadTimeout == 0 {
		errors.Add("nginx FastCGI read timeout is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}
