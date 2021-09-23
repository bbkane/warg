package command

import (
	"log"
	"strings"

	f "github.com/bbkane/warg/flag"
	v "github.com/bbkane/warg/value"
)

type Action = func(f.FlagValues) error
type CommandMap = map[string]Command
type CommandOpt = func(*Command)

type Command struct {
	Action Action
	Flags  f.FlagMap
	// Help is a required one-line description
	Help string
	// HelpLong is an optional longer description
	HelpLong string
}

func DoNothing(_ f.FlagValues) error {
	return nil
}

func NewCommand(helpShort string, action Action, opts ...CommandOpt) Command {
	category := Command{
		Help:   helpShort,
		Action: action,
		Flags:  make(map[string]f.Flag),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

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

func WithFlag(name string, helpShort string, empty v.EmptyConstructor, opts ...f.FlagOpt) CommandOpt {
	return AddFlag(name, f.NewFlag(helpShort, empty, opts...))
}

func HelpLong(helpLong string) CommandOpt {
	return func(cat *Command) {
		cat.HelpLong = helpLong
	}
}
