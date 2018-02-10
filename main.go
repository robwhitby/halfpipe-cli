package main

import (
	"os"

	"syscall"

	"fmt"

	"github.com/robwhitby/halfpipe-cli/controller"
	"github.com/spf13/afero"
)

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		syscall.Exit(1)
	}

	fileSystem := afero.NewOsFs()
	secretChecker := func(s string) bool { return false } //todo: vault checker

	ctrl := controller.NewController(fileSystem, pwd, os.Stdout, os.Stderr, secretChecker)

	if ok := ctrl.Run(); !ok {
		syscall.Exit(1)
	}

}
