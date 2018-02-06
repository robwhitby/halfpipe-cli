package model

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

func Parse(manifestYaml string) (*Manifest, *Failures) {
	man := new(Manifest)
	failures := new(Failures)
	if err := yaml.Unmarshal([]byte(manifestYaml), &man); err != nil {
		failures.Messages = append(failures.Messages, err.Error())
		return nil, failures
	}

	var rawTasks struct {
		Tasks []json.RawMessage
	}
	if err := yaml.Unmarshal([]byte(manifestYaml), &rawTasks); err != nil {
		failures.Messages = append(failures.Messages, err.Error())
		return nil, failures
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
		failures.Messages = append(failures.Messages, fmt.Sprintf("task %v is invalid", i+1))
	}
	if failures.IsEmpty() {
		return man, nil
	}
	return man, failures
}

func unmarshalTask(t task, taskName string, raw json.RawMessage) (error, bool) {
	if err := json.Unmarshal(raw, t); err != nil {
		return err, false
	}
	return nil, t.GetName() == taskName
}
