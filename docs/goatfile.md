# Goatfile

The `Goatfile` is used by goat both to find the root directory of the project
and to get information about all of the project's dependencies. Using the
`Goatfile` to set the root of the project is covered in the
[tutorial](/docs/tut.md), in this doc we'll be looking at the format of the
`Goatfile` and how it can be used.

# Basic Format

Here's a basic `Goatfile`:

```json
{
    "path":"github.com/username/myproject",
    "deps":[
        {
          "loc":"code.google.com/p/go.example/newmath",
          "type":"goget"
        },
        {
          "loc":"https://github.com:mediocregopher/goat.git",
          "type":"git",
          "reference":"master",
          "path":"github.com/mediocregopher/goat"
        }
    ]
}
```

The outer object describes attributes about the project. One of those is `deps`,
which describes the project's dependencies, each as a json object.

# path

The `path` field is required. It can either be a simple name for the project if
you don't plan on hosting it anywhere, or a full uri based on where it is hosted
(like in the example above). This uri can then be used to import your submodules
within eachother. For example if there was a folder `myproject/bettermath` in
your project you could import it with:

```
import "github.com/username/myproject/bettermath"
```

Note that your project this project wouldn't have to actually exist in the
`github.com/username` directory tree, it could be anywhere. Goat takes care of
the `GOPATH` for you.

# deps

The `deps` field is a list of dependency objects that this project requires.
There are four possible dependency fields.

## loc

This is the only required field. It is the actual url that goat wil use fetch
your project. This will be passed into whatever dependency fetcher you are using
(go get, git, mercurial, etc...).

## type

By default this field is `goget`, but currently `git` and `hg` are also
supported. This tells goat how to fetch the dependency. For example if `type` is
set to `goget` then `go get <path>` is used, while if it's `git` then `git clone
<path>` will be used.

## reference

The option is only valid for `type`s that have some kind of version reference
system (for example `git` does, `goget` doesn't). Here are the defaults for the
supported types:

* git: `master`
* hg:  `tip`

## path

The actual directory structure the dependency will be put into inside `lib`
folder. By default this is set to whatever `loc` is set to, which works for
`goget` dependencies but not so well for others where the `loc` is an actual url
(like `git`). In effect, `path` is the string you want to use for importing this
dependency inside your project.

An interesting effect is that you can have the same project being depended on
using different names. For instance if you want to have "version1" and
"version2" of "someproject" in your project you could have the `path`s be
`someproject-v1` and `someproject-v2` with the `reference`s set accordingly.
