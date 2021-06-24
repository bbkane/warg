package section

import (
	"log"

	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	v "github.com/bbkane/warg/value"
)

type CategoryMap = map[string]Category

type CategoryOpt = func(*Category)

type Category struct {
	Flags      f.FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands   c.CommandMap
	Categories CategoryMap
	HelpLong   string
	HelpShort  string
}

// New

func NewCategory(opts ...CategoryOpt) Category {
	category := Category{
		Flags:      make(map[string]f.Flag),
		Categories: make(map[string]Category),
		Commands:   make(map[string]c.Command),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
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

func AddCommand(name string, value c.Command) CategoryOpt {
	return func(app *Category) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Fatalf("command already exists: %#v\n", name)
		}
	}
}

func AddCategoryFlag(name string, value f.Flag) CategoryOpt {
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

func WithCategoryFlag(name string, empty v.Value, opts ...f.FlagOpt) CategoryOpt {
	return AddCategoryFlag(name, f.NewFlag(empty, opts...))
}

func WithCommand(name string, opts ...c.CommandOpt) CategoryOpt {
	return AddCommand(name, c.NewCommand(opts...))
}

func WithCategoryHelpLong(helpLong string) CategoryOpt {
	return func(cat *Category) {
		cat.HelpLong = helpLong
	}
}

func WithCategoryHelpShort(helpShort string) CategoryOpt {
	return func(cat *Category) {
		cat.HelpShort = helpShort
	}
}
