# DO NOT USE

This library is still under heavy development, will probably break, and will definitely change.

See ~/journal/arg_parsing.md and ~/Git/bakeoff_argparse

# Naming Conventions

- `NewXXX`: Create an `XXX` from options
- `AddXXX`: Add a created `XXX` to your tree
- `WithXXX`: Create and add an `XXX` to your tree

# TODO: Next milestone: grabbit

- finish tests for configreader/jsonreader - note that I'm using testify in new and exciting ways so that might be broken too...
- get grabbit working with YAML - turns out YAML and JSON need fundamentally incompatible ConfigMaps - see "The Case for ConfigReader" at the bottom and ~/warg_configreader.md
- Add ~/Code/Go/hello_testing/README.md to go notes on blog. Also expected, actual order convention
- write a good `errors` package - see bottom of README
- go through tests and change everything to `expected`, `actual`
- Fix grabbit subreddit-limit arg thing (it's set by appdefault to be a one element list, when the others are set by config to be a 2 element list) - this is probably going to be best handled by the user in docs - warg doesn't know these are related...
- Fix failing test derived from Grabbit! DONE!
- --help should never panic! Right now it does if it finds an improper config file
- Get errors a lot better... now that I'm actually trying to use it I'm running into them... Use ~/Code/Go/error_wrap_2
- upgrade config parser
- go through TODOs
- add required flag
- add type of flag to help output
- add envvar option to flag
- firm up tests - does got or expected come first when comparing - also use testify better - see configreader/jsonreader
- should my config paths start with . to be jq compatible? Nah... to be fully jq compatible, they'd also ahave to be surrounded by `[]` - DONE
- make help take an argument? - help = man, json, color, web, form, term, lsp, bash-completion, zsh-completion

# Links

# Working towards

```go
s.WithFlag(
    "--greeting",
    "help for --greeting",
    v.StringEmpty,
    f.Short("-g"),
    f.ConfigPath("people.greeting", NewStringValueFromInterface),
    f.Default("hi"),
    f.Choices("hi", "hello", "hola"), // TODO: how does this work with container type values? Probably just constrain what's passed to their update functions (i.e., not able to constrain length of one for example) - folks could also make custom values if they need something more specialized
    f.EnvVar("GREETING"),
    // TODO: chose a style for examples
    f.Example("hello", "use 'hello' for prim and proper greetings"),
    f.Example("hola", "use 'hola' for fun greetings"),
    f.Examples(
        {"hello", "use 'hello' for prim and proper greetings"},
        {"hola", "use 'hola' for fun greetings"},
    )
    f.Required(),
),
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

# The case for ConfigReader

When it comes to map interfaces, JSON can only decode into `map[string]interface{}` and YAML can only decode internal maps into `map[interface{}]interface{}` . This is bad because my code relies on `type ConfigMap = map[string]interface{}` everywhere.

Right now, I've got the JSON one working, but I need to change how that part works...

```go
type ConfigSearchResult struct {
	IFace      interface{}
	Exists     bool
	IsAggregated bool
}

type ConfigReader interface {
	Search(path string) (interface{}, error)
}

type NewConfigReader = func(filePath string) (ConfigReader, error)
```

with the following hierarchy:

```
configreader # contains above definitions
configreader/json  # implements them!
```

Other packages, such as `configreader_yaml` (which I can put in here for now, but should move into its own package when I want to shed dependencies), can use the interfaces too

# Ideas for errors

There are two types of errors I care about:

- errors I'm going to report (for dev or users)
  - I want error location, messages, and key/value pairs on these (similar to logging)
- errors I'm going to programmatically handle
  - I want a good type to check against on these - can be implemented with a custom struct or errors.New()

The new `errors` improvements (errors.Is/As/%w) + some custom code can give me this I think!

Maybe an output like:

```
> /path/to/file.go:133 myFunc : my message
  key: value
  key: value
> /path/to/other/file.go:145 DoTheThing : source
```

structerr - structured errors? saniterry? saniterror?

```
NewErr(msg string, keysAndValues ...interface{})

NewErrWithContext(err, msg string, keysAndValues ...interface{})

NewErrWithStack(err)

Format(filePath string, lineNumber int, funcName string, message string, keysAndValues ...interface{})
```
