package service

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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

func (c *FullConfig) GetServiceFiles() map[SupportedService][]*File {
	outputPath := c.getOutputPath()
	files := map[SupportedService][]*File{}

	for _, service := range c.Services.presentServices() {
		files[service] = getFilesForService(service, outputPath)
	}

	return files
}

func (c *FullConfig) getOutputPath() string {
	if c.OutputPath != "" {
		return c.OutputPath
	}

	return c.ProjectRoot
}

func LoadConfigFromFile(filepath string) (*FullConfig, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read config: %s", err)
	}

	conf := &FullConfig{}

	if err := yaml.Unmarshal(data, conf); err != nil {
		return nil, fmt.Errorf("parse config: %s", err)
	}

	conf.FillDefaultsIfNotSet()
	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %s", err)
	}

	return conf, nil
}
