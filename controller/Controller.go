package controller

import (
	"fmt"
	"strings"

	"io"

	"github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/robwhitby/halfpipe-cli/parser"
	"github.com/spf13/afero"
)

const (
	documentationRootUrl = "http://docs.halfpipe.io"
	manifestFilename     = ".halfpipe.io"
)

type controller struct {
	FileSystem    afero.Afero
	OutputWriter  io.Writer
	ErrorWriter   io.Writer
	SecretChecker model.SecretChecker
}

func NewController(fileSystem afero.Fs, repoDir string, outWriter, errWriter io.Writer, secretChecker model.SecretChecker) controller {
	return controller{
		FileSystem:    afero.Afero{Fs: afero.NewBasePathFs(fileSystem, repoDir)},
		OutputWriter:  outWriter,
		ErrorWriter:   errWriter,
		SecretChecker: secretChecker,
	}
}

func (c controller) Run() (ok bool) {
	//read manifest file
	yaml, err := readManifest(c.FileSystem)
	if err != nil {
		fmt.Fprintln(c.ErrorWriter, errorReport(err))
		return false
	}

	// parse it into a model.Manifest
	man, parseErrors := parser.Parse(yaml)
	if len(parseErrors) > 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(parseErrors...))
		return false
	}

	// lint it
	lintErrors := linter.LintManifest(man)
	lintErrors = append(lintErrors, linter.LintFiles(man, c.FileSystem)...)
	lintErrors = append(lintErrors, linter.LintSecrets(man, c.SecretChecker)...)

	if len(lintErrors) > 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(lintErrors...))
		return false
	}

	// TODO: generate the concourse yaml
	fmt.Fprintln(c.OutputWriter, "Good job")
	return true
}

func readManifest(fs afero.Afero) (string, error) {
	if err := linter.CheckFile(linter.RequiredFile{Path: manifestFilename}, fs); err != nil {
		return "", err
	}
	bytes, err := fs.ReadFile(manifestFilename)
	if err != nil {
		return "", model.NewFileError(manifestFilename, err.Error())
	}
	return string(bytes), nil
}

func errorReport(errs ...error) string {
	var lines []string
	lines = append(lines, "Found some problems:")
	for _, err := range errs {
		lines = append(lines, "- "+err.Error())
		if docs, ok := err.(model.Documented); ok {
			lines = append(lines, fmt.Sprintf("  rtfm: %s%s", documentationRootUrl, docs.DocumentationPath()))
		}
	}
	return strings.Join(lines, "\n")
}
