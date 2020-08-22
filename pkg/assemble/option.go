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
	databaseSystemInUse service.SupportedSystem
	compose             *dockercompose.Config
	serviceFiles        service.Files
	serviceEnv          service.Environment
}

func (o *optionsAssembler) assembleForService(serv service.SupportedService) []Option {
	var opts []Option

	if len(o.compose.Volumes) != 0 && serv == service.Database {
		databaseVols := o.getVolumesForDatabase()

		if len(databaseVols) != 0 {
			opts = append(opts, WithVolumes(databaseVols))
		}
	}

	if len(o.compose.Networks) != 0 {
		opts = append(opts, WithNetworks(o.compose.Networks.ToServiceNetworks()))
	}

	opts = append(opts, o.serviceFileOpts(serv)...)
	opts = append(opts, o.serviceEnvOpts(serv)...)

	return opts
}

func (o *optionsAssembler) getVolumesForDatabase() dockercompose.ServiceVolumes {
	vols := dockercompose.ServiceVolumes{}

	for _, vol := range o.compose.Volumes {
		vols = append(vols, &dockercompose.ServiceVolume{Source: vol.Name, Target: o.databaseSystemInUse.DataPath()})
	}

	return vols
}

func (o *optionsAssembler) serviceFileOpts(serv service.SupportedService) []Option {
	files, ok := o.serviceFiles[serv]

	if !ok {
		return nil
	}

	var opts []Option
	var volumes dockercompose.ServiceVolumes

	for _, file := range files {
		if file.Type == service.Dockerfile {
			opts = append(opts, WithDockerfilePath(file.PathOnHost))
		}

		if file.IsMountable() {
			volumes = append(volumes, &dockercompose.ServiceVolume{Source: file.PathOnHost, Target: file.PathInContainer})
		}
	}

	if len(volumes) != 0 {
		opts = append(opts, WithVolumes(volumes))
	}

	return opts
}

func (o *optionsAssembler) serviceEnvOpts(serv service.SupportedService) []Option {
	env, ok := o.serviceEnv[serv]

	if !ok {
		return nil
	}

	return []Option{WithEnvironment(env)}
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
