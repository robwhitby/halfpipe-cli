package model

type task interface {
	GetName() string
}

type Run struct {
	Name     string
	Script   string
	Username string
	Image    string
	Vars     Vars
}

func (t *Run) GetName() string {
	return t.Name
}

type DockerPush struct {
	Name     string
	Username string
	Password string
	Repo     string
	Vars     Vars
}

func (t *DockerPush) GetName() string {
	return t.Name
}

type DeployCF struct {
	Name     string
	Api      string
	Space    string
	Org      string
	Username string
	Password string
	Manifest string
	Vars     Vars
}

func (t *DeployCF) GetName() string {
	return t.Name
}

type Vars map[string]string

var allTasks = map[string]func() task{
	"run":         func() task { return new(Run) },
	"docker-push": func() task { return new(DockerPush) },
	"deploy-cf":   func() task { return new(DeployCF) },
}
