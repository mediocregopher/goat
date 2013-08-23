package dep

import (
	"errors"
	"fmt"
	. "goat/common"
	"goat/env"
	"goat/exec"
	"os"
	"path/filepath"
)

// FetchDependencies goes and retrieves the dependencies for a given GoatEnv. If
// the dependencies have any Goatfile's in them, this will fetch the
// dependencies listed therein as well. All dependencies are placed in the root
// project's lib directory, INCLUDING any sub-dependencies. This is done on
// purpose!
func FetchDependencies(genv *GoatEnv) error {
	var err error

	if _, err = os.Stat(genv.ProjRootLib); os.IsNotExist(err) {
		err = os.Mkdir(genv.ProjRootLib, 0755)
		if err != nil {
			return err
		}
	}

	for i := range genv.Dependencies {
		dep := &genv.Dependencies[i]
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
		if env.IsProjRoot(depprojroot) {
			depgenv,err := env.SetupGoatEnv(depprojroot)
			if err != nil {
				return err
			}
			env.ChrootEnv(depgenv, genv.ProjRoot)
			err = FetchDependencies(depgenv)
			if err != nil {
				return err
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
	"hg": hg,
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

func hg(genv *GoatEnv, dep *Dependency) error {
	localloc := filepath.Join(genv.ProjRootLib, "src", dep.Path)

	fmt.Println("hg","clone",dep.Location,localloc)
	err := exec.PipedCmd("hg","clone",dep.Location,localloc)
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

	fmt.Println("hg","pull")
	err = exec.PipedCmd("hg","pull")
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
