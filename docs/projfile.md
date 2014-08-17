# .go.yaml file

The `.go.yaml` is used by goat both to find the root directory of the project
and to get information about all of the project's dependencies. Using the
`.go.yaml` to set the root of the project is covered in the
[tutorial](/docs/tut.md), in this doc we'll be looking at the format of the
`.go.yaml` and how it can be used.

If you've never used yaml before, [wikipedia][yaml] has a fairly good outline of
it including some examples and explanations of the various data-types.

# Basic Format

Here's a basic `.go.yaml`:

```yaml
---
path: github.com/username/myproject
deps:
  - loc: code.google.com/p/go.example/newmath
    type: get

  - loc: https://github.com:mediocregopher/goat.git
    type: git
    ref: master
    path: github.com/mediocregopher/goat
```

The outer object describes attributes about the project. It has the keys `path`
and `deps`.

## path

The `path` field is required. It can either be a simple name for the project if
you don't plan on hosting it anywhere, or a full uri based on where it is hosted
(like in the example above). This uri can then be used to import your submodules
within eachother. For example if there was a folder in your project called
`bettermath` and you had `path: github.com/username/myproject`, you could import
it in another submodule in your project with:

```
import "github.com/username/myproject/bettermath"
```

Note that your project wouldn't have to actually exist in the
`github.com/username` directory tree, it could be anywhere. Goat takes care of
the `GOPATH` for you.

## deps

The `deps` field is a list of dependency objects that this project requires.
There are four possible dependency fields.

### loc

This is the only required field. It is the actual url that goat wil use fetch
your project. This will be passed into whatever dependency fetcher you are using
(go get, git, mercurial, etc...).

### type

By default this field is `get`, but currently `git` and `hg` are also
supported. This tells goat how to fetch the dependency. For example if `type` is
set to `get` then `go get <path>` is used, while if it's `git` then `git clone
<path>` will be used.

### ref

The option is only valid for `type`s that have some kind of version reference
system (for example `git` does, `get` doesn't). Here are the defaults for the
supported types:

* git: `master`
* hg:  `tip`

### path

The actual directory structure the dependency will be put into inside the `depdir`
folder. By default this is set to whatever `loc` is set to, which works for
`get` dependencies but not so well for others where the `loc` is an actual url
(like `git`). In effect, `path` is the string you want to use for importing this
dependency inside your project.

[yaml]: http://en.wikipedia.org/wiki/YAML
