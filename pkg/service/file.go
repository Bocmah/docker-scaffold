package service

import "os"

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

func getFilesForService(service SupportedService, outputPath string) []*File {
	switch service {
	case PHP:
		return []*File{
			{
				Type:         Dockerfile,
				PathOnHost:   outputPath + string(os.PathSeparator) + "php/Dockerfile",
				TemplatePath: "../../tmpl" + string(os.PathSeparator) + "php/php.dockerfile.gotmpl",
			},
		}
	case Nginx:
		return []*File{
			{
				Type:         ConfigFile,
				PathOnHost:   outputPath + string(os.PathSeparator) + "nginx/conf.d/app.conf",
				TemplatePath: "../../tmpl" + string(os.PathSeparator) + "nginx/conf.gotmpl",
			},
		}
	case NodeJS:
		return []*File{
			{
				Type:         Dockerfile,
				PathOnHost:   outputPath + string(os.PathSeparator) + "nodejs/Dockerfile",
				TemplatePath: "../../tmpl" + string(os.PathSeparator) + "nodejs/nodejs.dockerfile.gotmpl",
			},
		}
	default:
		return nil
	}
}
