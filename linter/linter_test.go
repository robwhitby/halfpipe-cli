package linter

import (
	"testing"

	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/stretchr/testify/assert"
)

func TestLint(t *testing.T) {
	man := &model.Manifest{
		Team: "meaT",
	}
	failures := Lint(man)

	assert.Equal(t, len(failures), 2, "total 2 failures")
	assert.Contains(t, failures, model.NewMissingField("repo.uri"))
	assert.Contains(t, failures, model.NewMissingField("tasks"))
}

func TestRepoUri(t *testing.T) {
	man := &model.Manifest{
		Repo: model.Repo{Uri: "uri"},
	}
	failures := Lint(man)

	assert.Contains(t, failures, model.NewInvalidField("repo.uri", "must contain 'github'"))
}
