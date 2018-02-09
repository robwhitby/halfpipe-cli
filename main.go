package main

import (
	"os"

	"syscall"

	"github.com/robwhitby/halfpipe-cli/controller"
	"github.com/spf13/afero"
)

func main() {

	ctrl := &controller.Controller{
		FileSystem:   afero.Afero{Fs: afero.NewOsFs()},
		RootDir:      os.Getenv("HOME") + "/go/src/github.com/robwhitby/halfpipe-cli",
		OutputWriter: os.Stdout,
		ErrorWriter:  os.Stderr,
	}

	if ok := ctrl.Run(); !ok {
		syscall.Exit(1)
	}

}
