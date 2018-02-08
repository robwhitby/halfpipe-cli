package controller

import (
	"fmt"
	"path"
	"strings"

	"io"

	"github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
)

type Controller struct {
	FileSystem   afero.Afero
	RootDir      string
	OutputWriter io.Writer
	ErrorWriter  io.Writer
}

func (c *Controller) Run() (ok bool) {
	manifestPath := path.Join(c.RootDir, ".halfpipe.io")

	// does manifest exist?
	if exists, _ := c.FileSystem.Exists(manifestPath); !exists {
		fmt.Fprintln(c.ErrorWriter, errorReport(model.NewMissingFile(manifestPath)))
		return false
	}

	// read the file
	bytes, err := c.FileSystem.ReadFile(manifestPath)
	if err != nil {
		fmt.Fprintln(c.ErrorWriter, errorReport(model.NewParseError(err.Error())))
		return false
	}

	// was the file empty?
	if len(bytes) == 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(model.NewParseError(manifestPath+" is empty")))
		return false
	}

	// parse it into a model.Manifest
	man, parseErrors := model.Parse(string(bytes))
	if len(parseErrors) > 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(parseErrors...))
		return false
	}

	// lint it
	if lintErrors := linter.Lint(man); len(lintErrors) > 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(lintErrors...))
		return false
	}

	// TODO: generate the concourse yaml
	fmt.Fprintln(c.OutputWriter, "Good job")
	return true
}

func errorReport(errs ...error) string {
	var lines []string
	lines = append(lines, "Found some problems:")
	for _, err := range errs {
		lines = append(lines, "- "+err.Error())
		if docs, ok := err.(model.Documented); ok {
			lines = append(lines, fmt.Sprintf("  rtfm: http://docs.halfpipe.io%s", docs.DocumentationPath()))
		}
	}
	return strings.Join(lines, "\n")
}
