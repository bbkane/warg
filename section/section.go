package section

import (
	"log"
	"strings"

	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/value"
)

// SectionMap holds Sections - used by other Sections
type SectionMap = map[string]SectionT

// SectionOpt customizes a Section on creation
type SectionOpt = func(*SectionT)

// Sections are like "folders" for Commmands.
// They should have noun names.
// Sections should not be created in place - use New/With/AddSection functions.
// SectionT is the type name because we need the more user-visible `Section` as a function name
type SectionT struct {
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
func New(helpShort string, opts ...SectionOpt) SectionT {
	section := SectionT{
		Help:     helpShort,
		Flags:    make(map[string]flag.Flag),
		Sections: make(map[string]SectionT),
		Commands: make(map[string]command.Command),
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

// ExistingSection adds an existing Section underneath this Section. Panics if a Section with the same name already exists
func ExistingSection(name string, value SectionT) SectionOpt {
	return func(app *SectionT) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// ExistingCommand adds an existing Command underneath this Section. Panics if a Command with the same name already exists
func ExistingCommand(name string, value command.Command) SectionOpt {
	return func(app *SectionT) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// ExistingFlag adds an existing Flag to be made availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
func ExistingFlag(name string, value flag.Flag) SectionOpt {
	if !strings.HasPrefix(name, "-") {
		log.Panicf("helpFlags should start with '-': %#v\n", name)
	}
	return func(sec *SectionT) {
		if _, alreadyThere := sec.Flags[name]; !alreadyThere {
			sec.Flags[name] = value
		} else {
			log.Panicf("flag already exists: %#v\n", name)
		}

	}
}

func ExistingFlags(flagMap flag.FlagMap) SectionOpt {
	// TODO: can I abstract this somehow? Until then - copy paste!
	for name := range flagMap {
		if !strings.HasPrefix(name, "-") {
			log.Panicf("helpFlags should start with '-': %#v\n", name)
		}
	}
	return func(sec *SectionT) {
		for name, value := range flagMap {
			if _, alreadyThere := sec.Flags[name]; !alreadyThere {
				sec.Flags[name] = value
			} else {
				log.Panicf("flag already exists: %#v\n", name)
			}
		}
	}
}

// Section creates a Section and adds it underneath this Section. Panics if a Section with the same name already exists
func Section(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return ExistingSection(name, New(helpShort, opts...))
}

// Flag creates a Flag and makes it availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
func Flag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) SectionOpt {
	return ExistingFlag(name, flag.New(helpShort, empty, opts...))
}

// Command creates a Command and adds it underneath this Section. Panics if a Command with the same name already exists
func Command(name string, helpShort string, action command.Action, opts ...command.CommandOpt) SectionOpt {
	return ExistingCommand(name, command.New(helpShort, action, opts...))
}

// Footer adds an optional help string to this Section
func Footer(footer string) SectionOpt {
	return func(cat *SectionT) {
		cat.Footer = footer
	}
}

// HelpLong adds an optional help string to this Section
func HelpLong(helpLong string) SectionOpt {
	return func(cat *SectionT) {
		cat.HelpLong = helpLong
	}
}
