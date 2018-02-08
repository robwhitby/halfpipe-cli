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
	bytes, err := c.FileSystem.ReadFile(path.Join(c.RootDir, ".halfpipe.io"))
	if err != nil {
		fmt.Fprintln(c.ErrorWriter, errorReport(err))
		return false
	}

	man, parseErrors := model.Parse(string(bytes))
	if len(parseErrors) > 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(parseErrors...))
		return false
	}

	if lintErrors := linter.Lint(man); len(lintErrors) > 0 {
		fmt.Fprintln(c.ErrorWriter, errorReport(lintErrors...))
		return false
	}

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
