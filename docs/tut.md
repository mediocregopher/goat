# Hello world with goat

This tutorial will walk you through the steps of setting up a simple goat
project and importing some dependencies using the `Goatfile`. The first step is
setting up the goat binary, which you can find instructions on in the main
[README.md](/). I will assume you have a binary called `goat` available on your
`PATH`.

## Reference

We will be building a project structure in this tutorial. I'll try to be as
explicit as possible when writing this, but in case you get lost here is what
the project will look like at the end of this tutorial:

```
/goatproject
    Goatfile
    /lib (contains some stuff)
    /src
        /foo
            foo.go
        main.go
```

Not too complicated (I hope!)

## Initialization and Goatfile

The first step is to create a new directory. This directory will be the root of
the rest of your project and can be located anywhere you want. I'm going to put
mine in `/tmp`:

```bash
> cd /tmp
> mkdir goatproject
> cd goatproject
```

Now we need to actually make this directory a goat project. To do this all that
goat needs is for a `Goatfile` to exist in the root, and for it to define the
project's import-able path.

```bash
echo '{
    "path":"github.com/yourusername/goatproject"
}' > Goatfile
```

You'll see the meaning of `path` a bit later.

Now whenever `goat` is used (in place of `go`) on the command-line and you're
current working directory is in `goatproject` or one of `goatproject`'s children
your `GOPATH` will be set to (assuming you put your project in `/tmp` like me):

```
/tmp/goatproject/lib:/tmp/goatproject:$OLDGOPATH
```

This isn't too important, but it may be useful in the future so there ya go.

Our project has some dependencies. Normally we would fetch these with `go get`,
but we're too cool for that. Make your `Goatfile`'s contents be the following:

```json
{
    "path":"github.com/yourusername/goatproject",
    "deps":[
        {
            "loc":"code.google.com/p/go.example/newmath"
        }
    ]
}
```

This is the equivalent of having goat do a
`go get code.google.com/p/go.example/newmath` inside our project. To see a full
write-up on `Goatfile`'s dependency syntax and how to use it see the
[Goatfile](/docs/goatfile.md) documentation.

To actually download the dependency do (you'll need mercurial installed):

```bash
> cd /tmp/goatproject #if you hadn't already
> goat deps
```

This should fetch the dependency and put it in the `lib/` directory in your
project. You shouldn't check this directory into your version control (if
you're using any), it's just a utility for goat.

If in the future you change the Goatfile you can call `goat deps` again and it
will re-setup your `lib` directory with the changes.

## Using the dependency

We will now import the dependency we downloaded and use it in some of our own
code. Create the file `src/foo/foo.go` and in it put:

```go
package foo

import (
    "code.google.com/p/go.example/newmath"
)

func SqrtTwo() float64 {
    return newmath.Sqrt(2)
}

func SqrtThree() float64 {
    return newmath.Sqrt(3)
}
```

You can see the syntax for using this library is exactly the same as if we had
installed the dependency globally. This is because the dependency is in our
`lib` directory, which is at the front of our `GOPATH` (when using `goat`).

## Using our own package

We have a package in our project now, `foo`, that we'd like to use. Create the
file `src/main.go` and put in it:

```go
package main

import (
    "github.com/yourusername/goatproject/src/foo"
    "fmt"
)

func main() {
    fmt.Printf("The square root of two is: %v\n", foo.SqrtTwo())
}
```

The `github.com/yourusername/goatproject` part corresponds to the `path` field
in the `Goatfile`, and when you use it goat does a bit of magic so that when go
searches for that path it finds the project's root directory, even though the
`goatproject` directory isn't in a folder called `github.com/yourusername`.

In fact, the project's directory name could be anything at all. The only thing
that matters is that the `path` field and any import statements used in the code
match up. Additionally, the `src` directory we're putting everything in isn't
explicitely needed either, but I think it's better for organization.

## Action!

To actually run our code:

```bash
> cd /tmp/goatproject #if you haven't already
> goat run src/main.go
```

That's it! `goat` passes through to `go` all commands that it doesn't recognize
after setting up the environment variables, so it's super easy to get existing
and new projects up and running.
