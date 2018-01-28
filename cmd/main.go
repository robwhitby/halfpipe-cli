package main

import (
	"fmt"
	"syscall"

	"github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/manifest"
	"github.com/robwhitby/halfpipe-cli/project"
	"github.com/robwhitby/halfpipe-cli/render"
	"github.com/spf13/afero"
)

func main() {

	//dependencies
	fileSystem := &afero.Afero{Fs: afero.NewOsFs()}
	config, err := project.NewConfig()
	exitOnError(err)

	//define linters
	linters := [3]linter.Linter{
		&linter.RequiredFieldsLinter{},
		&linter.RepoLinter{},
		&linter.RequiredFilesLinter{Fs: fileSystem, RepoRoot: config.RepoRoot},
	}

	//read manifest file
	manifest, err := manifest.NewManifestReader(fileSystem).ParseManifest(config.RepoRoot)
	exitOnError(err)

	// loop through linters
	hasLintErrors := false
	for _, linter := range linters {
		result, err := linter.Lint(manifest)
		exitOnError(err)
		hasLintErrors = hasLintErrors || result.HasLintErrors()
		fmt.Println(result.String())
	}

	if hasLintErrors {
		syscall.Exit(-1)
	}

	//no errors so output concourse pipeline
	concoursePipeline := render.ConcourseRenderer{}.RenderToString(manifest)
	fmt.Println(concoursePipeline)
}

func exitOnError(err interface{}) {
	if err != nil {
		fmt.Println(err)
		syscall.Exit(-1)
	}
}
