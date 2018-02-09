package linter

import (
	"strings"

	"github.com/robwhitby/halfpipe-cli/model"
)

func Lint(man model.Manifest) (errs []error) {
	addError := func(e error) {
		errs = append(errs, e)
	}

	if man.Team == "" {
		addError(model.NewMissingField("team"))
	}

	if man.Repo.Uri == "" {
		addError(model.NewMissingField("repo.uri"))
	} else if !strings.Contains(man.Repo.Uri, "github") {
		addError(model.NewInvalidField("repo.uri", "must contain 'github'"))
	}

	if len(man.Tasks) == 0 {
		addError(model.NewMissingField("tasks"))
	}

	return
}
