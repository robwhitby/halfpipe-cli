package model

import (
	"encoding/json"

	"fmt"

	"github.com/ghodss/yaml"
)

func Parse(manifestYaml string) (man *Manifest, errs []error) {
	addError := func(e error) {
		errs = append(errs, e)
	}

	if err := yaml.Unmarshal([]byte(manifestYaml), &man); err != nil {
		addError(err)
		return
	}

	var rawTasks struct {
		Tasks []json.RawMessage
	}
	if err := yaml.Unmarshal([]byte(manifestYaml), &rawTasks); err != nil {
		addError(err)
		return
	}

	for i, rawTask := range rawTasks.Tasks {
		done := false
		for taskName, taskFunc := range allTasks {
			newTask := taskFunc()
			if _, ok := unmarshalTask(newTask, taskName, rawTask); ok {
				man.Tasks = append(man.Tasks, newTask)
				done = true
				continue
			}
		}
		if done {
			continue
		}
		addError(NewInvalidField(fmt.Sprintf("task %v", i+1), "unknown task definition"))
	}
	return
}

func unmarshalTask(t task, taskName string, raw json.RawMessage) (error, bool) {
	if err := json.Unmarshal(raw, t); err != nil {
		return err, false
	}
	return nil, t.GetName() == taskName
}
