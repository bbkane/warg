package section

import (
	"log"

	"go.bbkane.com/warg/command"
)

// SectionOpt customizes a Section on creation
type SectionOpt func(*SectionT)

// NewSectionT creates a Section!
func NewSectionT(helpShort string, opts ...SectionOpt) SectionT {
	section := SectionT{
		HelpShort: helpShort,
		Sections:  make(SectionMapT),
		Commands:  make(command.CommandMap),
		HelpLong:  "",
		Footer:    "",
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

// Section adds an existing Section underneath this Section. Panics if a Section with the same name already exists
func Section(name string, value SectionT) SectionOpt {
	return func(app *SectionT) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// SectionMap adds existing Sections underneath this Section. Panics if a Section with the same name already exists
func SectionMap(sections SectionMapT) SectionOpt {
	return func(app *SectionT) {
		for name, value := range sections {
			Section(name, value)(app)
		}
	}
}

// Command adds an existing Command underneath this Section. Panics if a Command with the same name already exists
func Command(name string, value command.Command) SectionOpt {
	return func(app *SectionT) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// CommandMap adds existing Commands underneath this Section. Panics if a Command with the same name already exists
func CommandMap(commands command.CommandMap) SectionOpt {
	return func(app *SectionT) {
		for name, value := range commands {
			Command(name, value)(app)
		}
	}
}

// NewSection creates a NewSection and adds it underneath this NewSection. Panics if a NewSection with the same name already exists
func NewSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return Section(name, NewSectionT(helpShort, opts...))
}

// NewCommand creates a NewCommand and adds it underneath this Section. Panics if a NewCommand with the same name already exists
func NewCommand(name string, helpShort string, action command.Action, opts ...command.CommandOpt) SectionOpt {
	return Command(name, command.NewCommand(helpShort, action, opts...))
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
