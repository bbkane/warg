package command

import (
	"errors"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value"
)

// DoNothing is a command action that simply returns an error.
// Useful for prototyping
func DoNothing(_ warg.CommandContext) error {
	return errors.New("NOTE: replace this command.DoNothing call")
}

// Flag adds an existing flag to a Command. It panics if a flag with the same name exists
func Flag(name string, value warg.Flag) warg.CommandOpt {
	return func(com *warg.Command) {
		com.Flags.AddFlag(name, value)
	}
}

// FlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func FlagMap(flagMap warg.FlagMap) warg.CommandOpt {
	return func(com *warg.Command) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func NewFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...warg.FlagOpt) warg.CommandOpt {
	return Flag(name, warg.NewFlag(helpShort, empty, opts...))
}

// Footer adds an Help string to the command - useful from a help function
func Footer(footer string) warg.CommandOpt {
	return func(cat *warg.Command) {
		cat.Footer = footer
	}
}

// HelpLong adds an Help string to the command - useful from a help function
func HelpLong(helpLong string) warg.CommandOpt {
	return func(cat *warg.Command) {
		cat.HelpLong = helpLong
	}
}
