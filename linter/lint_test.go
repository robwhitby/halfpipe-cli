package linter

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/robwhitby/halfpipe-cli/model"
)

func TestLint(t *testing.T) {
	g := NewGomegaWithT(t)

	man := &model.Manifest{
		Team: "meaT",
	}

	failures := Lint(man)
	g.Expect(len(failures)).To(Equal(2))
	g.Expect(failures).To(ContainElement(missingField("repo.uri")))
	g.Expect(failures).To(ContainElement(missingField("tasks")))
}

func TestRepoUri(t *testing.T) {
	g := NewGomegaWithT(t)

	man := &model.Manifest{
		Repo: model.Repo{Uri: "uri"},
	}

	failures := Lint(man)
	g.Expect(failures).To(ContainElement(invalidValue("repo.uri", "must contain 'github'")))
}
