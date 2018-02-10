package linter

import (
	"path"

	. "github.com/robwhitby/halfpipe-cli/model"
	"github.com/spf13/afero"
)

func requiredFiles(man Manifest) (files []RequiredFile) {
	for _, t := range man.Tasks {
		switch task := t.(type) {
		case Run:
			files = append(files, RequiredFile{
				Path:       task.Script,
				Executable: true,
			})
		}
	}
	return
}

func LintFiles(man Manifest, rootDir string, fs afero.Afero) (errs []error) {
	for _, file := range requiredFiles(man) {
		if err := CheckFile(file, rootDir, fs); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func CheckFile(file RequiredFile, rootDir string, fs afero.Afero) error {
	absPath := path.Join(rootDir, file.Path)

	if exists, _ := fs.Exists(absPath); !exists {
		return NewFileError(absPath, "does not exist")
	}

	info, err := fs.Stat(absPath)
	if err != nil {
		return NewFileError(absPath, "cannot be read")
	}

	if !info.Mode().IsRegular() {
		return NewFileError(absPath, "is not a regular file")
	}

	if info.Size() == 0 {
		return NewFileError(absPath, "is empty")
	}

	if file.Executable && info.Mode()&0111 == 0 {
		return NewFileError(absPath, "is not executable")
	}

	return nil
}
