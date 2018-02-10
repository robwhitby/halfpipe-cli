package controller

import (
	"testing"

	"bytes"

	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	rootDir      = "/root/"
	manifestPath = rootDir + manifestFilename
)

func setup() (Controller, *bytes.Buffer, *bytes.Buffer) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")
	return Controller{
		FileSystem:   afero.Afero{Fs: afero.NewMemMapFs()},
		RootDir:      rootDir,
		OutputWriter: stdOut,
		ErrorWriter:  stdErr,
	}, stdOut, stdErr
}

func TestNoManifest(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewFileError(manifestPath, "does not exist")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestParseError(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ctrl.FileSystem.WriteFile(manifestPath, []byte("^&*(^&*"), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewParseError("")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestLintError(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ctrl.FileSystem.WriteFile(manifestPath, []byte("foo: bar"), 0777)
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
	ctrl.FileSystem.WriteFile(manifestPath, []byte(yaml), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewFileError("/root/build.sh", "does not exist")
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
  script: ./foo/bar.sh
  image: bar
`
	ctrl.FileSystem.WriteFile(manifestPath, []byte(yaml), 0777)
	ctrl.FileSystem.WriteFile("/root/foo/bar.sh", []byte("x"), 0777)
	ok := ctrl.Run()

	assert.True(t, ok)
	assert.Empty(t, stdErr.String())
	assert.Contains(t, stdOut.String(), "Good job")
}
