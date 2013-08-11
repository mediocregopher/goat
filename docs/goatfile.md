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
  "loc":"code.google.com/p/go.example/newmath",
  "type":"goget"
}

{
  "loc":"https://github.com:mediocregopher/goat.git",
  "type":"git",
  "reference":"master",
  "path":"github.com/mediocregopher/goat"
}
```

There are two dependencies defined, each being a json object. goat isn't picky
about formatting as long as the json objects are valid.

## Fields

There are four fields that can be used to define a dependency:

* `loc`: This is the only required field. It is the actual url that goat wil use
         fetch your project. This will be passed into whatever dependency
         fetcher you are using (go get, git, mercurial, etc...).

* `type`: By default this field is `goget`, but currently `git` is also
          supported. This tells goat how to fetch the dependency. For example if
          `type` is set to `goget` then `go get <path>` is used, while if it's
          `git` then `git clone <path>` will be used.

* `reference`: The option is only valid for `type`s that have some kind of
               version reference system (for example `git` does, `goget`
               doesn't). By default this is set to `master` for git, but it can
               be set to another branch name, a tag, or a commit hash.

* `path`: The actual directory structure the dependency will be put into inside
          `lib` folder. By default this is set to whatever `path` is set to,
          which works for `goget` dependencies but not so well for others where
          the `loc` is an actual url (like `git`). In effect, `path` is the
          string you want to use for importing this dependency inside your
          project. An interesting effect is that you can have the same project
          being depended on using different names. For instance if you want to
          have "version1" and "version2" of "someproject" in your project you
          could have the `path`s be `someproject-v1` and `someproject-v2` with
          the `reference`s set accordingly.
