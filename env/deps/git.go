package deps

import (
	"fmt"
	. "github.com/mediocregopher/goat/common"
	"github.com/mediocregopher/goat/exec"
	"os"
	"path/filepath"
	"strings"
)

func Git(depdir string, dep *Dependency) error {
	localloc := filepath.Join(depdir, "src", dep.Path)

	if _, err := os.Stat(localloc); os.IsNotExist(err) {
		fmt.Println("git", "clone", dep.Location, localloc)
		err := exec.PipedCmd("git", "clone", dep.Location, localloc)
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

	fmt.Println("git", "fetch", "-pv", "--all")
	err = exec.PipedCmd("git", "fetch", "-pv", "--all")
	if err != nil {
		return err
	}

	if dep.Reference == "" {
		dep.Reference = "master"
	}
	if exists, err := originBranchExists(dep.Reference); err != nil {
		return err
	} else if exists {
		dep.Reference = "origin/" + dep.Reference
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

// If 'origin/branch' exists, return true.
func originBranchExists(branch string) (bool, error) {
	branches, err := exec.TrimmedCmd("git", "branch", "--remote")
	if err != nil {
		return false, err
	}
	for _, b := range strings.Split(branches, "\n") {
		b := strings.TrimSpace(b)
		if strings.HasSuffix(b, "origin/"+branch) {
			return true, nil
		}
	}
	return false, nil
}
