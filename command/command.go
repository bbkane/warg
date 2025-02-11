package command

import (
	"context"
	"errors"
	"os"
	"sort"

	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value"
)

// PassedFlags holds a map of flag names to flag Values
type PassedFlags map[string]interface{} // This can just stay a string for the convenience of the user.

// Context holds everything a command needs.
type Context struct {
	AppName string

	// Context to smuggle user-defined state (i.e., not flags) into an Action. I use this for mocks when testing
	Context context.Context
	Flags   PassedFlags

	// Path passed either to a command or a section. Does not include executable name (os.Args[0])
	Path   []string
	Stderr *os.File
	Stdout *os.File

	// Version of this app
	Version string
}

// An Action is run as the result of a command
type Action func(Context) error

type HelpShort string

// Name of the command
type Name string

// A CommandMap holds Commands and is used by Sections
type CommandMap map[Name]Command

func (fm CommandMap) Empty() bool {
	return len(fm) == 0
}

func (fm CommandMap) SortedNames() []Name {
	keys := make([]Name, 0, len(fm))
	for k := range fm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	return keys
}

// A CommandOpt customizes a Command
type CommandOpt func(*Command)

// A Command will run code for you!
// The name of a Command should probably be a verb - add , edit, run, ...
// A Command should not be constructed directly. Use Command / New / ExistingCommand functions
type Command struct {
	// Action to run when command is invoked
	Action Action

	// Parsed Flags
	Flags flag.FlagMap

	// Footer is yet another optional longer description.
	Footer string

	// HelpLong is an optional longer description
	HelpLong string

	// HelpShort is a required one-line description
	HelpShort HelpShort
}

// DoNothing is a command action that simply returns an error.
// Useful for prototyping
func DoNothing(_ Context) error {
	return errors.New("NOTE: replace this command.DoNothing call")
}

// New builds a Command
func New(helpShort HelpShort, action Action, opts ...CommandOpt) Command {
	command := Command{
		HelpShort: helpShort,
		Action:    action,
		Flags:     make(flag.FlagMap),
		Footer:    "",
		HelpLong:  "",
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

// Flag adds an existing flag to a Command. It panics if a flag with the same name exists
func Flag(name flag.Name, value flag.Flag) CommandOpt {
	return func(com *Command) {
		com.Flags.AddFlag(name, value)
	}
}

// FlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func FlagMap(flagMap flag.FlagMap) CommandOpt {
	return func(com *Command) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func NewFlag(name flag.Name, helpShort flag.HelpShort, empty value.EmptyConstructor, opts ...flag.FlagOpt) CommandOpt {
	return Flag(name, flag.New(helpShort, empty, opts...))
}

// Footer adds an Help string to the command - useful from a help function
func Footer(footer string) CommandOpt {
	return func(cat *Command) {
		cat.Footer = footer
	}
}

// HelpLong adds an Help string to the command - useful from a help function
func HelpLong(helpLong string) CommandOpt {
	return func(cat *Command) {
		cat.HelpLong = helpLong
	}
}
