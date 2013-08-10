package goat

var GOATFILE = "Goatfile"

type Dependency struct {
    Location string `json:"loc"`
    Path string `json:"path"`
    Type string `json:"type"`
    Reference string `json:"reference"`
}

type GoatEnv struct {
    ProjRoot string
    ProjRootLib string
    Goatfile string
}
