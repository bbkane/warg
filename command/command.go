package command

import (
	"log"
	"strings"

	f "github.com/bbkane/warg/flag"
	v "github.com/bbkane/warg/value"
)

// An Action is run as the result of a command
type Action = func(f.PassedFlags) error

// A CommandMap holds Commands and is used by Sections
type CommandMap = map[string]Command

// A CommandOpt customizes a Command
type CommandOpt = func(*Command)

// A Command will run code for you!
// The name of a Command should probably be a verb - add , edit, run, ...
// It should not be constructed directly - use AddCommand/NewCommand/WithCommand functions
type Command struct {
	Action Action
	Flags  f.FlagMap
	// Help is a required one-line description
	Help string
	// Footer is yet another optional longer description.
	Footer string
	// HelpLong is an optional longer description
	HelpLong string
}

// DoNothing is a command action that simply returns nil
// Useful for prototyping
func DoNothing(_ f.PassedFlags) error {
	return nil
}

// New builds a Command
func New(helpShort string, action Action, opts ...CommandOpt) Command {
	command := Command{
		Help:   helpShort,
		Action: action,
		Flags:  make(map[string]f.Flag),
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

// AddFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func AddFlag(name string, value f.Flag) CommandOpt {
	if !strings.HasPrefix(name, "-") {
		log.Panicf("helpFlags should start with '-': %#v\n", name)
	}
	return func(app *Command) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Panicf("flag already exists: %#v\n", name)
		}
	}
}

// WithFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func WithFlag(name string, helpShort string, empty v.EmptyConstructor, opts ...f.FlagOpt) CommandOpt {
	return AddFlag(name, f.New(helpShort, empty, opts...))
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
