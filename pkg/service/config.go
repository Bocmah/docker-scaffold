package service

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"gopkg.in/yaml.v2"
)

// AppFs is the filesystem in use
var AppFs = afero.NewMemMapFs()

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

// FillDefaultsIfNotSet fills default parameters (if they are not present) for all services in the config
func (c *FullConfig) FillDefaultsIfNotSet() {
	c.Services.FillDefaultsIfNotSet()
}

// Validate validates all service parameters in the config
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
		if filesForService := getFilesForService(service, outputPath); filesForService != nil {
			files[service] = filesForService
		}
	}

	return files
}

// GetEnvironment returns collection of environment variables for services which require them
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
	data, err := afero.ReadFile(AppFs, filepath)
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
