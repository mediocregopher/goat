package env

import (
	"path/filepath"
	"errors"
	"os"
)

// FindGoatFile returns the directory name of the parent that contains the
// Goatfile
func FindGoatfile(dir string) (string, error) {

	if IsProjRoot(dir) {
		return dir, nil
	}

	parent := filepath.Dir(dir)
	if dir == parent {
		return "", errors.New("Goatfile not found")
	}

	return FindGoatfile(parent)
}

// IsProjRoot returns whether or not a particular directory is the project
// root for a goat project (aka, whether or not it has a goat file)
func IsProjRoot(dir string) bool {
	goatfile := filepath.Join(dir, GOATFILE)
	if _, err := os.Stat(goatfile); os.IsNotExist(err) {
		return false
	}
	return true
}
