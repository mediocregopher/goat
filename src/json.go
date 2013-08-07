package main

import (
    "encoding/json"
    "bufio"
)

type dependency struct {
    Location string `json:"loc"`
    Name string `json:"path"`
    Type string `json:"type"`
    Reference string `json:"reference"`
}

func readDependency(buf *bufio.Reader) (*dependency,error) {
    b,err := buf.ReadBytes('}')
    if err != nil {
        return nil,err
    }

    var dep dependency
    err = json.Unmarshal(b,&dep)
    return &dep,err
}
