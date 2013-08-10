package main

import (
    "encoding/json"
    "bufio"
)

func readDependency(buf *bufio.Reader) (*Dependency,error) {
    b,err := buf.ReadBytes('}')
    if err != nil {
        return nil,err
    }

    var dep Dependency
    err = json.Unmarshal(b,&dep)
    return &dep,err
}
