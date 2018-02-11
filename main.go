package main

import (
	"syscall"

	"os"

	"github.com/jessevdk/go-flags"
	"github.com/robwhitby/halfpipe-cli/controller"
	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
)

func main() {
	opts := model.Options{}
	flags.Parse(opts)

	fileSystem := afero.NewOsFs()
	secretChecker := func(s string) bool { return false } //todo: vault checker

	ctrl := controller.NewController(fileSystem, opts, os.Stdout, os.Stderr, secretChecker)

	if ok := ctrl.Run(); !ok {
		syscall.Exit(1)
	}
}
