package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
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
