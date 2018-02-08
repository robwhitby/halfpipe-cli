package controller

import (
	"testing"

	"bytes"

	"fmt"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")
	fs := afero.NewMemMapFs()

	ctrl := &Controller{
		FileSystem:   afero.Afero{Fs: fs},
		RootDir:      "/root",
		OutputWriter: stdOut,
		ErrorWriter:  stdErr,
	}

	ok := ctrl.Run()

	fmt.Println("stdout:")
	fmt.Println(stdOut.String())

	fmt.Println("stderr:")
	fmt.Println(stdErr.String())

	assert.False(t, ok)

}
