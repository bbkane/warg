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

func NewCommand(opts ...CommandOpt) Command {
	category := Command{
		Flags: make(map[string]f.Flag),
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

func WithAction(action Action) CommandOpt {
	return func(cmd *Command) {
		cmd.Action = action
	}
}

func WithFlag(name string, empty v.Value, opts ...f.FlagOpt) CommandOpt {
	return AddFlag(name, f.NewFlag(empty, opts...))
}

func HelpLong(helpLong string) CommandOpt {
	return func(cat *Command) {
		cat.HelpLong = helpLong
	}
}

func HelpShort(helpShort string) CommandOpt {
	return func(cat *Command) {
		cat.HelpShort = helpShort
	}
}
