package deps

import (
	. "github.com/mediocregopher/goat/src/goat/common"
	"github.com/mediocregopher/goat/src/goat/exec"
	"fmt"
)

func GoGet(genv *GoatEnv, dep *Dependency) error {
	fmt.Println("go", "get", dep.Location)
	return exec.PipedCmd("go", "get", dep.Location)
}
