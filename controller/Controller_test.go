package controller

import (
	"testing"

	"bytes"

	"os"

	"github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setup() (*Controller, *bytes.Buffer, *bytes.Buffer) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")
	return &Controller{
		FileSystem:   afero.Afero{Fs: afero.NewMemMapFs()},
		RootDir:      "/root",
		OutputWriter: stdOut,
		ErrorWriter:  stdErr,
	}, stdOut, stdErr
}

func TestNoManifest(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewMissingFile("/root/.halfpipe.io").Error()
	assert.Contains(t, stdErr.String(), expectedError)
}

func TestManifestParseError(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ctrl.FileSystem.WriteFile("/root/.halfpipe.io", []byte("^&*(^&*"), os.ModePerm)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewParseError("").Error()
	assert.Contains(t, stdErr.String(), expectedError)
}

func TestEmptyManifest(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	ctrl.FileSystem.WriteFile("/root/.halfpipe.io", []byte{}, os.ModePerm)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())
	expectedError := model.NewParseError("/root/.halfpipe.io is empty").Error()
	assert.Contains(t, stdErr.String(), expectedError)
}

// ignored cos permissions don't work in mem fs: https://github.com/spf13/afero/issues/150
//func TestPermissionsManifest(t *testing.T) {
//	ctrl, stdOut, stdErr := setup()
//	ctrl.FileSystem.WriteFile("/root/.halfpipe.io", []byte{}, 0)
//	ok := ctrl.Run()
//
//	assert.False(t, ok)
//	assert.Empty(t, stdOut.String())
//  expectedError := model.NewParseError("").Error()
//	//assert.Contains(t, stdErr.String(), expectedError)
//}

func TestValidManifest(t *testing.T) {
	ctrl, stdOut, stdErr := setup()

	yaml := `
team: foo
repo: 
  uri: git@github.com/foo/bar.git
tasks:
- name: run
  script: foo
`
	ctrl.FileSystem.WriteFile("/root/.halfpipe.io", []byte(yaml), os.ModePerm)
	ok := ctrl.Run()

	assert.True(t, ok)
	assert.Empty(t, stdErr.String())
	assert.Contains(t, stdOut.String(), "Good job")
}
