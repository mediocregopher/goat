package goat

import (
    "fmt"
    "os"
    "path/filepath"
)

type typefunc func(*GoatEnv,*Dependency)error

var TypeMap = map[string]typefunc{
    "": goget,
    "get": goget,
    "git": git,
}

func goget(genv *GoatEnv, dep *Dependency) error {
    fmt.Println("go","get",dep.Location)
    return PipedCmd("go","get",dep.Location)
}

func git(genv *GoatEnv, dep *Dependency) error {
    localloc := filepath.Join(genv.ProjRootLib,dep.Path)

    fmt.Println("git","clone",dep.Location,localloc)
    err := PipedCmd("git","clone",dep.Location,localloc)
    if err != nil {
        return err
    }

    origcwd,err := os.Getwd()
    if err != nil {
        return err
    }

    err = os.Chdir(localloc)
    if err != nil {
        return err
    }
    defer os.Chdir(origcwd)

    fmt.Println("git","fetch","-pv","--all")
    err = PipedCmd("git","fetch","-pv","--all")
    if err != nil {
        return err
    }

    if dep.Reference == "" {
        dep.Reference = "master"
    }
    fmt.Println("git","checkout",dep.Reference)
    err = PipedCmd("git","checkout",dep.Reference)
    if err != nil {
        return err
    }

    fmt.Println("git","clean","-f","-d")
    err = PipedCmd("git","clean","-f","-d")

    return err

}
