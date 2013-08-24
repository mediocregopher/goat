package env

import (
	. "github.com/mediocregopher/goat/src/goat/common"
	"encoding/json"
	"fmt"
	"io"
	"bytes"
)

func deprecationLog(genv *GoatEnv) {
	fmt.Printf("Goatfile at %s is using an old, deprecated format. It'll still work, but it may not in future versions\n",genv.Goatfile)
}

// UnmarshalGoat will try to fill the goat environment based on the given
// contents of a Goatfile. It will first try the most recent Goatfile format,
// then work backwards in time until it finds a compatible version. If no
// compatible version is found then and error is returned.
func UnmarshalGoat(genvraw []byte, genv *GoatEnv) error {
	var err error

	
	if err = json.Unmarshal(genvraw,genv); err == nil && (genv.Path != "" ||
	genv.Dependencies != nil) {
		return nil	
	}

	// First version of goat had Goatfiles just being a stream of json
	// dependency objects. Not sure why I thought that was a good idea.
	buf := bytes.NewBuffer(genvraw)
	deps := make([]Dependency,0,8)
	for {
		depraw,err := buf.ReadBytes('}')
		if err == io.EOF {
			if len(deps) < 1 {
				break
			} else {
				genv.Dependencies = deps
				deprecationLog(genv)
				return nil
			}
		} else if err != nil {
			return err
		} else {
			dep := Dependency{}
			err = json.Unmarshal(depraw,&dep)
			if err != nil {
				return err
			} else if dep.Location == "" {
				break
			}

			deps = append(deps,dep)
		}
	}

	// First version of goat had empty files being legitimate, since it just
	// meant no dependencies. Check for this as well.
	trimmed := bytes.Trim(genvraw," \n\r\t\b")
	if len(trimmed) == 0 {
		deprecationLog(genv)
		return nil
	}

	return fmt.Errorf("Goatfile at %s is not properly formatted",genv.Goatfile)
}
