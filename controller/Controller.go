package controller

import (
	"fmt"
	"path"
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

type Controller struct {
	FileSystem   afero.Afero
	RootDir      string
	OutputWriter io.Writer
	ErrorWriter  io.Writer
}

func (c *Controller) Run() (ok bool) {
	//read manifest file
	yaml, err := readManifest(c.FileSystem, c.RootDir)
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
	manifestErrors := linter.LintManifest(man)
	fileErrors := linter.LintFiles(man, c.RootDir, c.FileSystem)

	allErrors := append(manifestErrors, fileErrors...)

	if len(allErrors) > 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(allErrors...))
		return false
	}

	// TODO: generate the concourse yaml
	fmt.Fprintln(c.OutputWriter, "Good job")
	return true
}

func readManifest(fs afero.Afero, rootDir string) (string, error) {
	if err := linter.CheckFile(model.RequiredFile{Path: manifestFilename}, rootDir, fs); err != nil {
		return "", err
	}
	bytes, err := fs.ReadFile(path.Join(rootDir, manifestFilename))
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
