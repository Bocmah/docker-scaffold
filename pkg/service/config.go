package service

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config is interface for configs
type Config interface {
	// FillDefaultsIfNotSet fills values which are not present in the config with their default-equivalent
	FillDefaultsIfNotSet()
	// Validate validates config
	Validate() error
}

// FullConfig is user-filled config from which resulted docker files will be generated
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

// GetServiceFiles returns paths to service files (Dockerfiles, configs, etc.) for each service in the config
func (c *FullConfig) GetServiceFiles() Files {
	outputPath := c.GetOutputPath()
	files := map[SupportedService][]*File{}

	for _, service := range c.Services.presentServices() {
		files[service] = getFilesForService(service, outputPath)
	}

	return files
}

// GetEnvironment returns collection of environment variables for services which requires them
func (c *FullConfig) GetEnvironment() Environment {
	if !c.Services.IsPresent(Database) {
		return nil
	}

	return Environment{
		Database: c.Services.Database.Environment(),
	}
}

// GetOutputPath returns output path for resulting docker files
func (c *FullConfig) GetOutputPath() string {
	if c.OutputPath != "" {
		return c.OutputPath
	}

	return filepath.Join(c.ProjectRoot, ".docker")
}

// LoadConfigFromFile reads file at filepath, validates data and transforms it into FullConfig
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
