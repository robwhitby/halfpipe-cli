package main

import (
	"syscall"

	"os"

	"github.com/jessevdk/go-flags"
	"github.com/robwhitby/halfpipe-cli/controller"
	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
)

var version = "0.0" // should be set by build

func main() {
	config := model.Config{
		FileSystem:    afero.NewOsFs(),
		OutputWriter:  os.Stdout,
		ErrorWriter:   os.Stderr,
		SecretChecker: func(s string) bool { return false },
		Version:       version,
	}
	flags.Parse(&config.Options)

	if ok := controller.Process(config); !ok {
		syscall.Exit(1)
	}
}
