package linter

import (
	"fmt"
	"strings"

	"github.com/robwhitby/halfpipe-cli/model"
)

type Failure struct {
	Type    FailureType
	Message string
}

type FailureType int

const (
	MissingField FailureType = iota
	InvalidValue
)

func missingField(name string) *Failure {
	return &Failure{MissingField, name}
}

func invalidValue(name string, reason string) *Failure {
	return &Failure{InvalidValue, fmt.Sprintf("%s has invalid value: %s", name, reason)}
}

func Lint(man *model.Manifest) (failures []*Failure) {
	f := func(f *Failure) {
		failures = append(failures, f)
	}

	if man.Team == "" {
		f(missingField("repo.uri"))
	}

	if man.Repo.Uri == "" {
		f(missingField("repo.uri"))
	} else if !strings.Contains(man.Repo.Uri, "github") {
		f(invalidValue("repo.uri", "must contain 'github'"))
	}

	if len(man.Tasks) == 0 {
		f(missingField("tasks"))
	}

	return failures
}
