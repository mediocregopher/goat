package main

import (
    "fmt"
    "os"
    "path/filepath"
)

type typefunc func(*GoalEnv,*Dependency)error

var typemap = map[string]typefunc{
    "": goget,
    "get": goget,
    "git": git,
}

func goget(genv *GoalEnv, dep *Dependency) error {
    fmt.Println("go","get",dep.Location)
    return pipedCmd("go","get",dep.Location)
}

func git(genv *GoalEnv, dep *Dependency) error {
    localloc := filepath.Join(genv.ProjRootLib,dep.Path)

    fmt.Println("git","clone",dep.Location,localloc)
    err := pipedCmd("git","clone",dep.Location,localloc)
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
    err = pipedCmd("git","fetch","-pv","--all")
    if err != nil {
        return err
    }

    if dep.Reference == "" {
        dep.Reference = "master"
    }
    fmt.Println("git","checkout",dep.Reference)
    err = pipedCmd("git","checkout",dep.Reference)
    if err != nil {
        return err
    }

    fmt.Println("git","clean","-f","-d")
    err = pipedCmd("git","clean","-f","-d")

    return err

}
