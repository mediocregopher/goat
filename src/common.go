package main

type Dependency struct {
    Location string `json:"loc"`
    Path string `json:"path"`
    Type string `json:"type"`
    Reference string `json:"reference"`
}

type GoalEnv struct {
    ProjRoot string
    ProjRootLib string
    Goalfile string
}
