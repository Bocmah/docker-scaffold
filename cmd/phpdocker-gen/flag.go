package main

import (
	"bytes"
	"flag"
)

type Config struct {
	file string
	args []string
}

func parseFlags(progname string, args []string) (config *Config, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var conf Config
	flags.StringVar(&conf.file, "file", "", "File with services configuration")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}
	conf.args = flags.Args()
	return &conf, buf.String(), nil
}
