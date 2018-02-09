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

type Controller struct {
	FileSystem   afero.Afero
	RootDir      string
	OutputWriter io.Writer
	ErrorWriter  io.Writer
}

func (c *Controller) Run() (ok bool) {
	manifestPath := path.Join(c.RootDir, ".halfpipe.io")

	//read manifest file
	yaml, err := readFile(c.FileSystem, manifestPath)
	if err != nil {
		fmt.Fprintln(c.ErrorWriter, err)
		return false
	}

	// parse it into a model.Manifest
	man, parseErrors := parser.Parse(yaml)
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

func readFile(fs afero.Afero, path string) (string, error) {
	if exists, _ := fs.Exists(path); !exists {
		return "", model.NewMissingFile(path)
	}

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return "", model.NewParseError(err.Error())
	}

	if len(bytes) == 0 {
		return "", model.NewParseError(path + " is empty")
	}

	return string(bytes), nil
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
