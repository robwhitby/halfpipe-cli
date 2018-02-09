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

	ctrl := &controller.Controller{
		FileSystem:   afero.Afero{Fs: afero.NewOsFs()},
		RootDir:      pwd,
		OutputWriter: os.Stdout,
		ErrorWriter:  os.Stderr,
	}

	if ok := ctrl.Run(); !ok {
		syscall.Exit(1)
	}

}
