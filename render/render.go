package render

import (
	"fmt"
	"strings"

	"github.com/concourse/atc"
	"github.com/robwhitby/halfpipe-cli/manifest"
	"gopkg.in/yaml.v2"
)

type Render interface {
	Render(manifest.Manifest)
}

type ConcourseRenderer struct {
}

func (ConcourseRenderer) makeGitConfig(repo manifest.Repo) atc.ResourceConfig {
	source := atc.Source{
		"uri": repo.Uri,
	}
	if repo.PrivateKey != "" {
		source["private_key"] = repo.PrivateKey
	}
	return atc.ResourceConfig{
		Name:   repo.RepoName(),
		Type:   "git",
		Source: source,
	}
}

func (ConcourseRenderer) dockerImageAndTag(image string) (string, string) {
	if strings.Contains(image, ":") {
		split := strings.Split(image, ":")
		return split[0], split[1]
	}
	return image, "latest"
}

func (c ConcourseRenderer) makeRunJob(task manifest.RunTask, repo manifest.Repo) atc.JobConfig {
	image, tag := c.dockerImageAndTag(task.Image)
	return atc.JobConfig{
		Name:   task.Script,
		Serial: true,
		Plan: atc.PlanSequence{
			atc.PlanConfig{Get: repo.RepoName(), Trigger: true},
			atc.PlanConfig{
				Task: task.Script,
				TaskConfig: &atc.TaskConfig{
					Platform: "linux",
					Params:   task.Vars,
					ImageResource: &atc.ImageResource{
						Type: "docker-image",
						Source: atc.Source{
							"repository": image,
							"tag":        tag,
						},
					},
					Run: atc.TaskRunConfig{
						Path: "/bin/sh",
						Dir:  repo.RepoName(),
						Args: []string{"-exc", fmt.Sprintf("./%s", task.Script)},
					},
					Inputs: []atc.TaskInputConfig{
						atc.TaskInputConfig{Name: repo.RepoName()},
					},
				}}}}
}

func (c ConcourseRenderer) Render(man manifest.Manifest) atc.Config {
	config := atc.Config{}
	config.Resources = append(config.Resources, c.makeGitConfig(man.Repo))
	for _, task := range man.Tasks {
		switch task.(type) {
		case manifest.RunTask:
			config.Jobs = append(config.Jobs, c.makeRunJob(task.(manifest.RunTask), man.Repo))
		}
	}
	return config
}

func (c ConcourseRenderer) RenderToString(man manifest.Manifest) string {
	pipeline := c.Render(man)
	renderedPipeline, _ := yaml.Marshal(pipeline)
	return string(renderedPipeline)
}
