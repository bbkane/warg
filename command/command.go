package command

import (
	"errors"

	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value"
)

// A CommandOpt customizes a Command
type CommandOpt func(*Command)

// DoNothing is a command action that simply returns an error.
// Useful for prototyping
func DoNothing(_ Context) error {
	return errors.New("NOTE: replace this command.DoNothing call")
}

// NewCommand builds a Command
func NewCommand(helpShort string, action Action, opts ...CommandOpt) Command {
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
func Flag(name string, value flag.Flag) CommandOpt {
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
func NewFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) CommandOpt {
	return Flag(name, flag.NewFlag(helpShort, empty, opts...))
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
