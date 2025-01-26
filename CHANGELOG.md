# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

Note that I update this changelog as I make changes, so the top version (right
below this description) is likely unreleased.

# v0.0.28

## Added

- `section.CommandMap` and `section.SectionMap` for pre-existing sections and commands

## Changed

- Renamed `command.ExistingFlags` to `commamd.FlagMap`. The plural name is awkward when adding a one-element `FlagMap` is perfectly ok (and nice when you don't want to keep typing the name...)
- `version` is now a parameter to `warg.New` instead of the (now removed)
`warg.OverrideVersion`. Changed because I ALWAYS want to set the version and I
think thats common for real-world CLIs. If passed an empty string, warg will
attempt to use the go module version.

# v0.0.27

## Removed

- Removed `WARG_PRE_V0_0_26_PARSE_ALGORITHM` environment variable and the associated parse algorithm. We're 500 lines of code simpler now!

# v0.0.26

## Changed

- Moved `SetBy` into the `Value` interface (`value.UpdatedBy()` - this allows
`Flag` to be readonly and and makes the coupling between setting the value and
updating `UpdatedBy` explicit
- Flags must now be the last things passed - i.e. `<appname> <section>
<command> <flag>...`. In addition, the only flag allowed after a `<section>` is
the `--help` flag (unfortunately the `--color` flag is NOT currently allowed to
be passed for section help). This simplifies the parsing and will help with tab
completion, and after that's implemented I might try to restore the old
behavior if I get too annoyed with this limitation. Temporarily, the old
behavior can be restored by setting the `WARG_PRE_V0_0_26_PARSE_ALGORITHM`
environment variable, but I plan to remove that in the next version.

Examples:

```
$ go run ./examples/butler --color false -h
Parse err: Parse args error: expecting section or command, got --color
exit status 64
```

```
$ WARG_PRE_V0_0_26_PARSE_ALGORITHM=1 go run ./examples/butler --color false -h
A virtual assistant
# ... more text ...
```

# v0.0.25

## Added

- `path.Path` type that users should call `Expand`/`MustExpand` on. See `Removed`

## Changed

- The signature of `value.EmptyConstructor` no longer has the opportunity to
return an error. I don't think anything ever returned an error, so this was
never needed.

## Removed

- Removed `value.FromInstance`. This is a **nasty silent breaking change**
because it'll break users at runtime (not compile time) when they update
`warg`. However, `value.FromInstance` was only used for expanding `~` in
`Path`'s from default values - no other value types needed it. It also wasn't
really testable since the value `~` expanded to is different per machine
(without env var shennanigans). I'm replacing it with a `Path` type that the
user is expected to call `Path`/`MustExpand` on.

Old:

```go
path := cmdCtx.Flags["--mypath"].(string)
```

New:

```go
path := cmdCtx.Flags["--mypath"].(path.Path).MustExpand()
```

# v0.0.24

This is the start of some big changes to `warg`. I'm doing this mostly to make
adding tab completion a lot easier, but also because I've spent a few years
with `warg` now and I think I can simplify a few things.

I'm making sure all the tests pass, but I'm making enough changes to enough
code I no longer remember the exact context for that I'm sure I'll have to make
some bugfix releases. I think this is acceptable because I'm the only user of
this library and I think it'll be faster overall with my limited time and
energy.

I'm marking each breaking change with a new version so I can update piecemeal
or all at once.

## Changed

- Removed section flags in favor of app global flags. This is strictly less
  flexible than section flags, but paves the way to building easier tab
  completion, and it's easy use existing flags in multiple commands if they
  don't need to be global.

# v0.0.23

## Changed

- make `warg.GoldenTest` use `GoldenTestArgs` and `ParseOpt`s

# v0.0.22

## Changed

- make `warg.GoldenTest` accept `ParseOpt`s instead of a hardcoded list of options

# v0.0.21

## Added

- `command.Context`: `Version`, `AppName`, `Path` fields. Justification: I want
  to pass these fields to OpenTelemetry in `starghaze`.
- `command.Context.Context` field: Justification: I want to smuggle mocks into
  my `command.Action`s when testing. Before, this, I added an ugly "mock
  selection" flag, and this is much cleaner.

## Changed

- move `warg.ParseResult.Path` to `command.Context.Path`.
- rm `warg.AddVersionCommand()` in favor of `warg.VersionCommand()`. Use with
  `section.ExistingCommand("version", warg.VersionCommand()),`. Justification:
  more declarative - I'd like to define all commands inside the root section
  instead of having another way to add a flag as a warg option.
- rm `warg.AddColorFlag()` in favor of `warg.ColorFlag()`. Use with
  `section.ExistingFlag("--color", warg.ColorFlag()),`. Same justification as
  `warg.VersionCommand()`.
- update `Parse()` to use `ParseOpt`s instead of positional args:
  `OverrideArgs`, `OverrideLookupFunc`. Justification: these have obvious
  defaults that only need overriding for tests, which also probably want to use
  other `ParseOpt`s.
- move `warg.OverrideStderr` and `warg.OverrideStdout` to be `ParseOpt`s
  instead of `AppOpt`s. Justification: This removes the need for these public
  fields in `App` and nicer for callers.

# v0.0.20

## Fixed

- Fix YAML config parsing for `value.Dict`

# v0.0.19

## Fixed

- Fix panic when using a `value.Dict` and calling `detailed.DetailedCommandHelp`

# v0.0.18

## Added

- `contained.Addr` and `contained.AddrPort`
- `flag.UnsetSentinel` to allow for unsetting flags
- `value.Dict` container
- `warg.GoldenTest`
