package service

import (
	"path/filepath"
)

// FileType is one of supported file types
type FileType int

// All supported file types
const (
	Dockerfile FileType = iota + 1
	ConfigFile
)

// File represents a service file (e.g. service Dockerfile or service config)
type File struct {
	Type            FileType
	PathOnHost      string
	PathInContainer string
	TemplatePath    string
}

// GetTemplatePath returns path to template from which resulting service file can be rendered
func (f *File) GetTemplatePath() string {
	return f.TemplatePath
}

// GetOutputPath returns path to which resulting service file will be rendered
func (f *File) GetOutputPath() string {
	return f.PathOnHost
}

// IsMountable determines whether file can be mounted inside container
func (f *File) IsMountable() bool {
	return f.PathOnHost != "" && f.PathInContainer != ""
}

func getFilesForService(service SupportedService, outputPath string) []*File {
	switch service {
	case PHP:
		return []*File{
			{
				Type:         Dockerfile,
				PathOnHost:   filepath.Join(outputPath, "php/Dockerfile"),
				TemplatePath: "/php/php.dockerfile.gotmpl",
			},
		}
	case Nginx:
		return []*File{
			{
				Type:            ConfigFile,
				PathOnHost:      filepath.Join(outputPath, "nginx/conf.d/app.conf"),
				PathInContainer: "/etc/nginx/conf.d/app.conf",
				TemplatePath:    "/nginx/conf.gotmpl",
			},
		}
	case NodeJS:
		return []*File{
			{
				Type:         Dockerfile,
				PathOnHost:   filepath.Join(outputPath, "nodejs/Dockerfile"),
				TemplatePath: "/nodejs/nodejs.dockerfile.gotmpl",
			},
		}
	default:
		return nil
	}
}

// Files is a per-service collection of files
type Files map[SupportedService][]*File
