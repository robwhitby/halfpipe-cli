package linter

import (
	"testing"

	. "github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const rootDir = "/some/root/dir/"

func fs() afero.Afero {
	return afero.Afero{Fs: afero.NewMemMapFs()}
}

func TestFile_NotExists(t *testing.T) {
	fs := fs()
	err := CheckFile(RequiredFile{Path: ".halfpipe.io"}, rootDir, fs)

	assert.Equal(t, NewFileError(rootDir+".halfpipe.io", "does not exist"), err)
}

func TestFile_Empty(t *testing.T) {
	fs := fs()
	fs.WriteFile(rootDir+".halfpipe.io", []byte{}, 0777)

	err := CheckFile(RequiredFile{Path: ".halfpipe.io"}, rootDir, fs)
	assert.Equal(t, NewFileError(rootDir+".halfpipe.io", "is empty"), err)
}

func TestFile_IsDirectory(t *testing.T) {
	fs := fs()
	fs.Mkdir(rootDir+"build", 0777)

	err := CheckFile(RequiredFile{Path: "build"}, rootDir, fs)
	assert.Equal(t, NewFileError(rootDir+"build", "is not a regular file"), err)
}

func TestFile_NotExecutable(t *testing.T) {
	fs := fs()
	fs.WriteFile(rootDir+"build.sh", []byte("go test"), 0400)

	err := CheckFile(RequiredFile{Path: "build.sh", Executable: true}, rootDir, fs)
	assert.Equal(t, NewFileError(rootDir+"build.sh", "is not executable"), err)

	err = CheckFile(RequiredFile{Path: "build.sh", Executable: false}, rootDir, fs)
	assert.Nil(t, err)
}

func TestFile_Happy(t *testing.T) {
	fs := fs()
	fs.WriteFile(rootDir+".halfpipe.io", []byte("foo"), 0700)
	err := CheckFile(RequiredFile{Path: ".halfpipe.io", Executable: true}, rootDir, fs)

	assert.Nil(t, err)
}

var manifest = Manifest{
	Team: "ee",
	Repo: Repo{Uri: "http://github.com/foo/bar.git"},
	Tasks: []Task{Run{
		Script: "./build.sh",
		Image:  "alpine",
	}},
}

func TestRequiredFiles_RunTaskScript(t *testing.T) {
	files := requiredFiles(manifest)

	expected := []RequiredFile{{
		Path:       "./build.sh",
		Executable: true,
	}}

	assert.Equal(t, expected, files)
}
