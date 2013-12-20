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

	fmt.Println("hg", "clone", dep.Location, localloc)
	err := exec.PipedCmd("hg", "clone", dep.Location, localloc)
	if err != nil {
		return err
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

func HgInfo(path string) (string, string, error) {
	origcwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	err = os.Chdir(path)
	if err != nil {
		return "", "", err
	}
	defer os.Chdir(origcwd)

	ref, err := exec.TrimmedCmd("hg", "parent", "--template", "{node}")
	if err != nil {
		return "", "", err
	}
	url, err := exec.TrimmedCmd("hg", "paths", "default")
	if err != nil {
		return "", "", err
	}
	return url, ref, nil
}
