package deps

import (
	"fmt"
	. "github.com/mediocregopher/goat/src/goat/common"
	"github.com/mediocregopher/goat/src/goat/exec"
)

func GoGet(genv *GoatEnv, dep *Dependency) error {
	fmt.Println("go", "get", dep.Location)
	return exec.PipedCmd("go", "get", dep.Location)
}
