package linter

import "github.com/robwhitby/halfpipe-cli/manifest"

type RequiredFieldsLinter struct{}

func (RequiredFieldsLinter) Lint(man manifest.Manifest) (Result, error) {
	errors := []Error{}
	if man.Team == "" {
		errors = append(errors, Error{Message: "Required top level field 'team' missing"})
	}
	if (man.Repo == manifest.Repo{}) {
		errors = append(errors, Error{Message: "Required top level field 'repo' missing"})
	}
	if len(man.Tasks) == 0 {
		errors = append(errors, Error{Message: "Tasks is empty..."})
	}

	return Result{
		Linter: "Required Fields",
		Errors: errors,
	}, nil
}
