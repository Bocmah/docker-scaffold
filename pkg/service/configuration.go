package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	AppName  string `yaml:"appName"`
	PHP      *PHP
	NodeJS   *NodeJS
	Nginx    *Nginx
	Database *Database
}

func (c *Configuration) FillDefaultsIfNotSet() {
	c.Database.FillDefaultsIfNotSet()

	c.PHP.AddDatabaseExtension(c.Database.System)

	c.PHP.FillDefaultsIfNotSet()
	c.NodeJS.FillDefaultsIfNotSet()
	c.Nginx.FillDefaultsIfNotSet()
}

func (c *Configuration) Validate() error {
	errors := &ValidationErrors{}

	if c.AppName == "" {
		errors.Add("App name is required")
	}

	services := []Service{c.PHP, c.NodeJS, c.Nginx, c.Database}

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

func (c *Configuration) IsPresent(service string) bool {
	switch service {
	case "php":
		return c.PHP != nil && !c.PHP.IsEmpty()
	case "nodejs":
		return c.NodeJS != nil && !(*c.NodeJS == NodeJS{})
	case "nginx":
		return c.Nginx != nil && !(*c.Nginx == Nginx{})
	case "database":
		return c.Database != nil && !(*c.Database == Database{})
	default:
		return false
	}
}

func LoadConfigFromFile(filepath string) (*Configuration, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	conf := &Configuration{}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	conf.FillDefaultsIfNotSet()
	err = conf.Validate()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
