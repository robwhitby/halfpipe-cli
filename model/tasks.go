package model

type task interface {
	GetName() string
}

type Run struct {
	Name     string
	Script   string
	Username string
	Image    string
}

func (t *Run) GetName() string {
	return t.Name
}

type Docker struct {
	Name     string
	Username string
	Password string
	Repo     string
}

func (t *Docker) GetName() string {
	return t.Name
}

var allTasks = map[string]func() task{
	"run":    func() task { return new(Run) },
	"docker": func() task { return new(Docker) },
}
