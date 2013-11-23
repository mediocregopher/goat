package env

import (
	"errors"
	"os"
	"path/filepath"
)

// FindProjRoot returns the directory name of the parent that contains the
// .go.yaml file
func FindProjRoot(dir string) (string, error) {

	if IsProjRoot(dir) {
		return dir, nil
	}

	parent := filepath.Dir(dir)
	if dir == parent {
		return "", errors.New(".go.yaml not found")
	}

	return FindProjRoot(parent)
}

// IsProjRoot returns whether or not a particular directory is the project
// root for a goat project (aka, whether or not it has a goat file)
func IsProjRoot(dir string) bool {
	goatfile := filepath.Join(dir, PROJFILE)
	if _, err := os.Stat(goatfile); os.IsNotExist(err) {
		return false
	}
	return true
}
