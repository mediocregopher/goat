package main

import (
    "os"
    "path/filepath"
    "fmt"
    "bufio"
    "io"
    "errors"
    "syscall"
    . "goat"
)

// Returns the directory name of the parent that contains the
// Goatfile
func findGoatfile(dir string) (string,error) {
    dirh,err := os.Open(dir)
    defer dirh.Close()
    if err != nil {
        return "",err
    }

    dirnodes,err := dirh.Readdir(-1)
    if err != nil {
        return "",err
    }

    for _,n := range dirnodes {
        if n.IsDir() {
            continue
        } else if n.Name() == GOATFILE {
            return dir,nil
        }
    }

    parent := filepath.Dir(dir)
    if dir == parent {
        return "",errors.New("Goatfile not found")
    }

    return findGoatfile(parent)
}

func setupGoatEnv() (*GoatEnv,error) {
    cwd,err := os.Getwd()
    if err != nil {
        return nil,err
    }

    projroot,err := findGoatfile(cwd)
    if err != nil {
        return nil,err
    }

    goatfile := filepath.Join(projroot,GOATFILE)
    projrootlib := filepath.Join(projroot,"lib")
    if _, err := os.Stat(projrootlib); os.IsNotExist(err) {
        err = os.Mkdir(projrootlib,0755)
        if err != nil {
            return nil,err
        }
    }

    gopath,_ := syscall.Getenv("GOPATH")
    err = syscall.Setenv("GOPATH",projrootlib+":"+projroot+":"+gopath)
    if err != nil {
        return nil,err
    }

    return &GoatEnv{ ProjRoot: projroot,
                     ProjRootLib: projrootlib,
                     Goatfile: goatfile },nil
}

func fatal(err error) {
    fmt.Println(err)
    os.Exit(1)
}

func main() {

    genv,err := setupGoatEnv()
    if err != nil {
        fatal(err)
    }

    gfh,err := os.Open(genv.Goatfile)
    defer gfh.Close()
    if err != nil {
        fatal(err)
    }

    gfhbuf := bufio.NewReader(gfh)
    for {
        dep,err := ReadDependency(gfhbuf)
        if err == io.EOF {
            break
        } else if err != nil {
            fatal(err)
        } else {
            if dep.Path == "" {
                dep.Path = dep.Location
            }

            fun,ok := TypeMap[dep.Type]
            if !ok {
                fatal(errors.New("Unknown dependency type: "+dep.Type))
            }
            err = fun(genv,dep)
            if err != nil {
                fatal(err)
            }
        }
    }

    PipedCmd("/bin/go",os.Args[1:]...)
}
