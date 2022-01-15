package command

import (
	"log"
	"strings"

	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/value"
)

// An Action is run as the result of a command
type Action = func(flag.PassedFlags) error

type HelpShort string

// Name of the command
type Name string

// A CommandMap holds Commands and is used by Sections
type CommandMap = map[Name]Command

// A CommandOpt customizes a Command
type CommandOpt = func(*Command)

// A Command will run code for you!
// The name of a Command should probably be a verb - add , edit, run, ...
// It should not be constructed directly - use AddCommand/NewCommand/WithCommand functions
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
func DoNothing(_ flag.PassedFlags) error {
	return nil
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
	if !strings.HasPrefix(string(name), "-") {
		log.Panicf("flags should start with '-': %#v\n", name)
	}
	return func(app *Command) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Panicf("flag already exists: %#v\n", name)
		}
	}
}

func ExistingFlags(flagMap flag.FlagMap) CommandOpt {
	// TODO: can I abstract this somehow? Until then - copy paste!
	for name := range flagMap {
		if !strings.HasPrefix(string(name), "-") {
			log.Panicf("helpFlags should start with '-': %#v\n", name)
		}
	}
	return func(sec *Command) {
		for name, value := range flagMap {
			if _, alreadyThere := sec.Flags[name]; !alreadyThere {
				sec.Flags[name] = value
			} else {
				log.Panicf("flag already exists: %#v\n", name)
			}
		}
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
