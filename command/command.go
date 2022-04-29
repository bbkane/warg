package command

import (
	"errors"
	"sort"

	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value"
)

// PassedFlags holds a map of flag names to flag Values
type PassedFlags map[string]interface{} // This can just stay a string for the convenience of the user.

// Context holds everything a command needs. Right now that's just passed flags
type Context struct {
	Flags PassedFlags
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
	Action Action
	Flags  flag.FlagMap
	// HelpShort is a required one-line description
	HelpShort HelpShort
	// Footer is yet another optional longer description.
	Footer string
	// HelpLong is an optional longer description
	HelpLong string
}

// DoNothing is a command action that simply returns nil
// Useful for prototyping
func DoNothing(_ Context) error {
	return errors.New("TODO: replace this command.DoNothing call")
}

// New builds a Command
func New(helpShort HelpShort, action Action, opts ...CommandOpt) Command {
	command := Command{
		HelpShort: helpShort,
		Action:    action,
		Flags:     make(flag.FlagMap),
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

// ExistingFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func ExistingFlag(name flag.Name, value flag.Flag) CommandOpt {
	return func(com *Command) {
		com.Flags.AddFlag(name, value)
	}
}

func ExistingFlags(flagMap flag.FlagMap) CommandOpt {
	return func(com *Command) {
		com.Flags.AddFlags(flagMap)
	}
}

// Flag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func Flag(name flag.Name, helpShort flag.HelpShort, empty value.EmptyConstructor, opts ...flag.FlagOpt) CommandOpt {
	return ExistingFlag(name, flag.New(helpShort, empty, opts...))
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
