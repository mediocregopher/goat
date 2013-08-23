package env

import (
	"errors"
	. "goat/common"
	"os"
	"path/filepath"
	"goat/env/deps"
)

type typefunc func(*GoatEnv, *Dependency) error

var typemap = map[string]typefunc{
	"":    deps.GoGet,
	"get": deps.GoGet,
	"git": deps.Git,
	"hg": deps.Hg,
}

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
		if IsProjRoot(depprojroot) {
			depgenv,err := SetupGoatEnv(depprojroot)
			if err != nil {
				return err
			}
			ChrootEnv(depgenv, genv.ProjRoot)
			err = FetchDependencies(depgenv)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


