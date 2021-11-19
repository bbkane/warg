# warg

Build heirarchical CLI applications with warg!

- warg uses [funcopt](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) style APIs to keep CLI declaration readable and terse. It does not require code generation. Nested CLI command are indented, which makes apps easy to debug.
- warg is extremely interested in getting information into your app. Ensure a flag can be set from an environmental variable, configuration file, or default value by adding a single line to the flag (configuration files also take some app-level config).
- warg is customizable. Add new types of flag values, config file formats, or --help outputs using the public API.
- warg is easy to add to, maintain, and remove from your project (if necessary). This follows mostly from warg being terse and declarative. If you decide to remove warg, simply remove the app declaration and turn the passed flags into other types of function arguments for your command handlers. Done!


# Hello World

Also see the examples in [the docs](https://pkg.go.dev/github.com/bbkane/warg).

## Code

```go
package main

import (
	"fmt"
	"os"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
)

func hello(pf flag.PassedFlags) error {
	// this is a required flag, so we know it exists
	name := pf["--name"].(string)
	fmt.Printf("Hello %s!\n", name)
	return nil
}

func main() {
	app := warg.New(
		"say",
		section.New(
			"Make the terminal say things!!",
			section.WithCommand(
				"hello",
				"Say hello",
				hello,
				command.WithFlag(
					"--name",
					"Person we're talking to",
					value.String,
					flag.Alias("-n"),
					flag.EnvVars("SAY_NAME"),
					flag.Required(),
				),
			),
		),
	)
	app.MustRun(os.Args, os.LookupEnv)
}
```

## Run

By default, these help messages are in color. You'll have to imagine that within this README :)

```
$ ./say -h
Make the terminal say things!!

Commands

  hello : Say hello
```

The default help for a command dynamically includes each flag's **current** value and how it was was set (passed flag, config, envvar, app default).

```
$ ./say hello --name World -h
Say hello

Command Flags:

  --name , -n : Person we're talking to
    type : string
    envvars : [SAY_NAME]
    required : true
    value (set by passedflag) : World

Inherited Section Flags:

  --help , -h : Print help
    type : stringenum with choices: [default]
    default : default
    value (set by appdefault) : default
```
```
$ SAY_NAME=Bob ./say hello -h
Say hello

Command Flags:

  --name , -n : Person we're talking to
    type : string
    envvars : [SAY_NAME]
    required : true
    value (set by envvar) : Bob

Inherited Section Flags:

  --help , -h : Print help
    type : stringenum with choices: [default]
    default : default
    value (set by appdefault) : default
```

```
$ ./say hello --name World
Hello World!
```

# Should You Use warg?

I'm using warg for my personal projects, but the API is not finalized and there
are some known issues (see below). I will eventually improve warg, but I'm currently ( 2021-11-19 )
taking a break from developing on warg to develop some CLIs with warg.

## Known Issues

- warg does not warn you if a child section/command has a flag with the same name as a parent. The child flag essentially overwrites the parent flag. I'd like to check this at test time.
- lists containing aggregate values ( values in list objects from configs ) should be checked to have the same size and source but that must currently be done by the application ( see [grabbit](https://github.com/bbkane/grabbit/blob/d1f30b87c4e5c8112f08e9889fa541dbeab66842/main.go#L311) )
- Many more types of values need to implemented. Especially StringEnumSlice, StringMap and Duration

## Alternatives

- [cobra](https://github.com/spf13/cobra) is by far the most popular CLI framework for Go. It relies on codegen.
- [cli](https://github.com/urfave/cli) is also very popular.
- I've used the now unmaintained [kingpin](https://github.com/alecthomas/kingpin) fairly successfully.

# Concepts

## Sections, Commands, and Flags

warg is designed to create heirarchical CLI applications similar to [azure-cli](https://github.com/Azure/azure-cli) (just to be clear, azure-cli is not built with warg, but it was my inspiration for warg). These apps use sections to group subcommands, and pass information via flags, not positional arguments. A few examples:

### azure-cli

```
az keyvault certificate show --name <name> --vault-name <vault-name>
```

If we try to dissect the parts of this command, we see that it:

- Starts with the app name (`az`).
- Narrows down intent with a section (`keyvault`). Sections are usually nouns and function similarly to a directory heirarchy on a computer - used to group related sections and commands so they're easy to find and use together.
- Narrows down intent further with another section (`certificate`).
- Ends with a command (`show`). Commands are usually verbs and specify a single action to take within that section.
- Passes information to the command with flags (`--name`, `--vault-name`).

This structure is both readable and scalable. `az` makes hundreds of commands browsable with this strategy!

### grabbit

[grabbit](https://github.com/bbkane/grabbit) is a much smaller app to download wallpapers from Reddit that IS built with warg. It still benefits from the sections/commands/flags structure. Let's organize some of grabbit's components into a tree diagram:

```
grabbit                   # app name
├── --color               # section flag
├── --config-path         # section flag
├── --help                # section flag
├── config                # section
│   └── edit              # command
│       └── --editor      # command flag
├── grab                  # command
│   └── --subreddit-name  # command flag
└── version               # command
```

Similar to `az`, `grabbit` organizes its capabilities with sections, commands and flags. Sections are used to group commands. Flags defined in a "parent" section are available to child commands. for example, the `config edit` command has access to the parent `--config-path` flag, as does the `grab` command.

## Special Flags

TODO

--config

--help + --color

## Unsupported CLI Patterns

One of warg's tradeoffs is that it insists on only using sections, commands and flags. This means it is not possible (by design) to build some styles of CLI apps. warg does not support positional arguments. Instead, use a required flag: `git clone <url>` is spelled `git clone --url <url>`.

All warg apps must have at least one nested command.  It is not possible to design a warg app such that calling `<appname> --flag <value>` does useful work. Instead, `<appname> <command> --flag <value>` must be used.

##  API Naming Conventions

Replace `XXX` with `Section`, `Command`, or `Flag`:

- `xxx.New`: Create an `XXX` from options
- `AddXXX`: Add a created `XXX` to your tree
- `WithXXX`: Create and add an `XXX` to your tree

# TODO

- turn configreader into config and jsonreader.NewJSONConfigReader -> New
- make config/path package so the code isn't copied
- replace WithXXX with XXX ? that means renaming the XXX type to something like XXXT (SectionT)
- add screenshots for --help - colors look way better
- zsh completion with https://www.dolthub.com/blog/2021-11-15-zsh-completions-with-subcommands/
- Should I make commands not return an error? Maybe that should be handled by the app author?
- Ensure a flag created with flag.New can be used in multiple places! Probably with lots of tests...
- make help less verbose...
- don't skip the --help tests :)
- make an app.Test() method folks can add to their apps - should test for unique flag names between parent and child sections/commands for one thing
- go through TODOs in code
- --help ideas: man, json, web, form, term, lsp, bash-completion, zsh-completion, outline, compact
