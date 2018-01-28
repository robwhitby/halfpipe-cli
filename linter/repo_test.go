package linter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/manifest"
)

var _ = Describe("RepoLinter", func() {

	It("Returns error if repo.uri doesnt look like git repo", func() {
		man := manifest.Manifest{
			Repo: manifest.Repo{
				Uri: "gitasdasdt@github.com:simonjohansson/go-linter.git",
			},
		}
		results, _ := RepoLinter{}.Lint(man)
		Expect(len(results.Errors)).To(Equal(1))

		man = manifest.Manifest{
			Repo: manifest.Repo{
				Uri: "https://github.coasdm/cenkalti/backoff",
			},
		}
		results, _ = RepoLinter{}.Lint(man)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns error if repo.uri is private but no private key specified", func() {
		man := manifest.Manifest{
			Repo: manifest.Repo{
				Uri: "git@github.com:simonjohansson/go-linter.git",
			},
		}
		results, _ := RepoLinter{}.Lint(man)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns no errors if all is in order", func() {
		man := manifest.Manifest{
			Repo: manifest.Repo{
				Uri:        "git@github.com:simonjohansson/go-linter.git",
				PrivateKey: "asd",
			},
		}
		results, _ := RepoLinter{}.Lint(man)
		Expect(len(results.Errors)).To(Equal(0))

		man = manifest.Manifest{
			Repo: manifest.Repo{
				Uri: "https://github.com/cenkalti/backoff",
			},
		}
		results, _ = RepoLinter{}.Lint(man)
		Expect(len(results.Errors)).To(Equal(0))
	})
})
