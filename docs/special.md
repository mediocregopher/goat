# Goat special features

Here's some special features that goat has that don't really fit in with the
rest of the documentation.

## GOAT_ACTUALGO

When passing commands through to the actual go binary goat will, by default, use
the `env` command to find the binary and call the return from that. However, if
you've renamed the goat binary to `go` and put it on your path (useful for
easy-integrations with plugins like
[syntastic](https://github.com/scrooloose/syntastic)) then that's going to cause
an infinite loop. You can set an environment variable called `GOAT_ACTUALGO`
with the absolute path to your installed go binary. Goat will use this instead
of `env` when it's present.
