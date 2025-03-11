package section

import (
	"log"

	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
)

// SectionOpt customizes a Section on creation
type SectionOpt func(*cli.SectionT)

// NewSectionT creates a Section!
func NewSectionT(helpShort string, opts ...SectionOpt) cli.SectionT {
	section := cli.SectionT{
		HelpShort: helpShort,
		Sections:  make(cli.SectionMapT),
		Commands:  make(cli.CommandMap),
		HelpLong:  "",
		Footer:    "",
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

// Section adds an existing Section underneath this Section. Panics if a Section with the same name already exists
func Section(name string, value cli.SectionT) SectionOpt {
	return func(app *cli.SectionT) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// SectionMap adds existing Sections underneath this Section. Panics if a Section with the same name already exists
func SectionMap(sections cli.SectionMapT) SectionOpt {
	return func(app *cli.SectionT) {
		for name, value := range sections {
			Section(name, value)(app)
		}
	}
}

// Command adds an existing Command underneath this Section. Panics if a Command with the same name already exists
func Command(name string, value cli.Command) SectionOpt {
	return func(app *cli.SectionT) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// CommandMap adds existing Commands underneath this Section. Panics if a Command with the same name already exists
func CommandMap(commands cli.CommandMap) SectionOpt {
	return func(app *cli.SectionT) {
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
func NewCommand(name string, helpShort string, action cli.Action, opts ...command.CommandOpt) SectionOpt {
	return Command(name, command.NewCommand(helpShort, action, opts...))
}

// Footer adds an optional help string to this Section
func Footer(footer string) SectionOpt {
	return func(cat *cli.SectionT) {
		cat.Footer = footer
	}
}

// HelpLong adds an optional help string to this Section
func HelpLong(helpLong string) SectionOpt {
	return func(cat *cli.SectionT) {
		cat.HelpLong = helpLong
	}
}
