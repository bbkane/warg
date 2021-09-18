# DO NOT USE

This library is still under heavy development, will probably break, and will definitely change.

See ~/journal/arg_parsing.md and ~/Git/bakeoff_argparse

# Naming Conventions

- `NewXXX`: Create an `XXX` from options
- `AddXXX`: Add a created `XXX` to your tree
- `WithXXX`: Create and add an `XXX` to your tree

# TODO: Next milestone: grabbit

- upgrade config parser
- go through TODOs
- add required flag
- add type of flag to help output
- add envvar option to flag
- firm up tests - does got or expected come first when comparing
- should my config paths start with . to be jq compatible?
- make help take an argument? - help = man, json, color, web, form, term, lsp, bash-completion, zsh-completion

# Links

# Working towards

```go
root := NewSection(
    "help for example"
    s.HelpLong("example is an app that does x, y, z")
    s.WithFlag(
        "--sflag",
        "help for --sflag",
        v.NewEmptyStringValue(),
        f.Alias("-s"),
        f.ConfigPath("config.sflag", NewStringValueFromInterface), // interface[] -> (Value, error)
        f.Default(NewStringValue("hi")),
        f.EnvVar("example_sflag"),
        f.Examples([]string{"hola", "hello, governer!"})
        f.Required(),
    ),
    s.WithSection(
        "sec",
        "sec help",
        s.WithCommand(
            "com",
            "com help",
            func(vm ValueMap) error {
                cflag := vm["--cflag"].Get().(int)
                fmt.Println(cflag)
            }
            c.WithFlag(
                "--cflag",
                "help for --cflag",
                v.NewEmptyIntValue(),
            ),
            c.Examples(
                "do the thing: example sec com --cflag 2",
            )
        ),
    ),
)

app := a.NewApp(
    "example",
    "v0.0.0",
    a.Config(
        "--config",
        "path to config",
        map[string]Unmarshaller{
         ".json": a.JSONUnmarshaller,
         ".yaml": warg_yaml.YAMLUnmarshaller,
         ".yml": warg_yaml.YAMLUnmarshaller,
        },
       f.Default(v.NewStringValue("config.yml")),
    ),
    a.Help(
        []string{"-h", "--help"},
        a.DefaultSectionHelp,
        a.DefaultCommandHelp),
    a.Version([]string{"--version"}),
    a.AddRootSection(root),  // alternative to a.WithRootSection(s.SectionOpt...)
)

pr, _ := app.Parse(os.Args)
```
# Help

## Section

```
$ grabbit -h

helplong - Grab images from Reddit :)

Sections:
config : helpshort - change grabbit's config

Commands:
  grab : grab those images!!

Examples:

  # basic grab:
  grabbit grab \
      --subreddit bob \
      --limit 10 \
      --sort top \
      --timeframe weekly \
      --location ~/Downloads

  # config grab:
  grabbit grab  # all params in the config baby!
```

## Command

```
$ grabbit grab -h

helplong - do the grabby grabby

Flags:

  --subreddit : help - subreddit name
    short : -s
    configpath : subreddits[].name
    envvar : GRABBIT_SUBREDDIT_NAME
    examples :
      [art wallpapers]
    required : true
    setby : config
    type : stringslice
    value : [earthporn cityporn]

  --limit : how many pics to grab
    short : -s
    configpath : subreddits[].limit
    envvar : GRABBIT_SUBREDDIT_LIMIT
    examples :
      [art wallpapers]
    required : true
    setby : passedflag
    type : int
    value : 10

Examples:

  # basic grab:
  grabbit grab \
      --subreddit bob \
      --limit 10 \
      --sort top \
      --timeframe weekly \
      --location ~/Downloads

  # config grab:
  grabbit grab  # all params in the config baby!
```

# Cases for config:

These need to be added to get configs working well enough :)

config not passed -> flag not set
config passed, file doesn't exist -> flag not set  # command should error if a flag isn't set properly
config passed, file exists, can't unmarshall -> ERROR
config passed, file exists, can unmarshall, invalid path-> ERROR
config passed, file exists, can unmarshall, valid path, path not in config -> flag not set
config passed, file exists, can unmarshall, valid path, path in config, value error -> ERROR
config passed, file exists, can unmarshall, valid path, path in config, value created -> flag set
