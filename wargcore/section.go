package wargcore

import (
	"log"
)

// SectionOpt customizes a Section on creation
type SectionOpt func(*Section)

// NewSection creates a standalone [Section]. All section options are in the [go.bbkane.com/warg/section] package
func NewSection(helpShort string, opts ...SectionOpt) Section {
	section := Section{
		HelpShort: helpShort,
		Sections:  make(SectionMap),
		Commands:  make(CmdMap),
		HelpLong:  "",
		Footer:    "",
	}
	for _, opt := range opts {
		opt(&section)
	}
	return section
}

// ChildSection adds an existing ChildSection as a child of this ChildSection. Panics if a ChildSection with the same name already exists
func ChildSection(name string, value Section) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// ChildSectionMap adds existing Sections as a child of this Section. Panics if a Section with the same name already exists
func ChildSectionMap(sections SectionMap) SectionOpt {
	return func(app *Section) {
		for name, value := range sections {
			ChildSection(name, value)(app)
		}
	}
}

// ChildCmd adds an existing ChildCmd as a child of this Section. Panics if a ChildCmd with the same name already exists
func ChildCmd(name string, value Cmd) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// ChildCmdMap adds existing Commands as a child of this Section. Panics if a Command with the same name already exists
func ChildCmdMap(commands CmdMap) SectionOpt {
	return func(app *Section) {
		for name, value := range commands {
			ChildCmd(name, value)(app)
		}
	}
}

// NewChildSection creates a new Section as a child of this Section. Panics if a NewChildSection with the same name already exists
func NewChildSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return ChildSection(name, NewSection(helpShort, opts...))
}

// NewChildCmd creates a new Command as a child of this Section. Panics if a NewChildCmd with the same name already exists
func NewChildCmd(name string, helpShort string, action Action, opts ...CommandOpt) SectionOpt {
	return ChildCmd(name, NewCmd(helpShort, action, opts...))
}

// SectionFooter adds an optional help string to this Section
func SectionFooter(footer string) SectionOpt {
	return func(cat *Section) {
		cat.Footer = footer
	}
}

// SectionHelpLong adds an optional help string to this Section
func SectionHelpLong(helpLong string) SectionOpt {
	return func(cat *Section) {
		cat.HelpLong = helpLong
	}
}
