package linter

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/robwhitby/halfpipe-cli/model"
)

func TestLint(t *testing.T) {
	g := NewGomegaWithT(t)
	_ = g

	man := &model.Manifest{
		Team: "rarr",
	}

	failures := Lint(man)

	g.Expect(len(failures)).To(Equal(1))
	g.Expect(failures[0]).To(Equal(&Failure{MissingField, "tasks"}))
}
