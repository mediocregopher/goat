# goat

       _))
      /* \     _~
      `;'\\__-' \_     A simple go dependency manager
         | )  _ \ \
        / / ``   w w
       w w

* Unobtrusive dependency management
* Allows projects to be located anywhere, regardless of GOPATH
* Existing projects can be switched to goat without changing any code

See the [introduction][intro] for more details.

# Usage

**Pulling dependencies for an existing project:**

```bash
cd project/
goat deps
```

That's it. You just learned all the command-line stuff for goat.

**Creating a new project:**

```bash
mkdir newproject # Can be anywhere on the filesystem, regardless of GOPATH
cd newproject
vim .go.yaml
```

in .go.yaml put:

```yaml
---
path: github.com/user/newproject
```

**Adding a dependency to a project:**

Change .go.yaml to read:

```yaml
---
path: github.com/user/newproject
deps:
    - loc: gopkg.in/yaml.v1 # the same path you would use for go get
```

**Adding a dependency with a specific version:**

Change .go.yaml to read:

```yaml
---
path: github.com/user/newproject
deps:
    - loc: gopkg.in/yaml.v1 # the same path you would use for go get

    - loc: https://github.com/mediocregopher/flagconfig.git
      type: git
      ref: v0.4.2 # A tag in this case, could be commit hash or branch name
      path: github.com/mediocregopher/flagconfig
```

Run `goat deps` to automatically fetch both of your dependencies into your
project. Running `goat build` or `goat run` within your project will
automatically use any dependencies goat has fetched.

# More Usage

See the [tutorial][tutorial] for a basic use case for goat with more explanation
and a real project. After that check out the [.go.yaml file][projfile] for more
details on what kind of features goat has for dependency management. There are
also some [special features][special] that don't really fit in anywhere else
that might be useful to know about.

# Installation

To use goat you can either get a pre-compiled binary or build it yourself. Once
you get the binary I recommend renaming it as `go` (`alias go=goat`),
so that `goat` gets used whenever you use the `go` utility. Don't worry, all
`go` commands inside and outside of goat projects will behave the same way.

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

# Copyrights

Goat ASCII Art (c) 1997 ejm, Creative Commons

[intro]: /docs/introduction.md
[tutorial]: /docs/tut.md
[projfile]: /docs/projfile.md
[special]: /docs/special.md
