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

	//only 'valid.secret' exists
	secretChecker := func(s string) bool { return s == "valid.secret" }

	return NewController(afero.NewMemMapFs(), "/root", stdOut, stdErr, secretChecker), stdOut, stdErr
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

func TestManifestRequiredSecretError(t *testing.T) {
	ctrl, stdOut, stdErr := setup()
	yaml := `
team: foo
repo: 
  uri: git@github.com/foo/bar.git
tasks:
- name: run
  script: build.sh
  image: bar
  vars:
    badsecret: ((path.to.key))
    goodsecret: ((valid.secret))
`
	ctrl.FileSystem.WriteFile(manifestFilename, []byte(yaml), 0777)
	ctrl.FileSystem.WriteFile("build.sh", []byte("x"), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	assert.Empty(t, stdOut.String())

	expectedError := model.NewMissingSecret("path.to.key")
	assert.Contains(t, stdErr.String(), expectedError.Error())

	unexpected := model.NewMissingSecret("valid.secret")
	assert.NotContains(t, stdErr.String(), unexpected.Error())
}

func TestValidManifest(t *testing.T) {
	ctrl, stdOut, stdErr := setup()

	yaml := `
team: foo
repo: 
  uri: git@github.com/foo/bar.git
tasks:
- name: run
  script: build.sh
  image: bar
  vars:
    secret: ((valid.secret))
`
	ctrl.FileSystem.WriteFile(manifestFilename, []byte(yaml), 0777)
	ctrl.FileSystem.WriteFile("build.sh", []byte("x"), 0777)
	ok := ctrl.Run()

	assert.True(t, ok)
	assert.Empty(t, stdErr.String())
	assert.Contains(t, stdOut.String(), "Good job")
}
