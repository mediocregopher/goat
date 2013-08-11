# goat

A simple, json-based go dependency manager

Goat handles recursive, versioned, dependency management for go projects in an
unobtrusive way. To switch to using goat you probably won't even have to change
any code.

# The problem

There are two problems that goat aims to solve:

* `go get` does not allow for specifying versions of a library.

* `go get` does not have a easy of way of sandboxing your project's development
  environment. You can either muck up your global environment with dependencies
  or mess with your `GOPATH` everytime you want to develop for that project.
  Others that want to work on your project will have to do the same.

* Other dependency managers are strange and have weird command line arguments
  that I don't feel like learning.

# The solution

* The root of a goat project has a `Goatfile` which specifies a dependency's
  location, name, and version control reference if applicable. It is formatted
  using super-simple json objects, each having at most four fields.

* All dependencies are placed in a `lib` directory at the root of your project.
  goat will automatically look for a `Goatfile` in your current working
  directory or one of its parents, and call that the project root. For the rest
  of the command's duration your GOPATH will have the project root prepended to
  it, and `<projroot>/lib` prepended to that. This has many useful properties,
  most notably that your dependencies are sandboxed inside your code, but are
  still usable exactly as they would have been if they were global.

* Goat is a wrapper around the go command line utility. It adds one new command,
  all other commands are passed straight through to the normal go binary. This
  command is `goat deps`, and it retrieves all dependencies listed in your
  `Goatfile` and puts them into a folder called `lib` in your project. If any of
  those dependencies have `Goatfiles` then those are processed and put in your
  project's `lib` folder as well (this is done recursively).

# Installation

To use goat you can either get a pre-compiled binary or build it yourself. Once
you get the binary I recommend renaming it to `go` and putting it on your `PATH`
so that it gets used whenever you use the `go` utility. Don't worry, unless you
are in a directory tree with a `Goatfile` or use one of goat's like two commands
nothing will be different.

## Pre-built

Check the releases tab on github, the latest release will have pre-compiled
binaries for various systems, choose the one that applies to you.

## Build it yourself

To build goat yourself make sure you have `go` installed (go figure).

```bash
git clone https://github.com/mediocregopher/goat.git
cd goat
make
```

The binaries will be found in the `bin` directory.

# Usage

See the [tutorial](/docs/tut.md) for a basic use case for goat. After that
check out the [Goatfile](/docs/goatfile.md) for more details on what kind of
features goat has for dependency management.
