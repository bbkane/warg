package section

import (
	"log"
	"strings"

	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	v "github.com/bbkane/warg/value"
)

// SectionMap holds Sections - used by other Sections
type SectionMap = map[string]Section

// SectionOpt customizes a Section on creation
type SectionOpt = func(*Section)

// Sections are like "folders" for Commmands.
// They should have noun names.
// Sections should not be created in place - use New/With/AddSection functions
type Section struct {
	// Flags holds flags available to this Section and all subsections and Commands
	Flags f.FlagMap
	// Commands holds the Commands under this Section
	Commands c.CommandMap
	// Sections holds the Sections under this Section
	Sections SectionMap
	// Help is a required one-line descripiton of this section
	Help string
	// HelpLong is an optional longer description of this section
	HelpLong string
	// Footer is yet another optional longer description.
	Footer string
}

// New creates a Section!
func New(helpShort string, opts ...SectionOpt) Section {
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

// AddSection adds an existing Section underneath this Section. Panics if a Section with the same name already exists
func AddSection(name string, value Section) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// AddCommand adds an existing Command underneath this Section. Panics if a Command with the same name already exists
func AddCommand(name string, value c.Command) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// AddFlag adds an existing Flag to be made availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
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

// WithSection creates a Section and adds it underneath this Section. Panics if a Section with the same name already exists
func WithSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return AddSection(name, New(helpShort, opts...))
}

// WithFlag creates a Flag and makes it availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
func WithFlag(name string, helpShort string, empty v.EmptyConstructor, opts ...f.FlagOpt) SectionOpt {
	return AddFlag(name, f.New(helpShort, empty, opts...))
}

// WithCommand creates a Command and adds it underneath this Section. Panics if a Command with the same name already exists
func WithCommand(name string, helpShort string, action c.Action, opts ...c.CommandOpt) SectionOpt {
	return AddCommand(name, c.New(helpShort, action, opts...))
}

// Footer adds an optional help string to this Section
func Footer(footer string) SectionOpt {
	return func(cat *Section) {
		cat.Footer = footer
	}
}

// HelpLong adds an optional help string to this Section
func HelpLong(helpLong string) SectionOpt {
	return func(cat *Section) {
		cat.HelpLong = helpLong
	}
}
