package goat

import (
    "os"
    "path/filepath"
    "errors"
)

// FindGoatFile eturns the directory name of the parent that contains the
// Goatfile
func FindGoatfile(dir string) (string,error) {
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

    return FindGoatfile(parent)
}
