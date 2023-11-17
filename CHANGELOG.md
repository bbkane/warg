# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

Some versioning notes: warg's primary purpose is to serve my (@bbkane) CLIs, so
it's likely to have breaking changes whenever I discover a better way to do
something and remain < 1.0.0 forever.

## Unreleased

## Added

- `command.Context.Version` field and
- `warg.OverrideVersion()` to set version on app creation.

## Changed

- (WIP) - rm `AddVersionCommand()` in favor of `command.PrintVersion()`. Migrate by adding a "version" command manually

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
