package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Bocmah/phpdocker-scaffold/internal/dockercompose"

	"github.com/Bocmah/phpdocker-scaffold/pkg/assemble"
	"github.com/Bocmah/phpdocker-scaffold/pkg/render"
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

func main() {
	file := flag.String("file", "", "File with services configuration")

	flag.Parse()

	if *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	checkFileWithConfigurationExists(*file)

	conf := loadConfig(*file)

	renderServices(conf)

	composeConf := assemble.DockerCompose(conf)

	renderDockerCompose(composeConf, conf.GetOutputPath())
}

func checkFileWithConfigurationExists(filepath string) {
	if _, statErr := os.Stat(filepath); os.IsNotExist(statErr) {
		printAndExit("Provided file with configuration was not found")
	}
}

func loadConfig(filepath string) *service.FullConfig {
	conf, loadConfigErr := service.LoadConfigFromFile(filepath)
	checkErr(loadConfigErr)

	return conf
}

func renderServices(conf *service.FullConfig) {
	rendered, renderErr := render.RenderServices(conf)

	if renderErr == nil {
		return
	}

	if rendered != nil {
		if deleteErr := rendered.DeleteAllCreatedFiles(); deleteErr != nil {
			fmt.Println(deleteErr)
		}
	}

	printAndExit(renderErr.Error())
}

func renderDockerCompose(conf *dockercompose.Config, outputPath string) {
	renderErr := render.RenderDockerCompose(conf, outputPath)
	checkErr(renderErr)
}

func checkErr(err error) {
	if err != nil {
		printAndExit(err.Error())
	}
}

func printAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
