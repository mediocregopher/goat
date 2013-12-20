package env

import (
	"errors"
	"fmt"
	. "github.com/mediocregopher/goat/common"
	"github.com/mediocregopher/goat/env/deps"
	"path/filepath"
	"sort"
)

// depList implements sort.Interface for lists of Dependencies.
type depList []Dependency

func (deps depList) Len() int {
	return len(deps)
}

func (deps depList) Less(i, j int) bool {
	return deps[i].Location < deps[j].Location
}

func (deps depList) Swap(i, j int) {
	deps[i], deps[j] = deps[j], deps[i]
}

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

// DetectDependencies looks at the package referenced by genv and
// fills in the Dependency list based on what is found.
func (genv *GoatEnv) DetectDependencies() error {
	imports, err := getImports(genv.ProjRoot)
	if err != nil {
		return err
	}
	for _, imp := range imports {
		genv.Dependencies = append(genv.Dependencies, Dependency{Location: imp})
	}
	return nil
}

// UpdateDependencies looks at the package referenced by genv and
// updates the Dependency list based on what is found. Existing
// dependencies are not overwritten.
func (genv *GoatEnv) UpdateDependencies() error {
	imports, err := getImportsMap(genv.ProjRoot)
	if err != nil {
		return err
	}
	for _, dep := range genv.Dependencies {
		imports[dep.Location] = dep
	}
	genv.Dependencies = make(depList, 0, len(imports))
	for _, dep := range imports {
		genv.Dependencies = append(genv.Dependencies, dep)
	}
	sort.Sort(depList(genv.Dependencies))
	return nil
}

func (genv *GoatEnv) FreezeDependencies() error {
	combined := make(map[string]Dependency)
	for _, dep := range genv.Dependencies {
		frozenDep, err := FreezeImport(dep)
		if err != nil {
			fmt.Println(err)
		}
		combined[frozenDep.Location] = frozenDep
	}
	genv.Dependencies = make(depList, 0, len(combined))
	for _, dep := range combined {
		genv.Dependencies = append(genv.Dependencies, dep)
	}
	sort.Sort(depList(genv.Dependencies))
	return nil
}
