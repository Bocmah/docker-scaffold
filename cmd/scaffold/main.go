package main

import (
	"flag"
	"fmt"
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
	"os"
)

func main() {
	file := flag.String("file", "", "File with services configuration")
	//_ := flag.String("output", ".", "Directory to which docker configuration files will be written")

	flag.Parse()

	if *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if _, err := os.Stat(*file); os.IsNotExist(err) {
		fmt.Println("Provided file with configuration was not found")
		os.Exit(1)
	}

	conf, err := service.LoadConfigFromFile(*file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(conf)
}