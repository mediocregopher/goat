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
)

var srcDirs []string

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
	for _, importPath := range pkg.Imports {
		subPkg, err := tryImport(importPath)
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

func tryImport(path string) (*build.Package, error) {
	for _, srcDir := range srcDirs {
		pkg, err := build.Import(path, srcDir, 0)
		if err == nil {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find %s in any of %v", path, srcDirs)
}

func init() {
	srcDirs = build.Default.SrcDirs()
}

func FreezeImport(path string) (RepoType, string, string, error) {
	pkg, err := tryImport(path)
	if err != nil {
		return NoRepo, "", "", err
	}
	repoDir, repoType := repoRoot(pkg)

	var url, ref string
	switch repoType {
	case NoRepo:
		return NoRepo, "", "", fmt.Errorf("No repo found for %s", pkg.ImportPath)
	case Git:
		url, ref, err = deps.GitInfo(repoDir)
	case Hg:
		url, ref, err = deps.HgInfo(repoDir)
	}
	if err != nil {
		return NoRepo, "", "", err
	}

	return repoType, url, ref, nil
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

func repoRoot(pkg *build.Package) (string, RepoType) {
	for p := pkg.ImportPath; p != "."; p = path.Dir(p) {
		dir := path.Join(pkg.SrcRoot, p)
		if _, err := os.Stat(path.Join(dir, ".git")); !os.IsNotExist(err) {
			return dir, Git
		}
		if _, err := os.Stat(path.Join(dir, ".hg")); !os.IsNotExist(err) {
			return dir, Hg
		}
	}
	return "", NoRepo
}
