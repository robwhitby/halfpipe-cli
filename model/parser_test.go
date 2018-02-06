package model

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidYaml(t *testing.T) {
	g := NewGomegaWithT(t)

	man, err := Parse("team: my team")
	expected := &Manifest{Team: "my team"}
	g.Expect(man, err).To(Equal(expected))
}

func TestInvalidYaml(t *testing.T) {
	g := NewGomegaWithT(t)

	_, err := Parse("team : { foo")
	g.Expect(err).To(HaveOccurred())
}

func TestRepo(t *testing.T) {
	g := NewGomegaWithT(t)

	man, err := Parse("repo: { uri: myuri, private_key: mypk }")
	expected := &Manifest{
		Repo: Repo{
			Uri:        "myuri",
			PrivateKey: "mypk",
		},
	}
	g.Expect(man, err).To(Equal(expected))
}

func TestRunTask(t *testing.T) {
	g := NewGomegaWithT(t)

	man, err := Parse("tasks: [{ name: run, image: alpine, script: build.sh, vars: { FOO: Foo, BAR: Bar } }]")
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
	g.Expect(man, err).To(Equal(expected))
}

func TestMultipleTasks(t *testing.T) {
	g := NewGomegaWithT(t)

	man, err := Parse("tasks: [{ name: run, image: img, script: build.sh }, { name: docker-push, username: bob }, { name: run }, { name: deploy-cf, org: foo }]")
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
	g.Expect(man, err).To(Equal(expected))
}

func TestInvalidTask(t *testing.T) {
	g := NewGomegaWithT(t)

	_, err := Parse("tasks: [{ name: unknown, foo: bar }]")
	g.Expect(err).To(HaveOccurred())
}

func TestReportMultipleInvalidTasks(t *testing.T) {
	g := NewGomegaWithT(t)

	_, err := Parse("tasks: [{ name: unknown, foo: bar }, { name: run, image: alpine, script: build.sh }, { notname: foo }]")
	g.Expect(err.Messages).To(HaveLen(2))
}
