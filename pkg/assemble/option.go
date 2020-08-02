package assemble

import "github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"

type options struct {
	dockerfilePath string
	environment    dockercompose.Environment
	networks       dockercompose.ServiceNetworks
	volumes        dockercompose.ServiceVolumes
}

type Option interface {
	apply(*options)
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
