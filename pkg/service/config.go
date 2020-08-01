package service

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"

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

func (c *FullConfig) MapToDockerCompose() *dockercompose.Config {
	conf := dockercompose.Config{Version: "3.8"}

	appName := strings.ReplaceAll(strings.ToLower(c.AppName), " ", "-")

	network := dockercompose.Network{
		Name:   fmt.Sprintf("%s-network", appName),
		Driver: dockercompose.NetworkDriverBridge,
	}

	conf.Networks = dockercompose.Networks{&network}

	vol := dockercompose.NamedVolume{
		Name:   fmt.Sprintf("%s-data", appName),
		Driver: dockercompose.VolumeDriverLocal,
	}

	conf.Volumes = dockercompose.NamedVolumes{&vol}

	/*var services []*dockercompose.Service

	if c.Services.IsPresent() {

	}*/

	return &conf
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
