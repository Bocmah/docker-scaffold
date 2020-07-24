package service

import "fmt"

type Nginx struct {
	HttpPort   int     `yaml:"httpPort"`
	HttpsPort  int     `yaml:"httpsPort"`
	ServerName string  `yaml:"serverName"`
	FastCGI    FastCGI `yaml:"fastCGI"`
}

type FastCGI struct {
	PassPort           int `yaml:"passPort"`
	ReadTimeoutSeconds int `yaml:"readTimeoutSeconds"`
}

func (n *Nginx) FillDefaultsIfNotSet() {
	if n.HttpPort == 0 {
		n.HttpPort = 80
	}

	if n.HttpsPort == 0 {
		n.HttpsPort = 443
	}

	if n.FastCGI.PassPort == 0 {
		n.FastCGI.PassPort = 9000
	}

	if n.FastCGI.ReadTimeoutSeconds == 0 {
		n.FastCGI.ReadTimeoutSeconds = 60
	}
}

func (n *Nginx) Validate() error {
	errors := &ValidationErrors{}

	if n.HttpPort == 0 {
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

func (n *Nginx) String() string {
	return fmt.Sprintf(
		"Nginx{HttpPort: %d, HttpsPort: %d, ServerName: %s, FastCGI: %v}",
		n.HttpPort,
		n.HttpsPort,
		n.ServerName,
		n.FastCGI,
	)
}
