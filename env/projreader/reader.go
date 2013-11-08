package projreader

import (
	"launchpad.net/goyaml"
	"github.com/mediocregopher/goat/common"
	"io/ioutil"
)

// UnmarshalFile reads the data out of a file and puts it through Unmarshal
func UnmarshalFile(file string) (*common.GoatEnv, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return Unmarshal(b)
}

// Unmarshal takes in some bytes and tries to decode them into a GoatEnv
// structure
func Unmarshal(genvraw []byte) (*common.GoatEnv, error) {
	var genv *common.GoatEnv
	if err := goyaml.Unmarshal(genvraw, genv); err != nil {
		return nil, err	
	}
	return genv, nil
}
