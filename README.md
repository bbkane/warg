# warg

An opinionated CLI framework:

- declarative nested commands
- detailed, colored  `--help` output (including what a flag is currently set to)
- update flags from `os.Args`, config files, environment variables, and app defaults
- extend with new flag types, config file formats, or `--help` output
- snapshot testing support

# Project Status (2025-06-15)

`warg` is still a work in progress, but the "bones" are where I want them - I don't intend to change the basic way warg works, but I do plan to experiment with a few APIs (make the value and config APIs simpler), and add TUI generation. I think the largest breaking changes are behind me (see the [CHANGELOG](./CHANGELOG.md)). Read [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) to help understand the tooling.

I'm watching issues; please open one for any questions and especially BEFORE submitting a Pull request.

# Examples

All of the CLIs [on my profile](https://github.com/bbkane/bbkane) use warg.

See API docs (including code examples) at [pkg.go.dev](https://pkg.go.dev/go.bbkane.com/warg)

Simple "butler" example (full source [here](examples/butler/main.go)):

```go
app := warg.New(
  "butler",
  "v1.0.0",
  warg.NewSection(
    "A virtual assistant",
    warg.NewSubCmd(
      "present",
      "Formally present a guest (guests are never introduced, always presented).",
      present,
      warg.NewCmdFlag(
        "--name",
        "Guest to address.",
        scalar.String(),
        warg.Alias("-n"),
        warg.EnvVars("BUTLER_PRESENT_NAME", "USER"),
        warg.Required(),
      ),
    ),
  ),
)
```

<p align="center">
  <img src="img/image-20220114212104654.png" alt="Butler help screenshot"/>
</p>

# When to avoid warg

By design, warg apps have the following requirements:

- must contain at least one subcommand. This makes it easy to add further subcommands, such as a `version` subcommand.   It is not possible to design a warg app such that calling `<appname> --flag <value>` does useful work. Instead, `<appname> <command> --flag <value>` must be used.
- warg does not support positional arguments. Instead, use a required flag: `git clone <url>` would be `git clone --url <url>`. This makes parsing much easier, and I like the simplicity of it, even though it's more to type/tab-complete.

# Alternatives

- [cobra](https://github.com/spf13/cobra) is by far the most popular CLI framework for Go.
- [cli](https://github.com/urfave/cli) is also very popular.
- I haven't tried [ff](https://github.com/peterbourgon/ff), but it looks similar to warg, though less batteries-included
- I've used the now unmaintained [kingpin](https://github.com/alecthomas/kingpin) fairly successfully.

# Notes

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for notes on development tooling.
