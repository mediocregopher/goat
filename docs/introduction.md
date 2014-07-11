# Introduction

This doc gives a brief introduction to what problems goat solves, and how it
solves them.

# The problem

There are two problems that goat aims to solve:

* `go get` does not allow for specifying versions of a library.

* `go get` does not have an easy of way of sandboxing your project's development
  environment. You can either muck up your global environment with dependencies
  or mess with your `GOPATH` everytime you want to develop for that project.
  Others that want to work on your project will have to do the same.

* Other dependency managers are strange and have weird command line arguments
  that I don't feel like learning.

# The solution

* The root of a goat project has a `.go.yaml` file which specifies a
  dependency's location, name, and version control reference if applicable. It
  is formatted using super-simple yaml objects, each having at most four fields.

* All dependencies are placed in a `.goat/deps` directory at the root of your
  project.  goat will automatically look for a `.go.yaml` in your current
  working directory or one of its parents, and call that the project root. For
  the rest of the command's duration your GOPATH will have
  `<projroot>/.goat/deps` prepended to it. This has many useful properties, most
  notably that your dependencies are sandboxed inside your code, but are still
  usable exactly as they would have been if they were global.

* Goat is a wrapper around the go command line utility. It adds one new command,
  all other commands are passed straight through to the normal go binary. This
  command is `goat deps`, and it retrieves all dependencies listed in your
  `.go.yaml` and puts them into a folder called `.goat/deps` in your project. If
  any of those dependencies have `.go.yaml` files then those are processed and
  put in your project's `.goat/deps` folder as well (this is done recursively).

