package service

import "fmt"

type NginxConfig struct {
	HTTPPort   int     `yaml:"httpPort"`
	HTTPSPort  int     `yaml:"httpsPort"`
	ServerName string  `yaml:"serverName"`
	FastCGI    FastCGI `yaml:"fastCGI"`
}

type FastCGI struct {
	PassPort           int `yaml:"passPort"`
	ReadTimeoutSeconds int `yaml:"readTimeoutSeconds"`
}

func (n *NginxConfig) FillDefaultsIfNotSet() {
	if n.HTTPPort == 0 {
		n.HTTPPort = 80
	}

	if n.HTTPSPort == 0 {
		n.HTTPSPort = 443
	}

	if n.FastCGI.PassPort == 0 {
		n.FastCGI.PassPort = 9000
	}

	if n.FastCGI.ReadTimeoutSeconds == 0 {
		n.FastCGI.ReadTimeoutSeconds = 60
	}
}

func (n *NginxConfig) Validate() error {
	errors := &ValidationErrors{}

	if n.HTTPPort == 0 {
		errors.Add("nginx port is required")
	}

	if n.FastCGI.PassPort == 0 {
		errors.Add("nginx FastCGI pass port is required")
	}

	if n.FastCGI.ReadTimeoutSeconds == 0 {
		errors.Add("nginx FastCGI read timeout is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (n *NginxConfig) String() string {
	return fmt.Sprintf(
		"NginxConfig{HTTPPort: %d, HTTPSPort: %d, ServerName: %s, FastCGI: %v}",
		n.HTTPPort,
		n.HTTPSPort,
		n.ServerName,
		n.FastCGI,
	)
}
