package exec

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var trimstr = "\n\t\b\r "

// TrimmedCmd returns a command's output on stdout and stderr as
// a string and error object. Before returning both stdout and stderr
// have whitespace trimmed off both ends. Stderr will be nil if it was
// empty
func TrimmedCmd(cmdstr string, args ...string) (string, error) {
	cmd := exec.Command(cmdstr, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	bout, err := ioutil.ReadAll(stdout)
	strout := strings.Trim(string(bout), trimstr)
	if err != nil {
		return strout, err
	}

	berr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return strout, err
	}

	cmd.Wait()

	if len(berr) == 0 {
		return strout, nil
	}

	strerr := strings.Trim(string(berr), trimstr)
	return strout, errors.New(strerr)
}

// PipedCmd pipes a command's out/err to this process', and returns
// a channel which gives an err if anything went wrong, or returns
// nil when the command completes
func PipedCmd(cmdstr string, args ...string) error {
	cmd := exec.Command(cmdstr, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	err = cmd.Start()
	if err != nil {
		return err
	}

	cmd.Wait()
	return nil
}
