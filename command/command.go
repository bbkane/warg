package command

import (
	"log"

	f "github.com/bbkane/warg/flag"
	v "github.com/bbkane/warg/value"
)

type Action = func(v.ValueMap) error
type CommandMap = map[string]Command
type CommandOpt = func(*Command)

type Command struct {
	Action    Action
	Flags     f.FlagMap
	HelpLong  string
	HelpShort string
}

func DoNothing(_ v.ValueMap) error {
	return nil
}

func NewCommand(helpShort string, action Action, opts ...CommandOpt) Command {
	category := Command{
		HelpShort: helpShort,
		Action:    action,
		Flags:     make(map[string]f.Flag),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func AddFlag(name string, value f.Flag) CommandOpt {
	return func(app *Command) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}
	}
}

func WithFlag(name string, helpShort string, empty v.Value, opts ...f.FlagOpt) CommandOpt {
	return AddFlag(name, f.NewFlag(helpShort, empty, opts...))
}

func HelpLong(helpLong string) CommandOpt {
	return func(cat *Command) {
		cat.HelpLong = helpLong
	}
}
