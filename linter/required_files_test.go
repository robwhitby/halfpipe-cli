package linter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/robwhitby/halfpipe-cli/linter"
	"github.com/robwhitby/halfpipe-cli/manifest"
	"github.com/spf13/afero"
)

var _ = Describe("RequiredFiles", func() {
	var (
		fs             *afero.Afero
		repoRoot       string
		reqFilesLinter RequiredFilesLinter
	)

	BeforeEach(func() {
		fs = &afero.Afero{Fs: afero.NewMemMapFs()}
		repoRoot = "/path/to/repo"
		reqFilesLinter = RequiredFilesLinter{fs, repoRoot}
	})

	It("Returns error if README.md is missing", func() {
		results, _ := reqFilesLinter.Lint(manifest.Manifest{})
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns empty error if README.md is present", func() {
		fs.WriteFile("/path/to/repo/README.md", []byte(""), 0644)
		results, _ := reqFilesLinter.Lint(manifest.Manifest{})
		Expect(len(results.Errors)).To(Equal(0))
	})
})
