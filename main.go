package main

import (
	"fmt"

	"os"

	"syscall"

	"github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
)

func main() {

	fileSystem := &afero.Afero{Fs: afero.NewOsFs()}

	bytes, err := fileSystem.ReadFile(os.Getenv("HOME") + "/go/src/github.com/robwhitby/halfpipe-cli/.halfpipe.io")
	if err != nil {
		exitWithErrors(err)
	}

	manifestYaml := string(bytes)

	//parse
	man, parseErrors := model.Parse(manifestYaml)
	if len(parseErrors) > 0 {
		exitWithErrors(parseErrors...)
	}

	//lint
	if lintErrors := linter.Lint(man); len(lintErrors) > 0 {
		exitWithErrors(lintErrors...)
	}

	fmt.Println("Good job")

}

func exitWithErrors(errs ...error) {
	fmt.Println("Found some problems:")
	for _, e := range errs {
		fmt.Printf("- %+v\n", e)
	}
	syscall.Exit(-1)
}
