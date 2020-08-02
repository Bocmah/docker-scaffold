package assemble

import "github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"

type options struct {
	dockerfilePath string
	sharedNetwork  *dockercompose.Network
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

type sharedNetworkOption struct {
	Network *dockercompose.Network
}

func (sn sharedNetworkOption) apply(opts *options) {
	opts.sharedNetwork = sn.Network
}

func WithSharedNetwork(network *dockercompose.Network) Option {
	return sharedNetworkOption{Network: network}
}
