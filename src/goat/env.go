package goat

import (
    "syscall"
    "path/filepath"
)

// NewGoatEnv returns a new GoatEnv struct based on the directory passed in
func SetupGoatEnv(projroot string) *GoatEnv {

    goatfile := filepath.Join(projroot,GOATFILE)
    projrootlib := filepath.Join(projroot,"lib")
    //if _, err := os.Stat(projrootlib); os.IsNotExist(err) {
    //    err = os.Mkdir(projrootlib,0755)
    //    if err != nil {
    //        return nil,err
    //    }
    //}

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
