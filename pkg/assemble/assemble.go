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
	case service.Nginx:
		return nginxAssembler()
	case service.Database:
		return databaseAssembler()
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

		if len(options.volumes) != 0 {
			s.Volumes = append(s.Volumes, options.volumes...)
		}

		if len(options.networks) != 0 {
			s.Networks = options.networks
		}

		return &s
	}
}

func nginxAssembler() ServiceAssembler {
	return func(conf *service.FullConfig, opts ...Option) *dockercompose.Service {
		if !conf.Services.IsPresent(service.Nginx) {
			return nil
		}

		options := options{
			dockerfilePath: "",
		}

		for _, o := range opts {
			o.apply(&options)
		}

		var HTTPPort int

		if conf.Services.Nginx.HTTPPort != 0 {
			HTTPPort = conf.Services.Nginx.HTTPPort
		} else {
			HTTPPort = 80
		}

		var HTTPSPort int

		if conf.Services.Nginx.HTTPSPort != 0 {
			HTTPSPort = conf.Services.Nginx.HTTPSPort
		} else {
			HTTPSPort = 443
		}

		s := dockercompose.Service{
			Name: "webserver",
			Image: &dockercompose.Image{
				Name: "nginx",
				Tag:  "alpine",
			},
			ContainerName: "webserver",
			Restart:       dockercompose.RestartPolicyUnlessStopped,
			Ports: dockercompose.Ports{
				&dockercompose.PortsMapping{Host: HTTPPort, Container: HTTPPort},
				&dockercompose.PortsMapping{Host: HTTPSPort, Container: HTTPSPort},
			},
			Volumes: dockercompose.ServiceVolumes{
				&dockercompose.ServiceVolume{Source: conf.ProjectRoot, Target: "/var/www"},
			},
		}

		if len(options.volumes) != 0 {
			s.Volumes = append(s.Volumes, options.volumes...)
		}

		if len(options.networks) != 0 {
			s.Networks = append(s.Networks, options.networks...)
		}

		return &s
	}
}

func databaseAssembler() ServiceAssembler {
	return func(conf *service.FullConfig, opts ...Option) *dockercompose.Service {
		if !conf.Services.IsPresent(service.Database) {
			return nil
		}

		options := options{
			dockerfilePath: "",
		}

		for _, o := range opts {
			o.apply(&options)
		}

		s := dockercompose.Service{
			Name: "db",
			Image: &dockercompose.Image{
				Name: string(conf.Services.Database.System),
				Tag:  conf.Services.Database.Version,
			},
			ContainerName: "db",
			Restart:       dockercompose.RestartPolicyUnlessStopped,
			Ports: dockercompose.Ports{
				&dockercompose.PortsMapping{Host: conf.Services.Database.Port, Container: conf.Services.Database.Port},
			},
		}

		if len(options.volumes) != 0 {
			s.Volumes = append(s.Volumes, options.volumes...)
		}

		if len(options.networks) != 0 {
			s.Networks = append(s.Networks, options.networks...)
		}

		if len(options.environment) != 0 {
			if len(s.Environment) == 0 {
				s.Environment = dockercompose.Environment{}
			}

			for variable, val := range options.environment {
				s.Environment[variable] = val
			}
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
