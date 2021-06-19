package clide

import (
	"log"
)

type Action = func(ValueMap) error

type CategoryMap = map[string]Category
type CommandMap = map[string]Command
type FlagMap = map[string]Flag
type ValueMap = map[string]Value

type CategoryOpt = func(*Category)
type CommandOpt = func(*Command)
type FlagOpt = func(*Flag)

type Category struct {
	Flags      FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands   CommandMap
	Categories CategoryMap
}
type Command struct {
	Action Action

	Flags FlagMap
}

type Flag struct {
	// Default will be shoved into Value if needed
	// can be nil
	// TODO: actually use this
	Default Value
	SetBy   string
	// Value holds what gets passed to the flag: --myflag value
	// and should be initialized to the empty value
	Value Value
}

// New

func NewCategory(opts ...CategoryOpt) Category {
	category := Category{
		Flags:      make(map[string]Flag),
		Categories: make(map[string]Category),
		Commands:   make(map[string]Command),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func NewCommand(opts ...CommandOpt) Command {
	category := Command{
		Flags: make(map[string]Flag),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func NewFlag(empty Value, opts ...FlagOpt) Flag {
	flag := Flag{}
	flag.Value = empty
	for _, opt := range opts {
		opt(&flag)
	}
	return flag
}

// CategoryOpt functions

func AddCategory(name string, value Category) CategoryOpt {
	return func(app *Category) {
		if _, alreadyThere := app.Categories[name]; !alreadyThere {
			app.Categories[name] = value
		} else {
			log.Fatalf("category already exists: %#v\n", name)
		}
	}
}

func AddCommand(name string, value Command) CategoryOpt {
	return func(app *Category) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Fatalf("command already exists: %#v\n", name)
		}
	}
}

func AddCategoryFlag(name string, value Flag) CategoryOpt {
	return func(app *Category) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}

	}
}

func WithCategory(name string, opts ...CategoryOpt) CategoryOpt {
	return AddCategory(name, NewCategory(opts...))
}

func WithCategoryFlag(name string, empty Value, opts ...FlagOpt) CategoryOpt {
	return AddCategoryFlag(name, NewFlag(empty, opts...))
}

func WithCommand(name string, opts ...CommandOpt) CategoryOpt {
	return AddCommand(name, NewCommand(opts...))
}

// CommandOpt

func AddCommandFlag(name string, value Flag) CommandOpt {
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

func WithCommandFlag(name string, empty Value, opts ...FlagOpt) CommandOpt {
	return AddCommandFlag(name, NewFlag(empty, opts...))
}

// FlagOpt

func WithDefault(value Value) FlagOpt {
	return func(flag *Flag) {
		flag.Default = value
	}
}
