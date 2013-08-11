package dep

import (
	"bufio"
	"errors"
	"fmt"
	. "goat"
	"goat/env"
	"goat/exec"
	"io"
	"os"
	"path/filepath"
)

func FetchDependencies(genv *GoatEnv) error {
	gfh, err := os.Open(genv.Goatfile)
	defer gfh.Close()
	if err != nil {
		return err
	}

	if _, err := os.Stat(genv.ProjRootLib); os.IsNotExist(err) {
		err = os.Mkdir(genv.ProjRootLib, 0755)
		if err != nil {
			return err
		}
	}

	gfhbuf := bufio.NewReader(gfh)
	for {
		dep, err := ReadDependency(gfhbuf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			if dep.Path == "" {
				dep.Path = dep.Location
			}

			fun, ok := typemap[dep.Type]
			if !ok {
				return errors.New("Unknown dependency type: " + dep.Type)
			}
			err = fun(genv, dep)
			if err != nil {
				return err
			}

			depprojroot := filepath.Join(genv.ProjRootLib, "src", dep.Path)
			issubproj, err := env.IsProjRoot(depprojroot)
			if err != nil {
				return err
			} else if issubproj {
				depgenv := env.SetupGoatEnv(depprojroot)
				env.ChrootEnv(depgenv, genv.ProjRoot)
				err = FetchDependencies(depgenv)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

type typefunc func(*GoatEnv, *Dependency) error

var typemap = map[string]typefunc{
	"":    goget,
	"get": goget,
	"git": git,
}

func goget(genv *GoatEnv, dep *Dependency) error {
	fmt.Println("go", "get", dep.Location)
	return exec.PipedCmd("go", "get", dep.Location)
}

func git(genv *GoatEnv, dep *Dependency) error {
	localloc := filepath.Join(genv.ProjRootLib, "src", dep.Path)

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
