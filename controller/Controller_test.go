package controller

import (
	"testing"

	"bytes"

	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setup() (controller, *bytes.Buffer, *bytes.Buffer) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")
	return NewController(afero.NewMemMapFs(), "/root", stdOut, stdErr), stdOut, stdErr
}

func TestNoManifest(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewFileError(manifestFilename, "does not exist")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestParseError(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ctrl.FileSystem.WriteFile(manifestFilename, []byte("^&*(^&*"), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewParseError("")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestLintError(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ctrl.FileSystem.WriteFile(manifestFilename, []byte("foo: bar"), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewMissingField("team")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestRequiredFileError(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	yaml := `
team: foo
repo: 
  uri: git@github.com/foo/bar.git
tasks:
- name: run
  script: ./build.sh
  image: bar
`
	ctrl.FileSystem.WriteFile(manifestFilename, []byte(yaml), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewFileError("./build.sh", "does not exist")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestValidManifest(t *testing.T) {
	ctrl, stdOut, stdErr := setup()

	yaml := `
team: foo
repo: 
  uri: git@github.com/foo/bar.git
tasks:
- name: run
  script: foo/bar.sh
  image: bar
`
	ctrl.FileSystem.WriteFile(manifestFilename, []byte(yaml), 0777)
	ctrl.FileSystem.WriteFile("foo/bar.sh", []byte("x"), 0777)
	ok := ctrl.Run()

	assert.True(t, ok)
	assert.Empty(t, stdErr.String())
	assert.Contains(t, stdOut.String(), "Good job")
}
