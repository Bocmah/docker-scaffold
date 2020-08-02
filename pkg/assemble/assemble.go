package assemble

import (
	"fmt"
	"strings"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

type ServiceAssembler func(conf *service.FullConfig, opts ...Option) *dockercompose.Service

func NewServiceAssembler(s service.SupportedService) ServiceAssembler {
	switch s {
	case service.PHP:
		return phpAssembler()
	default:
		return unknownAssembler()
	}
}

func phpAssembler() ServiceAssembler {
	return func(conf *service.FullConfig, opts ...Option) *dockercompose.Service {
		if !conf.Services.IsPresent(service.PHP) {
			return nil
		}

		options := options{
			dockerfilePath: "",
		}

		for _, o := range opts {
			o.apply(&options)
		}

		workDir := "/var/www"
		appName := formatAppName(conf.AppName)

		s := dockercompose.Service{
			Name:          appName,
			ContainerName: appName,
			Restart:       dockercompose.RestartPolicyUnlessStopped,
			WorkingDir:    workDir,
			Volumes: dockercompose.ServiceVolumes{
				&dockercompose.ServiceVolume{Source: conf.ProjectRoot, Target: workDir},
			},
		}

		if options.dockerfilePath != "" {
			s.Build = &dockercompose.Build{
				Context:    conf.ProjectRoot,
				Dockerfile: options.dockerfilePath,
			}

			s.Image = &dockercompose.Image{Name: appName}
		} else {
			s.Image = &dockercompose.Image{
				Name: "php",
				Tag:  fmt.Sprintf("%s-fpm", conf.Services.PHP.Version),
			}
		}

		if options.sharedNetwork != nil {
			s.Networks = dockercompose.ServiceNetworks{options.sharedNetwork}
		}

		return &s
	}
}

func unknownAssembler() ServiceAssembler {
	return func(conf *service.FullConfig, opts ...Option) *dockercompose.Service {
		return nil
	}
}

func formatAppName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}
