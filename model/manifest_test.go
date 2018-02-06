package model

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseValidYaml(t *testing.T) {
	g := NewGomegaWithT(t)

	man, err := Parse("team: my team")
	expected := &Manifest{Team: "my team"}
	g.Expect(man, err).To(Equal(expected))
}

func TestParseInvalidYaml(t *testing.T) {
	g := NewGomegaWithT(t)

	_, err := Parse("team : { foo")
	g.Expect(err).To(HaveOccurred())
}

func TestParseRepo(t *testing.T) {
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

func TestParseRunTask(t *testing.T) {
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

func TestParseRunMultipleTasks(t *testing.T) {
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

func TestParseInvalidTask(t *testing.T) {
	g := NewGomegaWithT(t)

	_, err := Parse("tasks: [{ name: unknown, foo: bar }]")
	g.Expect(err).To(HaveOccurred())
}

func TestParseReportMultipleInvalidTasks(t *testing.T) {
	g := NewGomegaWithT(t)

	_, err := Parse("tasks: [{ name: unknown, foo: bar }, { name: run, image: alpine, script: build.sh }, { notname: foo }]")
	g.Expect(err.Messages).To(HaveLen(2))
}
