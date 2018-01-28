package render_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"fmt"

	"github.com/concourse/atc"
	"github.com/robwhitby/halfpipe-cli/manifest"
	. "github.com/robwhitby/halfpipe-cli/render"
)

func yamlString(s string) string {
	return strings.Replace(s, "\t", "  ", -1)
}

func expectEqualYAML(actual string, expected string) {
	var a interface{}
	var e interface{}

	var err error
	err = yaml.Unmarshal([]byte(actual), &a)
	if err != nil {
		fmt.Printf("Error parsing actual :: %s\n", err.Error())
	}
	err = yaml.Unmarshal([]byte(expected), &e)
	if err != nil {
		fmt.Printf("Error parsing expected :: %s\n", err.Error())
	}

	Expect(a).To(Equal(e))
}

var _ = Describe("ConcourseRenderer", func() {

	Context("Repo", func() {
		It("renders http repo to a git resource", func() {

			man := manifest.Manifest{
				Repo: manifest.Repo{
					Uri: "https://github.com/cloudfoundry/bosh-cli",
				},
			}
			pipeline := ConcourseRenderer{}.RenderToString(man)
			expected := yamlString(`
				groups: []
				resources:
				- name: bosh-cli
					type: git
					source:
						uri: https://github.com/cloudfoundry/bosh-cli
				resource_types: []
				jobs: []`)

			expectEqualYAML(pipeline, expected)
		})

		It("renders ssh repo to a git resource", func() {
			man := manifest.Manifest{
				Repo: manifest.Repo{
					Uri:        "git@github.com:springernature/ee-half-pipe-landing.git",
					PrivateKey: "((something.secret))",
				},
			}
			pipeline := ConcourseRenderer{}.RenderToString(man)
			expected := yamlString(`
				groups: []
				resources:
				- name: ee-half-pipe-landing
				  type: git
				  source:
						private_key: "((something.secret))"
						uri: git@github.com:springernature/ee-half-pipe-landing.git
				resource_types: []
				jobs: []`)

			expectEqualYAML(pipeline, expected)
		})
	})

	Context("Tasks", func() {
		Context("run", func() {
			It("Renders a task without vars correctly", func() {
				man := manifest.Manifest{
					Repo: manifest.Repo{
						Uri:        "git@github.com:springernature/ee-half-pipe-landing.git",
						PrivateKey: "((something.secret))",
					},
					Tasks: []manifest.Task{
						manifest.RunTask{
							Script: "yolo.sh",
							Image:  "something",
						},
					},
				}

				pipeline := ConcourseRenderer{}.Render(man)
				Expect(pipeline.Jobs[0]).To(Equal(atc.JobConfig{
					Name:   "yolo.sh",
					Serial: true,
					Plan: atc.PlanSequence{
						atc.PlanConfig{Get: man.Repo.RepoName(), Trigger: true},
						atc.PlanConfig{Task: "yolo.sh", TaskConfig: &atc.TaskConfig{
							Platform: "linux",
							ImageResource: &atc.ImageResource{
								Type: "docker-image",
								Source: atc.Source{
									"repository": "something",
									"tag":        "latest",
								},
							},
							Params: nil,
							Run: atc.TaskRunConfig{
								Path: "/bin/sh",
								Dir:  man.Repo.RepoName(),
								Args: []string{"-exc", "./yolo.sh"},
							},
							Inputs: []atc.TaskInputConfig{
								atc.TaskInputConfig{Name: man.Repo.RepoName()},
							},
						}},
					}}))
			})

			It("Renders a task with vars correctly", func() {
				man := manifest.Manifest{
					Repo: manifest.Repo{
						Uri:        "git@github.com:springernature/ee-half-pipe-landing.git",
						PrivateKey: "((something.secret))",
					},
					Tasks: []manifest.Task{
						manifest.RunTask{
							Script: "yolo.sh",
							Image:  "something",
							Vars: map[string]string{
								"VAR1": "Value",
								"VAR2": "Value",
							},
						},
					},
				}

				pipeline := ConcourseRenderer{}.Render(man)
				Expect(pipeline.Jobs[0]).To(Equal(atc.JobConfig{
					Name:   "yolo.sh",
					Serial: true,
					Plan: atc.PlanSequence{
						atc.PlanConfig{Get: man.Repo.RepoName(), Trigger: true},
						atc.PlanConfig{Task: "yolo.sh", TaskConfig: &atc.TaskConfig{
							Platform: "linux",
							Params: map[string]string{
								"VAR1": "Value",
								"VAR2": "Value",
							},
							ImageResource: &atc.ImageResource{
								Type: "docker-image",
								Source: atc.Source{
									"repository": "something",
									"tag":        "latest",
								},
							},
							Run: atc.TaskRunConfig{
								Path: "/bin/sh",
								Dir:  man.Repo.RepoName(),
								Args: []string{"-exc", fmt.Sprint("./yolo.sh")},
							},
							Inputs: []atc.TaskInputConfig{
								atc.TaskInputConfig{Name: man.Repo.RepoName()},
							},
						}},
					}}))
			})

		})
	})
})
