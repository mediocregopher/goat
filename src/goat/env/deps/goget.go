package deps

import (
	. "goat/common"
	"goat/exec"
	"fmt"
)

func GoGet(genv *GoatEnv, dep *Dependency) error {
	fmt.Println("go", "get", dep.Location)
	return exec.PipedCmd("go", "get", dep.Location)
}
