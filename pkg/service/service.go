package service

import "fmt"

type ServicesConfig struct {
	PHP      *PHP
	NodeJS   *NodeJS
	Nginx    *Nginx
	Database *Database
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

func (s *ServicesConfig) IsPresent(service string) bool {
	switch service {
	case "php":
		return s.PHP != nil && !s.PHP.IsEmpty()
	case "nodejs":
		return s.NodeJS != nil && !(*s.NodeJS == NodeJS{})
	case "nginx":
		return s.Nginx != nil && !(*s.Nginx == Nginx{})
	case "database":
		return s.Database != nil && !(*s.Database == Database{})
	default:
		return false
	}
}

func (s *ServicesConfig) String() string {
	return fmt.Sprintf(
		"ServicesConfig{PHP: %v, NodeJS: %v, Nginx: %v, Database: %v}",
		s.PHP,
		s.NodeJS,
		s.Nginx,
		s.Database,
	)
}
