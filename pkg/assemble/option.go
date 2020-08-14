package assemble

import (
	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

type options struct {
	dockerfilePath string
	environment    dockercompose.Environment
	networks       dockercompose.ServiceNetworks
	volumes        dockercompose.ServiceVolumes
}

type Option interface {
	apply(*options)
}

type optionsAssembler struct {
	compose      *dockercompose.Config
	serviceFiles map[service.SupportedService]ServiceFiles
}

func (o *optionsAssembler) assembleForService(serv service.SupportedService) []Option {
	var opts []Option

	if len(o.compose.Volumes) != 0 && serv == service.Database {
		opts = append(opts, WithVolumes(o.compose.Volumes.ToServiceVolumes()))
	}

	if len(o.compose.Networks) != 0 {
		opts = append(opts, WithNetworks(o.compose.Networks.ToServiceNetworks()))
	}

	opts = append(opts, o.serviceFileOpts(serv)...)

	return opts
}

func (o *optionsAssembler) serviceFileOpts(service service.SupportedService) []Option {
	files, ok := o.serviceFiles[service]

	if !ok {
		return nil
	}

	var opts []Option

	if files.DockerfilePath != "" {
		opts = append(opts, WithDockerfilePath(files.DockerfilePath))
	}

	if len(files.Mounts) > 0 {
		opts = append(opts, WithVolumes(files.Mounts))
	}

	if len(files.Environment) > 0 {
		opts = append(opts, WithEnvironment(files.Environment))
	}

	return opts
}

type dockerfilePathOption string

func (dp dockerfilePathOption) apply(opts *options) {
	opts.dockerfilePath = string(dp)
}

func WithDockerfilePath(path string) Option {
	return dockerfilePathOption(path)
}

type environmentOption struct {
	Environment dockercompose.Environment
}

func (e environmentOption) apply(opts *options) {
	opts.environment = e.Environment
}

func WithEnvironment(env dockercompose.Environment) Option {
	return environmentOption{Environment: env}
}

type networksOption struct {
	Networks dockercompose.ServiceNetworks
}

func (n networksOption) apply(opts *options) {
	opts.networks = n.Networks
}

func WithNetworks(networks dockercompose.ServiceNetworks) Option {
	return networksOption{Networks: networks}
}

type volumesOption struct {
	Volumes dockercompose.ServiceVolumes
}

func (v volumesOption) apply(opts *options) {
	opts.volumes = v.Volumes
}

func WithVolumes(volumes dockercompose.ServiceVolumes) Option {
	return volumesOption{Volumes: volumes}
}
