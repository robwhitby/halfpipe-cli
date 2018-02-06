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

	man, err := Parse("tasks: [{ name: run, image: alpine, script: build.sh }]")
	expected := &Manifest{
		Tasks: []task{
			&Run{
				Name:   "run",
				Image:  "alpine",
				Script: "build.sh",
			},
		},
	}
	g.Expect(man, err).To(Equal(expected))
}

func TestParseRunMultipleTasks(t *testing.T) {
	g := NewGomegaWithT(t)

	man, err := Parse("tasks: [{ name: run, image: alpine, script: build.sh }, { name: docker, username: bob }, { name: run, image: alpine2 }]")
	expected := &Manifest{
		Tasks: []task{
			&Run{
				Name:   "run",
				Image:  "alpine",
				Script: "build.sh",
			},
			&Docker{
				Name:     "docker",
				Username: "bob",
			},
			&Run{
				Name:  "run",
				Image: "alpine2",
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
