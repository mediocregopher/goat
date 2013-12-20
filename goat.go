package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mediocregopher/goat/env"
	"github.com/mediocregopher/goat/exec"
	"os"
	"syscall"
)

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func printGhelp() {
	fmt.Printf(
		`Goat is a command-line wrapper for go which handles dependency
management in a sane way. Check the goat docs at github.com/mediocregopher/goat
for a more in-depth overview.

Usage:

    %s command [arguments]

The commands are:

    deps    Read the .go.yaml file for this project and set up dependencies in
            the dependencies folder specified (default ".deps"). Recursively
            download dependencies wherever a .go.yaml file is encountered

    freeze  Scan through the list of dependencies in .go.yaml and update their
            references to the current ones installed.

    gen     Generate an initial .go.yaml file for this project based on imports
            done in this package and any imported packages (except for 
            sub-packages which also have .go.yaml files).

    ghelp   Show this dialog

All other commands are passed through to the go binary on your system. Try '%s
help' for its available commands
`, os.Args[0], os.Args[0])
	os.Exit(0)
}

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	projroot, err := env.FindProjRoot(cwd)
	var genv *env.GoatEnv
	if err == nil {
		genv, err = env.NewGoatEnv(projroot)
		if err != nil {
			fatal(err)
		}

		if err = genv.PrependToGoPath(); err != nil {
			fatal(err)
		}

		if err = genv.Setup(); err != nil {
			fatal(err)
		}
	}

	args := os.Args[1:]
	if len(args) < 1 {
		printGhelp()
	}

	switch args[0] {
	case "gen":
		var err error
		genFlags := flag.NewFlagSet("goat gen", flag.ExitOnError)
		forceUpdate := genFlags.Bool("update", false, "Update an existing .go.yaml file")
		freeze := genFlags.Bool("freeze", false, "Lock in the current git or hg revision of dependencies")
		genFlags.Parse(args[1:])
		if genv != nil {
			if !*forceUpdate {
				fatal(errors.New(".go.yaml file already exists, use -update to force update"))
			}
			err = genv.UpdateDependencies()
		} else {
			pwd, err := os.Getwd()
			if err != nil {
				fatal(err)
			}
			genv, err = env.GenGoatEnv(pwd)
			if err != nil {
				fatal(err)
			}
			err = genv.PrependToGoPath()
		}
		if err != nil {
			fatal(err)
		}
		if *freeze {
			err = genv.FreezeDependencies()
			if err != nil {
				fatal(err)
			}
		}
		err = genv.WriteConf()
		if err != nil {
			fatal(err)
		}
		fmt.Println(".go.yaml has been generated. Add a proper path setting if you haven't already")
	case "freeze":
		if genv != nil {
			err := genv.FreezeDependencies()
			if err != nil {
				fatal(err)
			}
		} else {
			fatal(errors.New(".go.yaml file not found on current path"))
		}
		err = genv.WriteConf()
		if err != nil {
			fatal(err)
		}
		fmt.Println(".go.yaml has been updated")
	case "deps":
		if genv != nil {
			err := genv.FetchDependencies(genv.AbsDepDir())
			if err != nil {
				fatal(err)
			}
		} else {
			fatal(errors.New(".go.yaml file not found on current path"))
		}
	case "ghelp":
		printGhelp()
	default:
		if actualgo, ok := ActualGo(); ok {
			exec.PipedCmd(actualgo, args...)
		} else {
			newargs := make([]string, len(args)+1)
			copy(newargs[1:], args)
			newargs[0] = "go"
			exec.PipedCmd("/usr/bin/env", newargs...)
		}
	}
}

// ActualGo returns the GOAT_ACTUALGO environment variable contents, and whether
// or not the variable was actually set
func ActualGo() (string, bool) {
	return syscall.Getenv("GOAT_ACTUALGO")
}
