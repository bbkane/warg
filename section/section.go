package section

import (
	"log"
	"sort"
	"strings"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value"
)

// Name of the section
type Name string

// HelpShort is a required short description of the section
type HelpShort string

// SectionMap holds Sections - used by other Sections
type SectionMap map[Name]SectionT

func (fm *SectionMap) SortedNames() []Name {
	keys := make([]Name, 0, len(*fm))
	for k := range *fm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	return keys
}

// SectionOpt customizes a Section on creation
type SectionOpt func(*SectionT)

// Sections are like "folders" for Commmands.
// They should have noun names.
// Sections should not be created in place - New/ExistingSection/Section functions.
// SectionT is the type name because we need the more user-visible `Section` as a function name.
type SectionT struct {
	// Flags holds flags available to this Section and all subsections and Commands
	Flags flag.FlagMap
	// Commands holds the Commands under this Section
	Commands command.CommandMap
	// Sections holds the Sections under this Section
	Sections SectionMap
	// HelpShort is a required one-line descripiton of this section
	HelpShort HelpShort
	// HelpLong is an optional longer description of this section
	HelpLong string
	// Footer is yet another optional longer description.
	Footer string
}

// New creates a Section!
func New(helpShort HelpShort, opts ...SectionOpt) SectionT {
	section := SectionT{
		HelpShort: helpShort,
		Flags:     make(flag.FlagMap),
		Sections:  make(SectionMap),
		Commands:  make(command.CommandMap),
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

// ExistingSection adds an existing Section underneath this Section. Panics if a Section with the same name already exists
func ExistingSection(name Name, value SectionT) SectionOpt {
	return func(app *SectionT) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// ExistingCommand adds an existing Command underneath this Section. Panics if a Command with the same name already exists
func ExistingCommand(name command.Name, value command.Command) SectionOpt {
	return func(app *SectionT) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// ExistingFlag adds an existing Flag to be made availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
func ExistingFlag(name flag.Name, value flag.Flag) SectionOpt {
	if !strings.HasPrefix(string(name), "-") {
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
		if !strings.HasPrefix(string(name), "-") {
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
func Section(name Name, helpShort HelpShort, opts ...SectionOpt) SectionOpt {
	return ExistingSection(name, New(helpShort, opts...))
}

// Flag creates a Flag and makes it availabe to subsections and subcommands. Panics if the flag name doesn't start with '-' or a flag with the same name exists already
func Flag(name flag.Name, helpShort flag.HelpShort, empty value.EmptyConstructor, opts ...flag.FlagOpt) SectionOpt {
	return ExistingFlag(name, flag.New(helpShort, empty, opts...))
}

// Command creates a Command and adds it underneath this Section. Panics if a Command with the same name already exists
func Command(name command.Name, helpShort command.HelpShort, action command.Action, opts ...command.CommandOpt) SectionOpt {
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
