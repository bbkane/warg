package section

import (
	"log"
	"strings"

	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/value"
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
	Flags flag.FlagMap
	// Commands holds the Commands under this Section
	Commands command.CommandMap
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
		Flags:    make(map[string]flag.Flag),
		Sections: make(map[string]Section),
		Commands: make(map[string]command.Command),
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
func AddCommand(name string, value command.Command) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// AddFlag adds an existing Flag to be made availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
func AddFlag(name string, value flag.Flag) SectionOpt {
	if !strings.HasPrefix(name, "-") {
		log.Panicf("helpFlags should start with '-': %#v\n", name)
	}
	return func(sec *Section) {
		if _, alreadyThere := sec.Flags[name]; !alreadyThere {
			sec.Flags[name] = value
		} else {
			log.Panicf("flag already exists: %#v\n", name)
		}

	}
}

func AddFlags(flagMap flag.FlagMap) SectionOpt {
	// TODO: can I abstract this somehow? Until then - copy paste!
	for name := range flagMap {
		if !strings.HasPrefix(name, "-") {
			log.Panicf("helpFlags should start with '-': %#v\n", name)
		}
	}
	return func(sec *Section) {
		for name, value := range flagMap {
			if _, alreadyThere := sec.Flags[name]; !alreadyThere {
				sec.Flags[name] = value
			} else {
				log.Panicf("flag already exists: %#v\n", name)
			}
		}
	}
}

// WithSection creates a Section and adds it underneath this Section. Panics if a Section with the same name already exists
func WithSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return AddSection(name, New(helpShort, opts...))
}

// WithFlag creates a Flag and makes it availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
func WithFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) SectionOpt {
	return AddFlag(name, flag.New(helpShort, empty, opts...))
}

// WithCommand creates a Command and adds it underneath this Section. Panics if a Command with the same name already exists
func WithCommand(name string, helpShort string, action command.Action, opts ...command.CommandOpt) SectionOpt {
	return AddCommand(name, command.New(helpShort, action, opts...))
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
