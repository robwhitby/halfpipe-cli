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

	ctrl := controller.NewController(afero.NewOsFs(), pwd, os.Stdout, os.Stderr)

	if ok := ctrl.Run(); !ok {
		syscall.Exit(1)
	}

}
