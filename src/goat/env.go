package goat

import (
    "syscall"
    "path/filepath"
    "errors"
    "os"
)

// FindGoatFile returns the directory name of the parent that contains the
// Goatfile
func FindGoatfile(dir string) (string,error) {

    isroot,err := IsProjRoot(dir)
    if err != nil {
        return "",err
    } else if isroot {
        return dir,nil
    }

    parent := filepath.Dir(dir)
    if dir == parent {
        return "",errors.New("Goatfile not found")
    }

    return FindGoatfile(parent)
}

// IsProjRoot returns whether or not a particular directory is the project
// root for a goat project (aka, whether or not it has a goat file)
func IsProjRoot(dir string) (bool,error) {
    dirh,err := os.Open(dir)
    defer dirh.Close()
    if err != nil {
        return false,err
    }

    dirnodes,err := dirh.Readdir(-1)
    if err != nil {
        return false,err
    }

    for _,n := range dirnodes {
        if n.IsDir() {
            continue
        } else if n.Name() == GOATFILE {
            return true,nil
        }
    }

    return false,nil
}

// NewGoatEnv returns a new GoatEnv struct based on the directory passed in
func SetupGoatEnv(projroot string) *GoatEnv {

    goatfile := filepath.Join(projroot,GOATFILE)
    projrootlib := filepath.Join(projroot,"lib")

    return &GoatEnv{ ProjRoot: projroot,
                     ProjRootLib: projrootlib,
                     Goatfile: goatfile }
}

func envPrepend(dir string) error {
    gopath,_ := syscall.Getenv("GOPATH")
    return syscall.Setenv("GOPATH",dir+":"+gopath)
}

// EnvPrependProj prepends a goat project's root and lib directories to the GOPATH
func EnvPrependProj(genv *GoatEnv) error {
    err := envPrepend(genv.ProjRoot)
    if err != nil {
        return err
    }

    return envPrepend(genv.ProjRootLib)
}
