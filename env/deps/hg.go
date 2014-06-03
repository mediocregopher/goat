package deps

import (
	"fmt"
	. "github.com/mediocregopher/goat/common"
	"github.com/mediocregopher/goat/exec"
	"os"
	"path/filepath"
)

func Hg(depdir string, dep *Dependency) error {
	localloc := filepath.Join(depdir, "src", dep.Path)

	if _, err := os.Stat(localloc); os.IsNotExist(err) {
		fmt.Println("hg", "clone", dep.Location, localloc)
		err := exec.PipedCmd("hg", "clone", dep.Location, localloc)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(localloc, "exists")
	}

	origcwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(localloc)
	if err != nil {
		return err
	}
	defer os.Chdir(origcwd)

	fmt.Println("hg", "pull")
	err = exec.PipedCmd("hg", "pull")
	if err != nil {
		return err
	}

	if dep.Reference == "" {
		dep.Reference = "tip"
	}
	fmt.Println("hg", "update", "-C", dep.Reference)
	err = exec.PipedCmd("hg", "update", "-C", dep.Reference)

	return err

}
