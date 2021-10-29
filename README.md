# DO NOT USE

This library is still under heavy development, will probably break, and will definitely change.

See ~/journal/arg_parsing.md and ~/Git/bakeoff_argparse

# Naming Conventions

- `NewXXX`: Create an `XXX` from options
- `AddXXX`: Add a created `XXX` to your tree
- `WithXXX`: Create and add an `XXX` to your tree

# TODO: Next milestone: grabbit

- make colored help
- make help less verbose...
- add flag.TypeInfo to be set to the value.TypeInfo - then I can make default a list when appropriate
- print long values more neatly.... make a value.StringList() method?
- figure out what to do with --color , --help requires it?
- make help take an argument? - help = man, json, color, web, form, term, lsp, bash-completion, zsh-completion
- Write good doc comments, README, examples
- make an app.Test() method folks can add to their apps - should test for unique flag names between parent and child sections/commands for one thing
- go through TODOs

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
