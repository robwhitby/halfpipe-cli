package linter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/manifest"
)

var _ = Describe("RequiredFields", func() {

	It("Returns error if team is missing", func() {
		manifest := manifest.Manifest{
			Repo:  manifest.Repo{Uri: "asdf"},
			Tasks: []manifest.Task{manifest.RunTask{}},
		}
		results, _ := RequiredFieldsLinter{}.Lint(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns error if repo is missing", func() {
		manifest := manifest.Manifest{
			Team:  "asdf",
			Tasks: []manifest.Task{manifest.RunTask{}},
		}
		results, _ := RequiredFieldsLinter{}.Lint(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns error if tasks is missing", func() {
		manifest := manifest.Manifest{
			Team: "asdf",
			Repo: manifest.Repo{Uri: "asdf"},
		}
		results, _ := RequiredFieldsLinter{}.Lint(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})
})
