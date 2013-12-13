package env

import (
	. "github.com/mediocregopher/goat/common"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"syscall"
)

var PROJFILE = ".go.yaml"

// unmarshal takes in some bytes and tries to decode them into a GoatEnv
// structure
func unmarshal(genvraw []byte) (*GoatEnv, error) {
	genv := GoatEnv{}
	if err := goyaml.Unmarshal(genvraw, &genv); err != nil {
		return nil, err
	}
	return &genv, nil
}

type GoatEnv struct {
	// ProjRoot is the absolute path to the project root in the current
	// environment
	ProjRoot string

	// Path is the path that the project will be using for its own internal
	// import statements, and consequently what other projects depending on this
	// one will use as well.
	Path string `yaml:"path"`

	// Dependencies are the dependencies listed in the project's project file
	Dependencies []Dependency `yaml:"deps"`
}

// NewGoatEnv takes in the directory where a goat project file can be found,
// creates the GoatEnv struct based on that file, and returns it
func NewGoatEnv(projroot string) (*GoatEnv, error) {
	projfile := filepath.Join(projroot, PROJFILE)
	b, err := ioutil.ReadFile(projfile)
	if err != nil {
		return nil, err
	}
	genv, err := unmarshal(b)
	if err != nil {
		return nil, err
	}

	genv.ProjRoot = projroot
	return genv, nil
}

// AbsDepDir is the absolute path to the dependency directory inside the goat
// data directory
func (genv *GoatEnv) AbsDepDir() string {
	return filepath.Join(genv.AbsGoatDir(), "deps")
}

// AbsGoatDir is the absolute path to the goat data directory
func (genv *GoatEnv) AbsGoatDir() string {
	return filepath.Join(genv.ProjRoot, ".goat")
}

// AbsProjFile is the absolute path to the goat project file for this
// environment
func (genv *GoatEnv) AbsProjFile() string {
	return filepath.Join(genv.ProjRoot, PROJFILE)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Setup makes sure the goat env has all the proper directories created inside
// of it. This includes the lib/ directory, and if it's specified the Path
// loopback in the lib/ directory
func (genv *GoatEnv) Setup() error {
	var err error

	// Make the dep directory if it doesn't exist
	depdir := genv.AbsDepDir()
	if !pathExists(depdir) {
		err = os.MkdirAll(depdir, 0755)
		if err != nil {
			return err
		}
	}

	if genv.Path != "" {
		loopbackPath := filepath.Join(depdir, "src", genv.Path)
		if !pathExists(loopbackPath) {
			loopbackDir := filepath.Dir(loopbackPath)
			if err = os.MkdirAll(loopbackDir, 0755); err != nil {
				return err
			} else if err = os.Symlink(genv.ProjRoot, loopbackPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (genv *GoatEnv) PrependToGoPath() error {
	gopath, _ := syscall.Getenv("GOPATH")
	return syscall.Setenv("GOPATH", genv.AbsDepDir()+":"+gopath)
}
