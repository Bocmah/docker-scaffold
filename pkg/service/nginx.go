package service

import "fmt"

type Nginx struct {
	HttpPort           int    `yaml:"httpPort"`
	HttpsPort          int    `yaml:"httpsPort"`
	ServerName         string `yaml:"serverName"`
	FastCGIPassPort    int    `yaml:"fastCGIPassPort"`
	FastCGIReadTimeout int    `yaml:"fastCGIReadTimeout"`
}

func (n *Nginx) FillDefaultsIfNotSet() {
	if n.HttpPort == 0 {
		n.HttpPort = 80
	}

	if n.HttpsPort == 0 {
		n.HttpsPort = 443
	}

	if n.FastCGIPassPort == 0 {
		n.FastCGIPassPort = 9000
	}

	if n.FastCGIReadTimeout == 0 {
		n.FastCGIReadTimeout = 60
	}
}

func (n *Nginx) Validate() error {
	errors := &ValidationErrors{}

	if n.HttpPort == 0 {
		errors.Add("nginx port is required")
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

func (n *Nginx) String() string {
	return fmt.Sprintf(
		"Nginx{HttpPort: %d, HttpsPort: %d, ServerName: %s, FastCGIPassPort: %d, FastCGIReadTimeout: %d}",
		n.HttpPort,
		n.HttpsPort,
		n.ServerName,
		n.FastCGIPassPort,
		n.FastCGIReadTimeout,
	)
}
