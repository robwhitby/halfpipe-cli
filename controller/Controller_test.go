package controller

import (
	"testing"

	"bytes"

	"os"

	. "github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const root = "/root/"

var opts = Options{
	ShowVersion: false,
	Args: Args{
		Dir: root,
	},
}

func setup() (controller, *bytes.Buffer, *bytes.Buffer) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	//only 'valid.secret' exists
	secretChecker := func(s string) bool { return s == "valid.secret" }

	ctrl := NewController(afero.NewMemMapFs(), opts, stdOut, stdErr, secretChecker)
	ctrl.FileSystem.Mkdir(root, 0777)
	return ctrl, stdOut, stdErr
}

func TestNoManifest(t *testing.T) {
	ctrl, _, stdErr := setup()
	ok := ctrl.Run()

	assert.False(t, ok)
	expectedError := NewFileError(manifestFilename, "does not exist")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestParseError(t *testing.T) {
	ctrl, _, stdErr := setup()
	afero.WriteFile(ctrl.FileSystem, root+manifestFilename, []byte("^&*(^&*"), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	expectedError := NewParseError("")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestLintError(t *testing.T) {
	ctrl, _, stdErr := setup()
	afero.WriteFile(ctrl.FileSystem, root+manifestFilename, []byte("foo: bar"), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	expectedError := NewMissingField("team")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestRequiredFileError(t *testing.T) {
	ctrl, _, stdErr := setup()
	yaml := `
team: foo
repo: 
  uri: git@github.com/foo/bar.git
tasks:
- name: run
  script: ./build.sh
  image: bar
`
	afero.WriteFile(ctrl.FileSystem, root+manifestFilename, []byte(yaml), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	expectedError := NewFileError("./build.sh", "does not exist")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestManifestRequiredSecretError(t *testing.T) {
	ctrl, _, stdErr := setup()
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
	afero.WriteFile(ctrl.FileSystem, root+manifestFilename, []byte(yaml), 0777)
	afero.WriteFile(ctrl.FileSystem, root+"build.sh", []byte("x"), 0777)
	ok := ctrl.Run()

	assert.False(t, ok)
	expectedError := NewMissingSecret("path.to.key")
	assert.Contains(t, stdErr.String(), expectedError.Error())

	unexpected := NewMissingSecret("valid.secret")
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
	afero.WriteFile(ctrl.FileSystem, root+manifestFilename, []byte(yaml), 0777)
	afero.WriteFile(ctrl.FileSystem, "/root/build.sh", []byte("x"), 0777)
	ok := ctrl.Run()

	assert.True(t, ok)
	assert.Empty(t, stdErr.String())
	assert.Contains(t, stdOut.String(), "Good job")
}

func TestController_ChecksRootDir(t *testing.T) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")
	secretChecker := func(s string) bool { return false }
	opts := opts
	opts.Args.Dir = "/invalid/root"
	ctrl := NewController(afero.NewMemMapFs(), opts, stdOut, stdErr, secretChecker)
	ok := ctrl.Run()

	assert.False(t, ok)
	expectedError := NewFileError("/invalid/root", "is not a directory")
	assert.Contains(t, stdErr.String(), expectedError.Error())
}

func TestAbsDirectory_Abs(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	fs.MkdirAll("/some/dir", 0777)

	dir, _ := absDir("/some/dir/", fs)
	assert.Equal(t, "/some/dir", dir)

	dir, _ = absDir("/some/dir/../dir", fs)
	assert.Equal(t, "/some/dir", dir)
}

func TestAbsDirectory_Relative(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	pwd, _ := os.Getwd()
	fs.MkdirAll(pwd, 0777)

	dir, _ := absDir(".", fs)
	assert.Equal(t, pwd, dir)

	dir, _ = absDir("", fs)
	assert.Equal(t, pwd, dir)

	fs.MkdirAll(pwd+"/foo", 0777)

	dir, _ = absDir("foo", fs)
	assert.Equal(t, pwd+"/foo", dir)

	dir, _ = absDir("./foo/", fs)
	assert.Equal(t, pwd+"/foo", dir)
}

func TestAbsDirectory_Errors(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	pwd, _ := os.Getwd()
	fs.MkdirAll(pwd, 0777)

	fileError := NewFileError("missing", "is not a directory")

	_, err := absDir("missing", fs)
	assert.Equal(t, fileError, err)

	fs.WriteFile("/file", []byte{}, 0777)
	_, err = absDir("/file", fs)
	assert.IsType(t, fileError, err)
}

func TestOption_Version(t *testing.T) {
	ctrl, stdOut, _ := setup()
	ctrl.Options.ShowVersion = true
	ok := ctrl.Run()

	assert.True(t, ok)
	assert.Equal(t, versionText()+"\n", stdOut.String())
}
