package main

import (
    "os"
    "fmt"
    . "goat"
)

func fatal(err error) {
    fmt.Println(err)
    os.Exit(1)
}

func main() {

    cwd,err := os.Getwd()
    if err != nil {
        fatal(err)
    }

    projroot,err := FindGoatfile(cwd)
    if err != nil {
        fatal(err)
    }

    genv := SetupGoatEnv(projroot)
    if err = EnvPrependProj(genv); err != nil {
        fatal(err)
    }

    PipedCmd("/bin/go",os.Args[1:]...)
}
