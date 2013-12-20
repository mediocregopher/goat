package env

import (
	"fmt"
	. "github.com/mediocregopher/goat/common"
	"github.com/mediocregopher/goat/env/deps"
	"go/build"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"syscall"
)

var (
	srcDirs     []string
	srcDirsOnce sync.Once
)

func getImports(path string) ([]string, error) {
	imports, err := getImportsMap(path)
	if err != nil {
		return nil, err
	}
	packages := make([]string, 0, len(imports))
	for k, _ := range imports {
		packages = append(packages, k)
	}
	sort.Strings(packages)
	return packages, nil
}

func getImportsMap(path string) (map[string]Dependency, error) {
	pkg, err := build.ImportDir(path, 0)
	if err != nil {
		return nil, err
	}
	imports := make(map[string]Dependency)
	walkImports(path, imports, pkg)
	return imports, nil
}

func walkImports(basePath string, imports map[string]Dependency, pkg *build.Package) {
	srcDirsOnce.Do(initSrcDirs)

	for _, importPath := range pkg.Imports {
		subPkg, err := tryImport(importPath, 0)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Skip stdlib packages and sub-packages
		if subPkg.SrcRoot == srcDirs[0] {
			continue
		}

		// Leave sub-packages out of the dependency list, but maybe still scan their deps
		if !strings.HasPrefix(subPkg.Dir, basePath) {
			imports[importPath] = Dependency{Location: importPath}
		}
		// Don't recurse into packages which have their own .go.yaml file
		if IsProjRoot(subPkg.Dir) {
			fmt.Printf("Not scanning %s since it has its own .go.yaml\n", subPkg.ImportPath)
			continue
		}
		walkImports(basePath, imports, subPkg)
	}
}

func tryImport(path string, mode build.ImportMode) (*build.Package, error) {
	srcDirsOnce.Do(initSrcDirs)
	for _, srcDir := range srcDirs {
		pkg, err := build.Import(path, srcDir, mode)
		if err == nil {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find %s in any of %v", path, srcDirs)
}

// Defer getting the list of srcDirs until after path prepend is done
func initSrcDirs() {
	gopath, _ := syscall.Getenv("GOPATH")
	build.Default.GOPATH = gopath
	srcDirs = build.Default.SrcDirs()
}

func FreezeImport(dep Dependency) (Dependency, error) {
	if dep.Path == "" {
		dep.Path = dep.Location
	}

	pkg, err := tryImport(dep.Path, build.FindOnly)
	if err != nil {
		return dep, err
	}
	importPath, repoDir, repoType := repoRoot(pkg)

	var url, ref string
	switch repoType {
	case NoRepo:
		return dep, fmt.Errorf("No repo found for %s", pkg.ImportPath)
	case Git:
		url, ref, err = deps.GitInfo(repoDir)
	case Hg:
		url, ref, err = deps.HgInfo(repoDir)
	}
	if err != nil {
		return dep, err
	}

	dep.Path = importPath
	dep.Location = url
	dep.Reference = ref

	return dep, nil
}

func repoRoot(pkg *build.Package) (string, string, RepoType) {
	for importPath := pkg.ImportPath; importPath != "."; importPath = path.Dir(importPath) {
		dir := path.Join(pkg.SrcRoot, importPath)
		if _, err := os.Stat(path.Join(dir, ".git")); !os.IsNotExist(err) {
			return importPath, dir, Git
		}
		if _, err := os.Stat(path.Join(dir, ".hg")); !os.IsNotExist(err) {
			return importPath, dir, Hg
		}
	}
	return "", "", NoRepo
}

type RepoType int

const (
	NoRepo RepoType = iota
	Git
	Hg
)

func (rt RepoType) String() string {
	switch rt {
	case NoRepo:
		return ""
	case Git:
		return "git"
	case Hg:
		return "hg"
	default:
		return ""
	}
}
