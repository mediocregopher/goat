package deps

import (
	"fmt"
	. "github.com/mediocregopher/goat/common"
	"github.com/mediocregopher/goat/exec"
	"os"
	"path/filepath"
)

func Git(depdir string, dep *Dependency) error {
	localloc := filepath.Join(depdir, "src", dep.Path)

	fmt.Println("git", "clone", dep.Location, localloc)
	err := exec.PipedCmd("git", "clone", dep.Location, localloc)
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

	fmt.Println("git", "fetch", "-pv", "--all")
	err = exec.PipedCmd("git", "fetch", "-pv", "--all")
	if err != nil {
		return err
	}

	if dep.Reference == "" {
		dep.Reference = "master"
	}
	fmt.Println("git", "checkout", dep.Reference)
	err = exec.PipedCmd("git", "checkout", dep.Reference)
	if err != nil {
		return err
	}

	fmt.Println("git", "clean", "-f", "-d")
	err = exec.PipedCmd("git", "clean", "-f", "-d")

	return err

}

func GitInfo(path string) (string, string, error) {
	origcwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	err = os.Chdir(path)
	if err != nil {
		return "", "", err
	}
	defer os.Chdir(origcwd)

	ref, err := exec.TrimmedCmd("git", "rev-parse", "--verify", "HEAD")
	if err != nil {
		return "", "", err
	}
	url, err := exec.TrimmedCmd("git", "config", "--get", "remote.origin.url")
	if err != nil {
		return "", "", err
	}

	return url, ref, nil
}
