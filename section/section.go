package section

import (
	"log"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/wargcore"
)

// SectionOpt customizes a Section on creation
type SectionOpt func(*wargcore.Section)

// New creates a standalone [wargcore.Section]. All section options are in the [go.bbkane.com/warg/section] package
func New(helpShort string, opts ...SectionOpt) wargcore.Section {
	section := wargcore.Section{
		HelpShort: helpShort,
		Sections:  make(wargcore.SectionMap),
		Commands:  make(wargcore.CmdMap),
		HelpLong:  "",
		Footer:    "",
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

// Section adds an existing Section as a child of this Section. Panics if a Section with the same name already exists
func Section(name string, value wargcore.Section) SectionOpt {
	return func(app *wargcore.Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// SectionMap adds existing Sections as a child of this Section. Panics if a Section with the same name already exists
func SectionMap(sections wargcore.SectionMap) SectionOpt {
	return func(app *wargcore.Section) {
		for name, value := range sections {
			Section(name, value)(app)
		}
	}
}

// Command adds an existing Command as a child of this Section. Panics if a Command with the same name already exists
func Command(name string, value wargcore.Cmd) SectionOpt {
	return func(app *wargcore.Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// CommandMap adds existing Commands as a child of this Section. Panics if a Command with the same name already exists
func CommandMap(commands wargcore.CmdMap) SectionOpt {
	return func(app *wargcore.Section) {
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
func NewCommand(name string, helpShort string, action wargcore.Action, opts ...command.CommandOpt) SectionOpt {
	return Command(name, command.New(helpShort, action, opts...))
}

// Footer adds an optional help string to this Section
func Footer(footer string) SectionOpt {
	return func(cat *wargcore.Section) {
		cat.Footer = footer
	}
}

// HelpLong adds an optional help string to this Section
func HelpLong(helpLong string) SectionOpt {
	return func(cat *wargcore.Section) {
		cat.HelpLong = helpLong
	}
}
