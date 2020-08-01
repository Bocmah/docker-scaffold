package service

import "fmt"

type PHPConfig struct {
	Version    string
	Extensions []string
}

func (p *PHPConfig) FillDefaultsIfNotSet() {
	if p.Version == "" {
		p.Version = "7.4"
	}

	if len(p.Extensions) == 0 {
		p.Extensions = []string{"mbstring", "zip", "exif", "pcntl", "gd"}
	}
}

func (p *PHPConfig) AddDatabaseExtension(db SupportedSystem) {
	switch db {
	case MySQL:
		p.Extensions = append(p.Extensions, "pdo_mysql")
	case PostgreSQL:
		p.Extensions = append(p.Extensions, "pdo_pgsql")
	}
}

func (p *PHPConfig) Validate() error {
	errors := &ValidationErrors{}

	if p.Version == "" {
		errors.Add("PHPConfig version is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (p *PHPConfig) String() string {
	return fmt.Sprintf("PHPConfig{Version: %s, Extensions: %v}", p.Version, p.Extensions)
}

func (p *PHPConfig) IsEmpty() bool {
	return p.Version == "" && len(p.Extensions) == 0
}
