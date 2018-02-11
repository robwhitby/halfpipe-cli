package model

import (
	"io"

	"github.com/spf13/afero"
)

type Config struct {
	FileSystem    afero.Fs
	Options       Options
	OutputWriter  io.Writer
	ErrorWriter   io.Writer
	SecretChecker SecretChecker
	Version       string
}
