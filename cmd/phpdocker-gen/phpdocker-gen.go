package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/Bocmah/phpdocker-gen/internal/dockercompose"

	"github.com/Bocmah/phpdocker-gen/pkg/assemble"
	"github.com/Bocmah/phpdocker-gen/pkg/render"
	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

// AppFs is a filesystem in use
var AppFs = afero.NewOsFs()

func generateDocker(conf *Config) {
	checkFileWithConfigurationExists(conf.file)

	serviceConf := loadConfig(conf.file)

	renderServices(serviceConf)

	composeConf := assemble.DockerCompose(serviceConf)

	renderDockerCompose(composeConf, filepath.Join(serviceConf.GetOutputPath(), "docker-compose.yml"))
}

func checkFileWithConfigurationExists(filepath string) {
	if _, statErr := AppFs.Stat(filepath); os.IsNotExist(statErr) {
		printAndExit("Provided file with configuration was not found")
	}
}

func loadConfig(filepath string) *service.FullConfig {
	conf, loadConfigErr := service.LoadConfigFromFile(filepath)

	if loadConfigErr != nil {
		if _, ok := loadConfigErr.(*service.ValidationErrors); ok {
			printAndExit(fmt.Sprintf("File contains errors:\n\n%v", loadConfigErr))
		} else {
			printAndExit(fmt.Sprintf("Encountered error while loading config file:\n\n%v", loadConfigErr))
		}
	}

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

func main() {
	flagConf, output, err := parseFlags(os.Args[0], os.Args[1:])

	if err == flag.ErrHelp {
		fmt.Println(output)
		os.Exit(2)
	} else if err != nil {
		fmt.Println("got error:", err)
		fmt.Println("output:\n", output)
		os.Exit(1)
	}

	generateDocker(flagConf)
}
