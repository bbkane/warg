package section

import (
	"log"

	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	v "github.com/bbkane/warg/value"
)

type SectionMap = map[string]Section

type SectionOpt = func(*Section)

type Section struct {
	Flags     f.FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands  c.CommandMap
	Sections  SectionMap
	HelpLong  string
	HelpShort string
}

func NewSection(opts ...SectionOpt) Section {
	category := Section{
		Flags:    make(map[string]f.Flag),
		Sections: make(map[string]Section),
		Commands: make(map[string]c.Command),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func AddSection(name string, value Section) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Fatalf("category already exists: %#v\n", name)
		}
	}
}

func AddCommand(name string, value c.Command) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Fatalf("command already exists: %#v\n", name)
		}
	}
}

func AddFlag(name string, value f.Flag) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}

	}
}

func WithSection(name string, opts ...SectionOpt) SectionOpt {
	return AddSection(name, NewSection(opts...))
}

func WithFlag(name string, empty v.Value, opts ...f.FlagOpt) SectionOpt {
	return AddFlag(name, f.NewFlag(empty, opts...))
}

func WithCommand(name string, opts ...c.CommandOpt) SectionOpt {
	return AddCommand(name, c.NewCommand(opts...))
}

func HelpLong(helpLong string) SectionOpt {
	return func(cat *Section) {
		cat.HelpLong = helpLong
	}
}

func HelpShort(helpShort string) SectionOpt {
	return func(cat *Section) {
		cat.HelpShort = helpShort
	}
}
