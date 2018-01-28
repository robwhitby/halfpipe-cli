package manifest

import (
	"regexp"
	"strings"
)

type Manifest struct {
	Team  string `yaml:"team"`
	Repo  Repo   `yaml:"repo"`
	Tasks []Task `yaml:"tasks"`
}

type Repo struct {
	Uri        string `yaml:"uri"`
	PrivateKey string `yaml:"private_key"`
}

func (r Repo) RepoName() string {
	if strings.HasPrefix(r.Uri, "git@github.com") {
		re, _ := regexp.Compile(`.*/(.*).git$`)
		return re.FindAllStringSubmatch(r.Uri, -1)[0][1]
	}
	if strings.HasPrefix(r.Uri, "https://github.com/") {
		re, _ := regexp.Compile(`.*/(.*)$`)
		return re.FindAllStringSubmatch(r.Uri, -1)[0][1]
	}
	return ""
}

type Task interface{}

type RunTask struct {
	Script string
	Image  string
	Vars   map[string]string
}

type DockerTask struct {
	Username   string
	Password   string
	Repository string
}

type DeployTask struct {
	Api      string
	Org      string
	Space    string
	Username string
	Password string
	Manifest string
	Vars     map[string]string
}
