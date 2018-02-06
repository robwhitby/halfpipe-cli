package linter

import "github.com/robwhitby/halfpipe-cli/model"

type Failure struct {
	Type    FailureType
	Message string
}

func NewFailure(t FailureType, message string) *Failure {
	return &Failure{t, message}
}

type FailureType int

const (
	MissingField FailureType = iota
	InvalidValue
)

func Lint(man *model.Manifest) (failures []*Failure) {

	if len(man.Tasks) == 0 {
		failures = append(failures, NewFailure(MissingField, "tasks"))
	}

	return failures
}
