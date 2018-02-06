package main

import (
	"fmt"
	"syscall"

	"os"

	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
)

func main() {

	fileSystem := &afero.Afero{Fs: afero.NewOsFs()}

	bytes, err := fileSystem.ReadFile(os.Getenv("HOME") + "/go/src/github.com/robwhitby/halfpipe-cli/.halfpipe.io")
	exitWithError(err)

	manifestYaml := string(bytes)
	man, failures := model.Parse(manifestYaml)

	if !failures.IsEmpty() {
		fmt.Println("Failed to parse manifest:")
		exitWithError(failures)
	}

	fmt.Println("Manifest object:")
	fmt.Printf("%+v\n", man)

}

func exitWithError(err interface{}) {
	if err != nil {
		fmt.Println(err)
		syscall.Exit(-1)
	}
}
