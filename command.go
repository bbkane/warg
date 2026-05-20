package warg

import (
	"errors"
	"os"
	"sort"

	"go.bbkane.com/warg/metadata"
	"go.bbkane.com/warg/value"
)

// CmdOpt is a functional option for configuring a [Cmd] during creation.
type CmdOpt func(*Cmd)

// Unimplemented returns an [Action] that always returns an error.
// Useful as a placeholder while prototyping commands.
func Unimplemented() Action {
	return func(_ CmdContext) error {
		return errors.New("TODO: implement this command")
	}
}

// NewCmd creates a [Cmd] with the given short help text, action, and options.
// Use [NewSubCmd] to simultaneously create and attach a command to a [Section].
func NewCmd(helpShort string, action Action, opts ...CmdOpt) Cmd {
	command := Cmd{
		HelpShort:          helpShort,
		Action:             action,
		Flags:              make(FlagMap),
		AllowForwardedArgs: false,
		Footer:             "",
		HelpLong:           "",
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

// CmdFlag attaches an existing [Flag] to a command. Panics if a flag with the same name exists.
func CmdFlag(name string, value Flag) CmdOpt {
	return func(com *Cmd) {
		com.Flags.AddFlag(name, value)
	}
}

// CmdFlagMap attaches multiple existing flags to a command. Panics if any name already exists.
func CmdFlagMap(flagMap FlagMap) CmdOpt {
	return func(com *Cmd) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewCmdFlag creates a new [Flag] and attaches it to a command. Panics if a flag with the same name exists.
func NewCmdFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...FlagOpt) CmdOpt {
	return CmdFlag(name, NewFlag(helpShort, empty, opts...))
}

// CmdFooter sets an optional footer text displayed at the end of help output for this command.
func CmdFooter(footer string) CmdOpt {
	return func(cat *Cmd) {
		cat.Footer = footer
	}
}

// CmdHelpLong sets an optional extended description displayed in detailed help output.
func CmdHelpLong(helpLong string) CmdOpt {
	return func(cat *Cmd) {
		cat.HelpLong = helpLong
	}
}

// AllowForwardedArgs enables passing extra arguments after "--" to this command.
// Forwarded args are accessible via [CmdContext].ForwardedArgs.
//
// Example usage:
//
//	enventory exec --env prod -- go run .
func AllowForwardedArgs() CmdOpt {
	return func(cmd *Cmd) {
		cmd.AllowForwardedArgs = true
	}
}

// PassedFlags is a map of flag names to their resolved values, containing only flags
// that were set from any source (CLI, config, env var, or default).
// TODO: is this true?
type PassedFlags map[string]interface{} // This can just stay a string for the convenience of the user.

// CmdContext holds all parsed information passed to an [Action] when a command is executed.
// It includes resolved flags, forwarded args, I/O streams, and parse metadata.
type CmdContext struct {
	App           *App
	Flags         PassedFlags
	ForwardedArgs []string

	ParseState *ParseState

	// ParseMetadata to smuggle user-defined state (i.e., not flags) into an Action. I use this for mocks when testing
	ParseMetadata metadata.Metadata

	Stderr *os.File
	Stdin  *os.File
	Stdout *os.File
}

// Action is a function executed when a [Cmd] is matched during parsing.
// Return nil on success; return an error to signal failure.
type Action func(CmdContext) error

// CmdMap maps command names to [Cmd] instances within a [Section].
type CmdMap map[string]Cmd

// Empty reports whether the map contains no commands.
func (fm CmdMap) Empty() bool {
	return len(fm) == 0
}

// SortedNames returns the command names in alphabetical order.
func (fm CmdMap) SortedNames() []string {
	keys := make([]string, 0, len(fm))
	for k := range fm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	return keys
}

// Cmd represents an executable command within the CLI.
// Command names should be verbs (e.g., "add", "edit", "run").
// Do not construct directly; use [NewCmd] or [NewSubCmd].
type Cmd struct {
	// Action to run when command is invoked
	Action Action

	// Parsed Flags
	Flags FlagMap

	// AllowForwardedArgs indicates whether or not extra args are allowed after flags and following `--`.
	// These args will be accessible in CmdContext.ForwardedArgs.
	AllowForwardedArgs bool

	// Footer is yet another optional longer description.
	Footer string

	// HelpLong is an optional longer description
	HelpLong string

	// HelpShort is a required one-line description
	HelpShort string
}
