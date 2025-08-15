package warg

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

// SubSection adds an existing SubSection as a child of this SubSection. Panics if a SubSection with the same name already exists
func SubSection(name string, value Section) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// SubSectionMap adds existing Sections as a child of this Section. Panics if a Section with the same name already exists
func SubSectionMap(sections SectionMap) SectionOpt {
	return func(app *Section) {
		for name, value := range sections {
			SubSection(name, value)(app)
		}
	}
}

// SubCmd adds an existing SubCmd as a child of this Section. Panics if a SubCmd with the same name already exists
func SubCmd(name string, value Cmd) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// SubCmdMap adds existing Commands as a child of this Section. Panics if a Command with the same name already exists
func SubCmdMap(commands CmdMap) SectionOpt {
	return func(app *Section) {
		for name, value := range commands {
			SubCmd(name, value)(app)
		}
	}
}

// NewSubSection creates a new Section as a child of this Section. Panics if a NewSubSection with the same name already exists
func NewSubSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return SubSection(name, NewSection(helpShort, opts...))
}

// NewSubCmd creates a new Command as a child of this Section. Panics if a NewSubCmd with the same name already exists
func NewSubCmd(name string, helpShort string, action Action, opts ...CmdOpt) SectionOpt {
	return SubCmd(name, NewCmd(helpShort, action, opts...))
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
