package env

import (
	. "github.com/mediocregopher/goat/src/goat/common"
	"os"
	"path/filepath"
)

func pathExists (path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Setup takes a goat environment and makes sure it has all the proper
// directories created inside of it. This includes the lib/ directory, and if
// it's specified the Path loopback in the lib/ directory
func Setup(genv *GoatEnv) error {
	var err error

	// Make the lib directory if it doesn't exist
	if !pathExists(genv.ProjRootLib) {
		err = os.Mkdir(genv.ProjRootLib, 0755)
		if err != nil {
			return err
		}
	}
	
	if genv.Path != "" {
		loopbackPath := filepath.Join(genv.ProjRootLib,"src",genv.Path)
		if !pathExists(loopbackPath) {
			loopbackDir := filepath.Dir(loopbackPath)
			if err = os.MkdirAll(loopbackDir,0755); err != nil {
				return err
			} else if err = os.Symlink(genv.ProjRoot,loopbackPath); err != nil {
				return err
			}
		}
	}

	return nil
}
