package linter

import (
	"testing"

	. "github.com/robwhitby/halfpipe-cli/model"
	"github.com/stretchr/testify/assert"
)

func TestLint(t *testing.T) {
	man := Manifest{
		Team: "meaT",
	}
	failures := Lint(man)

	assert.Equal(t, len(failures), 2, "total 2 failures")
	assert.Contains(t, failures, NewMissingField("repo.uri"))
	assert.Contains(t, failures, NewMissingField("tasks"))
}

func TestRepoUri(t *testing.T) {
	man := Manifest{
		Repo: Repo{Uri: "uri"},
	}
	failures := Lint(man)

	assert.Contains(t, failures, NewInvalidField("repo.uri", "must contain 'github'"))
}
