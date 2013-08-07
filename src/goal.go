package main

import (
    "os"
    "os/exec"
    "path/filepath"
    "fmt"
    "bufio"
    "io"
    "errors"
    "syscall"
)

// pipedCmd pipes a command's out/err to this process', and returns
// a channel which gives an err if anything went wrong, or returns
// nil when the command completes
func pipedCmd(cmdstr string, args... string) chan error {
    cmd := exec.Command(cmdstr,args...)
    rch := make(chan error,1)

    stdout,err := cmd.StdoutPipe()
    if err != nil {
        rch <- err
        return rch
    }
    stderr,err := cmd.StderrPipe()
    if err != nil {
        rch <- err
        return rch
    }

    go io.Copy(os.Stdout, stdout)
    go io.Copy(os.Stderr, stderr)

    err = cmd.Start()
    if err != nil {
        rch <- err
        return rch
    }

    go func(){
        cmd.Wait()
        rch <- nil
    }()

    return rch
}

// Returns the directory name of the parent that contains the
// Goalfile
func findGoalfile(dir string) (string,error) {
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
        } else if n.Name() == "Goalfile" {
            return dir,nil
        }
    }

    parent := filepath.Dir(dir)
    if dir == parent {
        return "",errors.New("Goalfile not found")
    }

    return findGoalfile(parent)
}

func fatal(err error) {
    fmt.Println(err)
    os.Exit(1)
}

func main() {

    cwd,err := os.Getwd()
    if err != nil {
        fatal(err)
    }

    projroot,err := findGoalfile(cwd)
    if err != nil {
        fatal(err)
    }

    goalfile := filepath.Join(projroot,"Goalfile")
    projrootlib := filepath.Join(projroot,"lib")
    if _, err := os.Stat(projrootlib); os.IsNotExist(err) {
        err = os.Mkdir(projrootlib,0755)
        if err != nil {
            fatal(err)
        }
    }

    gopath,_ := syscall.Getenv("GOPATH")
    err = syscall.Setenv("GOPATH",projrootlib+":"+projroot+":"+gopath)
    if err != nil {
        fatal(err)
    }
    
    gfh,err := os.Open(goalfile)
    defer gfh.Close()
    if err != nil {
        fatal(err)
    }

    gfhbuf := bufio.NewReader(gfh)
    for {
        dep,err := readDependency(gfhbuf)
        if err == io.EOF {
            break
        } else if err != nil {
            fatal(err)
        } else {
            if dep.Name == "" {
                dep.Name = dep.Location
            }

            fun,ok := typemap[dep.Type]
            if !ok {
                fatal(errors.New("Unknown dependency type: "+dep.Type))
            }
            err = fun(projrootlib,dep)
            if err != nil {
                fatal(err)
            }
        }
    }

    //fmt.Println( <- pipedCmd("go",os.Args[1:]...) )
}
