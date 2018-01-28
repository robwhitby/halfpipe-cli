package manifest

import (
	"path"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type ManifestReader struct {
	fs *afero.Afero
}

func NewManifestReader(fs *afero.Afero) *ManifestReader {
	return &ManifestReader{
		fs: fs,
	}
}

func (m *ManifestReader) parseVars(vars interface{}) map[string]string {
	returnMap := map[string]string{}
	switch vars.(type) {
	case map[interface{}]interface{}:

		for k, v := range vars.(map[interface{}]interface{}) {
			switch v.(type) {
			case string:
				returnMap[k.(string)] = v.(string)
			case int:
				returnMap[k.(string)] = strconv.Itoa(v.(int))
			}

		}
	}
	return returnMap
}

func (m *ManifestReader) parseRunTask(task map[interface{}]interface{}) RunTask {
	runTask := RunTask{}
	mapstructure.Decode(task, &runTask)
	if runTask.Vars != nil {
		if vars, ok := task["vars"]; ok {
			runTask.Vars = m.parseVars(vars)
		}
	} else {
		runTask.Vars = make(map[string]string)
	}
	return runTask
}

func (m *ManifestReader) parseDeployTask(task map[interface{}]interface{}, team string) DeployTask {
	deployTask := DeployTask{}
	mapstructure.Decode(task, &deployTask)
	if deployTask.Org == "" {
		deployTask.Org = team
	}
	if deployTask.Password == "" {
		deployTask.Password = "((cf-credentials.password))"
	}
	if deployTask.Username == "" {
		deployTask.Username = "((cf-credentials.username))"
	}
	if deployTask.Manifest == "" {
		deployTask.Manifest = "manifest.yml"
	}
	if deployTask.Vars != nil {
		if vars, ok := task["vars"]; ok {
			deployTask.Vars = m.parseVars(vars)
		}
	} else {
		deployTask.Vars = make(map[string]string)
	}
	return deployTask
}

func (m *ManifestReader) parseDockerTask(task map[interface{}]interface{}) DockerTask {
	dockerTask := DockerTask{}
	mapstructure.Decode(task, &dockerTask)
	return dockerTask
}

func (m *ManifestReader) ParseManifest(repoRoot string) (Manifest, error) {
	halfPipePath := path.Join(repoRoot, ".halfpipe.io")
	manifest := Manifest{}
	exists, err := m.fs.Exists(halfPipePath)
	if err != nil {
		return manifest, err
	}
	if exists == false {
		return manifest, errors.Errorf("Manifest at path '%s' does not exist", halfPipePath)
	}

	bytes, err := m.fs.ReadFile(halfPipePath)
	if err != nil {
		return manifest, err
	}

	yaml.Unmarshal(bytes, &manifest)
	tasks := []Task{}
	for _, task := range manifest.Tasks {
		t := task.(map[interface{}]interface{})
		switch t["task"] {
		case "run":
			tasks = append(tasks, m.parseRunTask(t))
		case "docker":
			tasks = append(tasks, m.parseDockerTask(t))
		case "deploy":
			tasks = append(tasks, m.parseDeployTask(t, manifest.Team))
		case nil:
			return Manifest{}, errors.New("Task is missing 'task' key")
		default:
			return Manifest{}, errors.Errorf("Task '%s' is not supported", t["task"].(string))
		}
	}

	if len(tasks) > 0 {
		manifest.Tasks = tasks
	}
	return manifest, nil
}
