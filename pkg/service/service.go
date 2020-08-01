package service

import "fmt"

type SupportedService int

func (s SupportedService) String() string {
	if s < PHP || s > Database {
		return "Unknown"
	}

	services := [...]string{
		"PHP",
		"NodeJS",
		"Nginx",
		"Database",
	}

	return services[s-1]
}

const (
	PHP SupportedService = iota + 1
	NodeJS
	Nginx
	Database
)

type ServicesConfig struct {
	PHP      *PHPConfig
	NodeJS   *NodeJSConfig
	Nginx    *NginxConfig
	Database *DatabaseConfig
}

func (s *ServicesConfig) FillDefaultsIfNotSet() {
	s.Database.FillDefaultsIfNotSet()

	s.PHP.AddDatabaseExtension(s.Database.System)

	s.PHP.FillDefaultsIfNotSet()
	s.NodeJS.FillDefaultsIfNotSet()
	s.Nginx.FillDefaultsIfNotSet()
}

func (s *ServicesConfig) Validate() error {
	errors := &ValidationErrors{}

	services := []Config{s.PHP, s.NodeJS, s.Nginx, s.Database}

	for _, s := range services {
		errs := s.Validate()
		if errs != nil {
			if e, ok := errs.(*ValidationErrors); ok {
				errors.Merge(e)
			} else {
				errors.Add(errs.Error())
			}
		}
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (s *ServicesConfig) IsPresent(service SupportedService) bool {
	switch service {
	case PHP:
		return s.PHP != nil && !s.PHP.IsEmpty()
	case NodeJS:
		return s.NodeJS != nil && !(*s.NodeJS == NodeJSConfig{})
	case Nginx:
		return s.Nginx != nil && !(*s.Nginx == NginxConfig{})
	case Database:
		return s.Database != nil && !(*s.Database == DatabaseConfig{})
	default:
		return false
	}
}

func (s *ServicesConfig) String() string {
	return fmt.Sprintf(
		"ServicesConfig{PHPConfig: %v, NodeJSConfig: %v, NginxConfig: %v, DatabaseConfig: %v}",
		s.PHP,
		s.NodeJS,
		s.Nginx,
		s.Database,
	)
}
