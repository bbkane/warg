# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

Some versioning notes: warg's primary purpose is to serve my (@bbkane) CLIs, so
it's likely to have breaking changes whenever I discover a better way to do
something and remain < 1.0.0 forever.

## Unreleased

## Added

- `command.Context`: `Version`, `AppName`, `Path` fields. Justification: I want
  to pass these fields to OpenTelemetry in `starghaze`

## Changed

- rm `warg.AddVersionCommand()` in favor of `warg.VersionCommand()`. Use with
  `section.ExistingCommand("version", warg.VersionCommand()),`. Justification:
  more declarative - I'd like to define all commands inside the root section
  instead of having another way to add a flag as a warg option.
- rm `warg.AddColorFlag()` in favor of `warg.ColorFlag()`. Use with
  `section.ExistingFlag("--color", warg.ColorFlag()),`. Same justification as
  `warg.VersionCommand()`

## v0.0.20

### Fixed

- Fix YAML config parsing for `value.Dict`

## v0.0.19

### Fixed

- Fix panic when using a `value.Dict` and calling `detailed.DetailedCommandHelp`

## v0.0.18

### Added

- `contained.Addr` and `contained.AddrPort`
- `flag.UnsetSentinel` to allow for unsetting flags
- `value.Dict` container
- `warg.GoldenTest`
