package section

import (
	"log"

	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
)

// SectionOpt customizes a Section on creation
type SectionOpt func(*cli.Section)

// New creates a standalone [cli.Section]. All section options are in the [go.bbkane.com/warg/section] package
func New(helpShort string, opts ...SectionOpt) cli.Section {
	section := cli.Section{
		HelpShort: helpShort,
		Sections:  make(cli.SectionMap),
		Commands:  make(cli.CommandMap),
		HelpLong:  "",
		Footer:    "",
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

// Section adds an existing Section as a child of this Section. Panics if a Section with the same name already exists
func Section(name string, value cli.Section) SectionOpt {
	return func(app *cli.Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// SectionMap adds existing Sections as a child of this Section. Panics if a Section with the same name already exists
func SectionMap(sections cli.SectionMap) SectionOpt {
	return func(app *cli.Section) {
		for name, value := range sections {
			Section(name, value)(app)
		}
	}
}

// Command adds an existing Command as a child of this Section. Panics if a Command with the same name already exists
func Command(name string, value cli.Command) SectionOpt {
	return func(app *cli.Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// CommandMap adds existing Commands as a child of this Section. Panics if a Command with the same name already exists
func CommandMap(commands cli.CommandMap) SectionOpt {
	return func(app *cli.Section) {
		for name, value := range commands {
			Command(name, value)(app)
		}
	}
}

// NewSection creates a new Section as a child of this Section. Panics if a NewSection with the same name already exists
func NewSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return Section(name, New(helpShort, opts...))
}

// NewCommand creates a new Command as a child of this Section. Panics if a NewCommand with the same name already exists
func NewCommand(name string, helpShort string, action cli.Action, opts ...command.CommandOpt) SectionOpt {
	return Command(name, command.New(helpShort, action, opts...))
}

// Footer adds an optional help string to this Section
func Footer(footer string) SectionOpt {
	return func(cat *cli.Section) {
		cat.Footer = footer
	}
}

// HelpLong adds an optional help string to this Section
func HelpLong(helpLong string) SectionOpt {
	return func(cat *cli.Section) {
		cat.HelpLong = helpLong
	}
}
