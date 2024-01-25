# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

Note that I update this changelog as I make changes, so the top version (right
below this description) is likely unreleased.

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

### Fixed

- Fix YAML config parsing for `value.Dict`

# v0.0.19

### Fixed

- Fix panic when using a `value.Dict` and calling `detailed.DetailedCommandHelp`

# v0.0.18

### Added

- `contained.Addr` and `contained.AddrPort`
- `flag.UnsetSentinel` to allow for unsetting flags
- `value.Dict` container
- `warg.GoldenTest`
