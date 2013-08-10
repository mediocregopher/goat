package dep

import (
    "encoding/json"
    "bufio"
    . "goat"
)

// ReadDependency reads in a dependency json object from a buffered reader and
// returns it
func ReadDependency(buf *bufio.Reader) (*Dependency,error) {
    b,err := buf.ReadBytes('}')
    if err != nil {
        return nil,err
    }

    var dep Dependency
    err = json.Unmarshal(b,&dep)
    return &dep,err
}
