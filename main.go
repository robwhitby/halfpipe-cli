package main

import (
	"fmt"
	"syscall"

	"os"

	"github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
)

func main() {

	fileSystem := &afero.Afero{Fs: afero.NewOsFs()}

	bytes, err := fileSystem.ReadFile(os.Getenv("HOME") + "/go/src/github.com/robwhitby/halfpipe-cli/.halfpipe.io")
	exitWithError(err)

	manifestYaml := string(bytes)

	//parse
	man, parseFailures := model.Parse(manifestYaml)
	if !parseFailures.IsEmpty() {
		fmt.Println("Failed to parse manifest:")
		exitWithError(parseFailures)
	}

	//lint
	lintFailures := linter.Lint(man)

	if len(lintFailures) > 0 {
		fmt.Printf("Found %v issues:\n", len(lintFailures))
		for _, f := range lintFailures {
			fmt.Printf("- %s\n", f.Message)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Println("Manifest object:")
	fmt.Printf("%+v\n", man)

}

func exitWithError(err interface{}) {
	if err != nil {
		fmt.Println(err)
		syscall.Exit(-1)
	}
}
