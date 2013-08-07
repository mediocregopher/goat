package main

import (
    "fmt"
)

type typefunc func(string,*dependency)error

var typemap = map[string]typefunc{
    "": goget,
    "get": goget,
    //"git": git,
}

func goget(projrootlib string, dep *dependency) error {
    fmt.Println("go","get",dep.Location)
    return pipedCmd("go","get",dep.Location)
}

//func git(projrootlib string, dep *dependency) error {
//    fmt.Println("git","clone",dep.Location,dep.Name)
//    err := <- pipedCmd("git","clone",dep.Location,dep.Name)
//    if err != nil {
//        return err
//    }
//
//}
