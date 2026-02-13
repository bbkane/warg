# AGENTS.md

## Project Overview

**warg** is an opinionated CLI (Command Line Interface) framework for Go. It provides a declarative approach to building CLI applications with nested commands, detailed help output, and flexible flag configuration from multiple sources.

**Repository**: `go.bbkane.com/warg`
**License**: MIT
**Language**: Go

## Core Concepts

### Architecture

warg applications are built from four structural types:

1. **App** ([app.go](app.go)) - The top-level application container
2. **Section** ([section.go](section.go)) - "Folders" for organizing commands (noun names)
3. **Cmd** ([command.go](command.go)) - Executable commands (verb names like `add`, `edit`, `run`)
4. **Flag** ([flag.go](flag.go)) - Configuration options for commands

### Key Design Decisions

- **Subcommands required**: Apps must have at least one subcommand. `<appname> --flag <value>` alone won't work; use `<appname> <command> --flag <value>`
- **No positional arguments**: Use required flags instead (e.g., `git clone --url <url>` instead of `git clone <url>`)
- **Flag resolution order**: `os.Args` → config files → environment variables → app defaults

## Project Structure

```
warg/
├── app.go              # App type, New(), MustRun(), Parse()
├── app_parse.go        # Parsing logic, ParseState, ValueMap
├── command.go          # Cmd type, CmdContext, Action
├── flag.go             # Flag type and options
├── section.go          # Section type for organizing commands
├── help_*.go           # Help output formatters (detailed, outline, allcommands)
├── testing.go          # GoldenTest for snapshot testing
├── completion/         # Shell completion (zsh)
├── config/             # Config file readers (YAML, JSON)
├── value/              # Flag value types
│   ├── scalar/         # Single values (String, Int, Bool, Path, etc.)
│   ├── slice/          # List values
│   └── contained/      # Type definitions
├── examples/           # Example applications
│   ├── butler/         # Simple example
│   ├── fling/          # Link/unlink utility
│   ├── grabbit/        # Reddit image downloader
│   └── starghaze/      # GitHub stars manager
└── metadata/           # Key/value parse metadata
```

## Key APIs

### Creating an App

```go
app := warg.New(
    "appname",
    "v1.0.0",
    warg.NewSection(
        "Section help text",
        warg.NewSubCmd(
            "command",
            "Command help",
            actionFunc,
            warg.NewCmdFlag(
                "--flag",
                "Flag help",
                scalar.String(),
            ),
        ),
    ),
)
```

### Important Types

| Type | Purpose | File |
|------|---------|------|
| [`App`](app.go) | Application container | app.go |
| [`Section`](section.go) | Command grouping | section.go |
| [`Cmd`](command.go) | Executable command | command.go |
| [`CmdContext`](command.go) | Context passed to actions | command.go |
| [`Action`](command.go) | Command handler function `func(CmdContext) error` | command.go |
| [`Flag`](flag.go) | Command flag definition | flag.go |
| [`ParseResult`](app_parse.go) | Result of parsing | app_parse.go |
| [`ParseState`](app_parse.go) | Current parsing state | app_parse.go |

### App Options (AppOpt)

- `warg.GlobalFlag()` / `warg.NewGlobalFlag()` - Add global flags
- `warg.ConfigFlag()` - Enable config file support
- `warg.HelpFlag()` - Customize help
- `warg.SkipAll()` - Skip all auto-added features (for tests)
- `warg.SkipCompletionCmds()` - Skip completion commands
- `warg.SkipGlobalColorFlag()` - Skip `--color` flag
- `warg.SkipVersionCmd()` - Skip `version` command
- `warg.SkipValidation()` - Skip startup validation

### Flag Options (FlagOpt)

- `warg.Alias()` - Short alias (e.g., `-n`)
- `warg.EnvVars()` - Environment variable names
- `warg.Required()` - Mark as required
- `warg.ConfigPath()` - Path in config file
- `warg.FlagCompletions()` - Tab completion function

## Testing

### Running Tests

```bash
go test ./...
```

### Updating Golden Files

```bash
WARG_TEST_UPDATE_GOLDEN=1 go test ./...
```

### Golden Test Helper

Use [`warg.GoldenTest`](testing.go) for snapshot testing:

```go
warg.GoldenTest(
    t,
    warg.GoldenTestArgs{
        App:             app(),
        UpdateGolden:    updateGolden,
        ExpectActionErr: false,
    },
    warg.ParseWithArgs([]string{"app", "cmd", "--flag", "value"}),
    warg.ParseWithLookupEnv(warg.LookupMap(nil)),
)
```

## Code Style Guidelines

1. **Package imports**: Most warg types are in the main `warg` package. Value types are in `value/scalar` and `value/slice`
2. **Error handling**: Commands return errors; `MustRun()` exits on error
3. **Naming conventions**:
   - Commands use verb names (`add`, `edit`, `delete`)
   - Sections use noun names (`config`, `comments`)
   - Flags start with `--` and optionally have `-` aliases

## Development Tooling

- **Linting**: `.golangci.yml` configures golangci-lint
- **Git hooks**: `lefthook.yml` for pre-commit hooks
- **YAML linting**: `.yamllint.yml`

See [Go Project Notes](https://www.bbkane.com/blog/go-project-notes/) for additional tooling documentation.

## Common Patterns

### Config File Support

```go
warg.ConfigFlag(
    yamlreader.New,
    warg.FlagMap{
        "--config": warg.NewFlag(
            "Path to config file",
            scalar.Path(),
        ),
    },
)
```

### Forwarded Arguments

For commands that wrap other tools:

```go
warg.NewSubCmd(
    "exec",
    "Execute with environment",
    execAction,
    warg.AllowForwardedArgs(),
)
// Usage: app exec --env prod -- go run .
```

### Custom Completions

```go
warg.FlagCompletions(warg.CompletionsValues([]string{"option1", "option2"}))
warg.FlagCompletions(warg.CompletionsDirectories())
warg.FlagCompletions(warg.CompletionsDirectoriesFiles())
```

