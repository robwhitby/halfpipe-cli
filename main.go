package main

import (
	"fmt"

	"syscall"

	"os"

	"github.com/jessevdk/go-flags"
	"github.com/robwhitby/halfpipe-cli/controller"
	"github.com/spf13/afero"
)

const version = "0.01"

type options struct {
	Version bool `short:"v" long:"version" description:"Display version"`
	Args    struct {
		Dir string `positional-arg-name:"directory" description:"Path to process. Defaults to pwd."`
	} `positional-args:"true"`
}

func main() {
	opts := new(options)
	flags.Parse(opts)

	if opts.Version {
		fmt.Println("halfpipe v" + version)
		return
	}

	fileSystem := afero.NewOsFs()
	secretChecker := func(s string) bool { return false } //todo: vault checker

	ctrl := controller.NewController(fileSystem, opts.Args.Dir, os.Stdout, os.Stderr, secretChecker)

	if ok := ctrl.Run(); !ok {
		syscall.Exit(1)
	}
}
