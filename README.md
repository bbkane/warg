# warg

Build hierarchical CLI applications with warg!

- warg uses [funcopt](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) style declarative APIs to keep CLIs readable (nested commands can be indented!) and terse. warg does not require code generation.
- warg is extremely interested in getting information into your app. Ensure a flag can be set from an environmental variable, configuration file, or default value by adding a single line to the flag declaration (configuration files also take some app-level config).
- warg is customizable. Add new types of flag values, config file formats, or --help outputs using the public API.
- warg is easy to add to, maintain, and remove from your project (if necessary). This follows mostly from warg being terse and declarative. If you decide to remove warg, simply remove the app declaration and turn the passed flags into other types of function arguments for your command handlers. Done!

## Butler Example

In addition to a few [examples in the docs](https://pkg.go.dev/go.bbkane.com/warg#pkg-examples), warg includes a small [example program](./examples/butler/main.go).

```go
app := warg.New(
	"butler",
	section.New(
		section.HelpShort("A virtual assistant"),
		section.Command(
			command.Name("present"),
			command.HelpShort("Formally present a guest (guests are never introduced, always presented)."),
			present,
			command.Flag(
				flag.Name("--name"),
				flag.HelpShort("Guest to address."),
				value.String,
				flag.Alias("-n"),
				flag.EnvVars("BUTLER_PRESENT_NAME", "USER"),
				flag.Required(),
			),
		),
	),
)
```

## Run Butler

Color can be toggled on/off/auto with the `--color` flag:

<p align="center">
  <img src="img/image-20220114210824919.png" alt="Sublime's custom image"/>
</p>

The default help for a command dynamically includes each flag's **current** value and how it was was set (passed flag, config, envvar, app default).

<p align="center">
  <img src="img/image-20220114212104654.png" alt="Sublime's custom image"/>
</p>

Of course, running it with the flag also works

<p align="center">
  <img src="img/image-20220114212309862.png" alt="Sublime's custom image"/>
</p>

## Apps Using Warg

- [fling](https://github.com/bbkane/fling/) - GNU Stow replacement to manage my dotfiles
- [grabbit](https://github.com/bbkane/grabbit) - Grab images from Reddit
- [starghaze](https://github.com/bbkane/starghaze/) - Save GitHub Starred repos to GSheets, Zinc

# Should You Use warg?

I'm using warg for my personal projects, but the API is not finalized and there
are some known issues (see below). I will eventually improve warg, but I'm currently ( 2021-11-19 )
taking a break from developing on warg to develop some CLIs with warg.

## Known Issues

- lists containing aggregate values ( values in list objects from configs ) should be checked to have the same size and source but that must currently be done by the application ( see [grabbit](https://github.com/bbkane/grabbit/blob/d1f30b87c4e5c8112f08e9889fa541dbeab66842/main.go#L311) )
- Many more types of values need to implemented. Especially StringEnumSlice, StringMap and Duration

## Alternatives

- [cobra](https://github.com/spf13/cobra) is by far the most popular CLI framework for Go. It relies on codegen.
- [cli](https://github.com/urfave/cli) is also very popular.
- I've used the now unmaintained [kingpin](https://github.com/alecthomas/kingpin) fairly successfully.

# Concepts

## Sections, Commands, and Flags

warg is designed to create hierarchical CLI applications similar to [azure-cli](https://github.com/Azure/azure-cli) (just to be clear, azure-cli is not built with warg, but it was my inspiration for warg). These apps use sections to group subcommands, and pass information via flags, not positional arguments. A few examples:

### azure-cli

```
az keyvault certificate show --name <name> --vault-name <vault-name>
```

If we try to dissect the parts of this command, we see that it:

- Starts with the app name (`az`).
- Narrows down intent with a section (`keyvault`). Sections are usually nouns and function similarly to a directory hierarchy on a computer - used to group related sections and commands so they're easy to find and use together.
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

# TODO

- make outline help for a command just show the flags
- should I auto add a color flag, what about a version subcommand. I literally want these in all my apps
- Add a sentinal value (UNSET?) to be used with optional flags that unsets the flag? sets the flag to the default value? So I can use fling without passing -i 'README.*' all the time :)
- use https://stackoverflow.com/a/16946478/2958070 for better number handling?
- zsh completion with https://www.dolthub.com/blog/2021-11-15-zsh-completions-with-subcommands/
- go through TODOs in code
- --help ideas: man, json, web, form, term, lsp, bash-completion, zsh-completion, compact

# TODO - value2

- rm flag Enum, Default
- update calling code
- update docs
- rm value package
