package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidYaml(t *testing.T) {
	man, errs := Parse("team: my team")
	expected := &Manifest{Team: "my team"}

	assert.Nil(t, errs)
	assert.Equal(t, man, expected)
}

func TestInvalidYaml(t *testing.T) {
	_, errs := Parse("team : { foo")

	assert.Equal(t, len(errs), 1)
}

func TestRepo(t *testing.T) {
	man, errs := Parse("repo: { uri: myuri, private_key: mypk }")
	expected := &Manifest{
		Repo: Repo{
			Uri:        "myuri",
			PrivateKey: "mypk",
		},
	}

	assert.Nil(t, errs)
	assert.Equal(t, man, expected)
}

func TestRunTask(t *testing.T) {
	man, errs := Parse("tasks: [{ name: run, image: alpine, script: build.sh, vars: { FOO: Foo, BAR: Bar } }]")
	expected := &Manifest{
		Tasks: []task{
			&Run{
				Name:   "run",
				Image:  "alpine",
				Script: "build.sh",
				Vars: Vars{
					"FOO": "Foo",
					"BAR": "Bar",
				},
			},
		},
	}

	assert.Nil(t, errs)
	assert.Equal(t, man, expected)
}

func TestMultipleTasks(t *testing.T) {
	man, errs := Parse("tasks: [{ name: run, image: img, script: build.sh }, { name: docker-push, username: bob }, { name: run }, { name: deploy-cf, org: foo }]")
	expected := &Manifest{
		Tasks: []task{
			&Run{
				Name:   "run",
				Image:  "img",
				Script: "build.sh",
			},
			&DockerPush{
				Name:     "docker-push",
				Username: "bob",
			},
			&Run{
				Name: "run",
			},
			&DeployCF{
				Name: "deploy-cf",
				Org:  "foo",
			},
		},
	}

	assert.Nil(t, errs)
	assert.Equal(t, man, expected)
}

func TestInvalidTask(t *testing.T) {
	_, errs := Parse("tasks: [{ name: unknown, foo: bar }]")

	assert.Equal(t, len(errs), 1)
}

func TestReportMultipleInvalidTasks(t *testing.T) {
	_, errs := Parse("tasks: [{ name: unknown, foo: bar }, { name: run, image: alpine, script: build.sh }, { notname: foo }]")

	assert.Equal(t, len(errs), 2)
	assert.IsType(t, errs[0], &invalidField{})
	assert.IsType(t, errs[1], &invalidField{})
}
