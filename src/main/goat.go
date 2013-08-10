package main

import (
    "os"
    "fmt"
    "goat"
    "goat/env"
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

    projroot,err := env.FindGoatfile(cwd)
    if err != nil {
        fatal(err)
    }

    genv := env.SetupGoatEnv(projroot)
    if err = env.EnvPrependProj(genv); err != nil {
        fatal(err)
    }

    goat.PipedCmd("/bin/go",os.Args[1:]...)
}
