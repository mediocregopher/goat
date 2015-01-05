package deps

import (
	"fmt"

	. "github.com/mediocregopher/goat/common"
	"github.com/mediocregopher/goat/exec"
)

func GoGet(depdir string, dep *Dependency) error {
	fmt.Println("go", "get", dep.Location)
	return exec.PipedCmd("go", "get", dep.Location)
}
