package linter

import (
	"github.com/robwhitby/halfpipe-cli/manifest"
)

type Linter interface {
	Lint(manifest manifest.Manifest) (Result, error)
}
