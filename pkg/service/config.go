package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config interface {
	FillDefaultsIfNotSet()
	Validate() error
}

type FullConfig struct {
	AppName     string `yaml:"appName"`
	ProjectRoot string `yaml:"projectRoot"`
	OutputPath  string `yaml:"outputPath"`
	Services    *ServicesConfig
}

func (c *FullConfig) FillDefaultsIfNotSet() {
	c.Services.FillDefaultsIfNotSet()
}

func (c *FullConfig) Validate() error {
	errors := &ValidationErrors{}

	if c.AppName == "" {
		errors.Add("App name is required")
	}

	if c.ProjectRoot == "" {
		errors.Add("Project root is required")
	}

	errs := c.Services.Validate()

	if errs != nil {
		if e, ok := errs.(*ValidationErrors); ok {
			errors.Merge(e)
		} else {
			errors.Add(errs.Error())
		}
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func LoadConfigFromFile(filepath string) (*FullConfig, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	conf := &FullConfig{}

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
