package service

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"gopkg.in/yaml.v2"
)

// AppFs is the filesystem in use
var AppFs = afero.NewOsFs()

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
	if c.Services != nil {
		c.Services.FillDefaultsIfNotSet()
	}
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

	if c.Services == nil || c.Services.PresentServicesCount() == 0 {
		errors.Add("At least one service is required")
	}

	if c.Services != nil {
		errs := c.Services.Validate()

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
	data, readFileErr := afero.ReadFile(AppFs, filepath)
	if readFileErr != nil {
		return nil, fmt.Errorf("read config: %s", readFileErr)
	}

	conf := &FullConfig{}

	if unmarshallErr := yaml.Unmarshal(data, conf); unmarshallErr != nil {
		return nil, fmt.Errorf("parse config: %s", unmarshallErr)
	}

	conf.FillDefaultsIfNotSet()
	if validateErr := conf.Validate(); validateErr != nil {
		return nil, validateErr
	}

	return conf, nil
}
