package main

import (
    "io"
    "os"
    "os/exec"
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

