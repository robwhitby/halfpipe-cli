package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robwhitby/halfpipe-cli/manifest"
	"github.com/spf13/afero"
)

var _ = Describe("RequiredFiles", func() {
	var (
		fs           *afero.Afero
		reader       *manifest.ManifestReader
		repoRoot     = "/path/to/repo"
		manifestPath = repoRoot + "/.halfpipe.io"
	)

	BeforeEach(func() {
		fs = &afero.Afero{Fs: afero.NewMemMapFs()}
		reader = manifest.NewManifestReader(fs)
	})

	It("Returns error if .halfpipe.io is missing", func() {
		_, err := reader.ParseManifest(repoRoot)
		Expect(err).To(HaveOccurred())
	})

	It("Returns empty manifest if .halfpipe.io is empty", func() {
		fs.WriteFile(manifestPath, []byte(""), 0644)

		man, err := reader.ParseManifest(repoRoot)
		Expect(err).To(Not(HaveOccurred()))
		Expect(man).To(Equal(manifest.Manifest{}))
	})

	It("Parses empty .halfpipe.io to empty manifest", func() {
		content := ``
		fs.WriteFile(manifestPath, []byte(content), 0644)

		man, err := reader.ParseManifest(repoRoot)
		Expect(err).To(Not(HaveOccurred()))
		Expect(man).To(Equal(manifest.Manifest{}))
	})

	It("Parses minimal .halfpipe.io to minimal manifest", func() {
		content := `
team: engineering-enablement
`
		fs.WriteFile(manifestPath, []byte(content), 0644)

		man, err := reader.ParseManifest(repoRoot)
		Expect(err).To(Not(HaveOccurred()))
		Expect(man).To(Equal(manifest.Manifest{
			Team: "engineering-enablement",
		}))
	})

	It("Parses .halfpipe.io to manifest", func() {
		content := `
team: engineering-enablement
repo:
  uri: https://....
  private_key: asdf
tasks:
- task: run
  script: ./test.sh
  image: openjdk:8-slim
- task: docker
  username: ((docker.username))
  password: ((docker.password))
  repository: simonjohansson/half-pipe-linter
- task: deploy
  space: test
  api: https://api.europe-west1.cf.gcp.springernature.io
- task: run
  script: ./asdf.sh
  image: openjdk:8-slim
  vars:
    A: asdf
    B: 1234
- task: deploy
  space: test
  api: https://api.europe-west1.cf.gcp.springernature.io
  vars:
    VAR1: asdf1234
    VAR2: 9876
`
		fs.WriteFile(manifestPath, []byte(content), 0644)

		man, err := reader.ParseManifest(repoRoot)
		Expect(err).To(Not(HaveOccurred()))
		Expect(man).To(Equal(manifest.Manifest{
			Team: "engineering-enablement",
			Repo: manifest.Repo{
				Uri:        "https://....",
				PrivateKey: "asdf",
			},
			Tasks: []manifest.Task{
				manifest.RunTask{
					Script: "./test.sh",
					Image:  "openjdk:8-slim",
					Vars:   make(map[string]string),
				},
				manifest.DockerTask{
					Username:   "((docker.username))",
					Password:   "((docker.password))",
					Repository: "simonjohansson/half-pipe-linter",
				},
				manifest.DeployTask{
					Username: "((cf-credentials.username))",
					Password: "((cf-credentials.password))",
					Api:      "https://api.europe-west1.cf.gcp.springernature.io",
					Org:      "engineering-enablement",
					Space:    "test",
					Manifest: "manifest.yml",
					Vars:     make(map[string]string),
				},
				manifest.RunTask{
					Script: "./asdf.sh",
					Image:  "openjdk:8-slim",
					Vars: map[string]string{
						"A": "asdf",
						"B": "1234",
					},
				},
				manifest.DeployTask{
					Username: "((cf-credentials.username))",
					Password: "((cf-credentials.password))",
					Api:      "https://api.europe-west1.cf.gcp.springernature.io",
					Org:      "engineering-enablement",
					Space:    "test",
					Manifest: "manifest.yml",
					Vars: map[string]string{
						"VAR1": "asdf1234",
						"VAR2": "9876",
					},
				},
			},
		}))
	})

})
