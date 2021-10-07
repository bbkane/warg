# DO NOT USE

This library is still under heavy development, will probably break, and will definitely change.

See ~/journal/arg_parsing.md and ~/Git/bakeoff_argparse

# Naming Conventions

- `NewXXX`: Create an `XXX` from options
- `AddXXX`: Add a created `XXX` to your tree
- `WithXXX`: Create and add an `XXX` to your tree

# TODO: Next milestone: grabbit

- dear lord make a path value type that auto-expands home
- need WAY better error messages - figure out whether I want errors for users or errors for devs (with stacktraces)
- rm --version - folks can just use version subcommmand
- the commnd handlers need to be able to see app.Version - new context?
- write down a list of differences between container type values and scalar values - what I wish the methods were named etc
- turn --help Example tests into regular tests with an update flag to generate the help text golden file - need to make the app output more flexible for this - add color=true|false|auto support here?
- make OverrideVersionFlag customizable (with an action)
- Add ~/Code/Go/hello_testing/README.md to go notes on blog. Also expected, actual order convention
- write a good `errors` package - see bottom of README
- go through tests and change everything to `expected`, `actual`
- Fix grabbit subreddit-limit arg thing (it's set by appdefault to be a one element list, when the others are set by config to be a 2 element list) - this is probably going to be best handled by the user in docs - warg doesn't know these are related...
- Fix failing test derived from Grabbit! DONE!
- --help should never panic! Right now it does if it finds an improper config file
- Get errors a lot better... now that I'm actually trying to use it I'm running into them... Use ~/Code/Go/error_wrap_2
- go through TODOs
- add required flag
- add type of flag to help output
- add envvar option to flag
- firm up tests - does got or expected come first when comparing - also use testify better - see configreader/jsonreader
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
