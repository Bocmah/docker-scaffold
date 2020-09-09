package service

import "fmt"

// SupportedService is one of the services supported by the tool
type SupportedService int

func (s SupportedService) String() string {
	if s < PHP || s > NodeJS {
		return "Unknown"
	}

	services := [...]string{
		"PHP",
		"Nginx",
		"Database",
		"NodeJS",
	}

	return services[s-1]
}

// SupportedServices returns a collection of all supported services
func SupportedServices() []SupportedService {
	return []SupportedService{
		PHP,
		Nginx,
		Database,
		NodeJS,
	}
}

const (
	PHP SupportedService = iota + 1
	Nginx
	Database
	NodeJS
)

// ServicesConfig contains config for each service
type ServicesConfig struct {
	PHP      *PHPConfig
	Nginx    *NginxConfig
	Database *DatabaseConfig
	NodeJS   *NodeJSConfig
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

// IsPresent contains whether service is present in the config
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

// PresentServicesCount counts how many supported services are present in the config
func (s *ServicesConfig) PresentServicesCount() int {
	count := 0

	for _, sup := range SupportedServices() {
		if s.IsPresent(sup) {
			count++
		}
	}

	return count
}

func (s *ServicesConfig) presentServices() []SupportedService {
	var services []SupportedService

	for _, sup := range SupportedServices() {
		if s.IsPresent(sup) {
			services = append(services, sup)
		}
	}

	return services
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
