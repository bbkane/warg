package section

import (
	"log"
	"strings"

	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	v "github.com/bbkane/warg/value"
)

type SectionMap = map[string]Section

type SectionOpt = func(*Section)

type Section struct {
	Flags    f.FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands c.CommandMap
	Sections SectionMap
	// Help is a required one-line descripiton of this section
	Help string
	// HelpLong is an optional longer description of this section
	HelpLong string
	// Footer is yet another optional longer description.
	Footer string
}

func NewSection(helpShort string, opts ...SectionOpt) Section {
	section := Section{
		Help:     helpShort,
		Flags:    make(map[string]f.Flag),
		Sections: make(map[string]Section),
		Commands: make(map[string]c.Command),
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

func AddSection(name string, value Section) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

func AddCommand(name string, value c.Command) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

func AddFlag(name string, value f.Flag) SectionOpt {
	if !strings.HasPrefix(name, "-") {
		log.Panicf("helpFlags should start with '-': %#v\n", name)
	}
	return func(app *Section) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Panicf("flag already exists: %#v\n", name)
		}

	}
}

func WithSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return AddSection(name, NewSection(helpShort, opts...))
}

func WithFlag(name string, helpShort string, empty v.EmptyConstructor, opts ...f.FlagOpt) SectionOpt {
	return AddFlag(name, f.NewFlag(helpShort, empty, opts...))
}

func WithCommand(name string, helpShort string, action c.Action, opts ...c.CommandOpt) SectionOpt {
	return AddCommand(name, c.NewCommand(helpShort, action, opts...))
}

func Footer(footer string) SectionOpt {
	return func(cat *Section) {
		cat.Footer = footer
	}
}

func HelpLong(helpLong string) SectionOpt {
	return func(cat *Section) {
		cat.HelpLong = helpLong
	}
}
