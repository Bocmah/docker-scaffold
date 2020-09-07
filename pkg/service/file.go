package service

import (
	"path/filepath"
)

type FileType int

const (
	Dockerfile FileType = iota + 1
	ConfigFile
)

type File struct {
	Type            FileType
	PathOnHost      string
	PathInContainer string
	TemplatePath    string
}

func (f *File) GetTemplatePath() string {
	return f.TemplatePath
}

func (f *File) GetOutputPath() string {
	return f.PathOnHost
}

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
				TemplatePath: filepath.Join("../../tmpl", "php/php.dockerfile.gotmpl"),
			},
		}
	case Nginx:
		return []*File{
			{
				Type:            ConfigFile,
				PathOnHost:      filepath.Join(outputPath, "nginx/conf.d/app.conf"),
				PathInContainer: "/etc/nginx/conf.d/app.conf",
				TemplatePath:    filepath.Join("../../tmpl", "nginx/conf.gotmpl"),
			},
		}
	case NodeJS:
		return []*File{
			{
				Type:         Dockerfile,
				PathOnHost:   filepath.Join(outputPath, "nodejs/Dockerfile"),
				TemplatePath: filepath.Join("../../tmpl", "nodejs/nodejs.dockerfile.gotmpl"),
			},
		}
	default:
		return nil
	}
}

type Files map[SupportedService][]*File
