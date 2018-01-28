package linter

import (
	"fmt"
	"path"

	"github.com/robwhitby/halfpipe-cli/manifest"
	"github.com/spf13/afero"
)

type RequiredFilesLinter struct {
	Fs       *afero.Afero
	RepoRoot string
}

func (r RequiredFilesLinter) Lint(man manifest.Manifest) (Result, error) {
	result := Result{
		Linter: "Required Files",
	}

	files := [1]string{
		"README.md",
	}

	for _, reqFile := range files {
		filePath := path.Join(r.RepoRoot, reqFile)
		exists, err := r.Fs.Exists(filePath)
		if err != nil {
			return result, err
		}
		if !exists {
			result.Errors = append(result.Errors, Error{
				Message:       fmt.Sprintf("Missing: %s", filePath),
				Documentation: "That's a schoolboy error",
			})
		}
	}

	return result, nil
}
