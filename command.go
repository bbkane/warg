package warg

import (
	"context"
	"errors"
	"os"
	"sort"

	"go.bbkane.com/warg/value"
)

// A CmdOpt customizes a Command
type CmdOpt func(*Cmd)

// Unimplemented() is an Action that simply returns an error.
// Useful for prototyping
func Unimplemented() Action {
	return func(_ CmdContext) error {
		return errors.New("TODO: implement this command")
	}
}

// NewCmd builds a Cmd
func NewCmd(helpShort string, action Action, opts ...CmdOpt) Cmd {
	command := Cmd{
		HelpShort: helpShort,
		Action:    action,
		Flags:     make(FlagMap),
		Footer:    "",
		HelpLong:  "",
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

// CmdFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func CmdFlag(name string, value Flag) CmdOpt {
	return func(com *Cmd) {
		com.Flags.AddFlag(name, value)
	}
}

// CmdFlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func CmdFlagMap(flagMap FlagMap) CmdOpt {
	return func(com *Cmd) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewCmdFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func NewCmdFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...FlagOpt) CmdOpt {
	return CmdFlag(name, NewFlag(helpShort, empty, opts...))
}

// CmdFooter adds an Help string to the command - useful from a help function
func CmdFooter(footer string) CmdOpt {
	return func(cat *Cmd) {
		cat.Footer = footer
	}
}

// CmdHelpLong adds an Help string to the command - useful from a help function
func CmdHelpLong(helpLong string) CmdOpt {
	return func(cat *Cmd) {
		cat.HelpLong = helpLong
	}
}

// PassedFlags holds a map of flag names to flag Values
type PassedFlags map[string]interface{} // This can just stay a string for the convenience of the user.

// CmdContext holds everything a command needs.
type CmdContext struct {
	App   *App
	Flags PassedFlags

	ParseState *ParseState

	// Context to smuggle user-defined state (i.e., not flags) into an Action. I use this for mocks when testing
	Context context.Context

	Stderr *os.File
	Stdout *os.File
}

// An Action is run as the result of a command
type Action func(CmdContext) error

// A CmdMap holds Commands and is used by Sections
type CmdMap map[string]Cmd

func (fm CmdMap) Empty() bool {
	return len(fm) == 0
}

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

// A Cmd will run code for you!
// The name of a Cmd should probably be a verb - add , edit, run, ...
// A Cmd should not be constructed directly. Use Cmd / New / ExistingCommand functions
type Cmd struct {
	// Action to run when command is invoked
	Action Action

	// Parsed Flags
	Flags FlagMap

	// Footer is yet another optional longer description.
	Footer string

	// HelpLong is an optional longer description
	HelpLong string

	// HelpShort is a required one-line description
	HelpShort string
}
