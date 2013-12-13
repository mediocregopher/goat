package env

import (
	"errors"
	"fmt"
	. "github.com/mediocregopher/goat/common"
	"github.com/mediocregopher/goat/env/deps"
	"path/filepath"
)

type typefunc func(string, *Dependency) error

var typemap = map[string]typefunc{
	"":    deps.GoGet,
	"get": deps.GoGet,
	"git": deps.Git,
	"hg":  deps.Hg,
}

func header(c string, strs ...interface{}) {
	fmt.Printf("\n")
	for i := 0; i < 80; i++ {
		fmt.Printf(c)
	}
	fmt.Printf("\n")

	fmt.Println(strs...)

	for i := 0; i < 80; i++ {
		fmt.Printf(c)
	}
	fmt.Printf("\n")
	fmt.Printf("\n")
}

// FetchDependencies goes and retrieves the dependencies for the GoatEnv. If the
// dependencies have any goat project files in them, this will fetch the
// dependencies listed therein as well. All dependencies are placed in the root
// project's lib directory, INCLUDING any sub-dependencies. This is done on
// purpose!
func (genv *GoatEnv) FetchDependencies(depdir string) error {
	var err error

	if len(genv.Dependencies) > 0 {
		header("#", "Downloading dependencies listed in", genv.AbsProjFile())

		for i := range genv.Dependencies {
			dep := &genv.Dependencies[i]

			header("=", "Retrieving dependency at:", dep.Location)

			if dep.Path == "" {
				dep.Path = dep.Location
			}

			fun, ok := typemap[dep.Type]
			if !ok {
				return errors.New("Unknown dependency type: " + dep.Type)
			}
			err = fun(depdir, dep)
			if err != nil {
				return err
			}

			depprojroot := filepath.Join(depdir, "src", dep.Path)

			if IsProjRoot(depprojroot) {
				header("-", "Reading", depprojroot, "'s dependencies")
				depgenv, err := NewGoatEnv(depprojroot)
				if err != nil {
					return err
				}

				err = depgenv.FetchDependencies(depdir)
				if err != nil {
					return err
				}
			} else {
				header("-", "No "+PROJFILE+" found in", depprojroot)
			}
		}

		header("#", "Done downloading dependencies for", genv.AbsProjFile())
	} else {
		header("-", "No dependencies listed in", genv.AbsProjFile())
	}

	return nil
}
