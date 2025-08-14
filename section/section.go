package section

import (
	"log"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/wargcore"
)

// SectionOpt customizes a Section on creation
type SectionOpt func(*wargcore.Section)

// NewSection creates a standalone [wargcore.Section]. All section options are in the [go.bbkane.com/warg/section] package
func NewSection(helpShort string, opts ...SectionOpt) wargcore.Section {
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

// ChildSection adds an existing ChildSection as a child of this ChildSection. Panics if a ChildSection with the same name already exists
func ChildSection(name string, value wargcore.Section) SectionOpt {
	return func(app *wargcore.Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// ChildSectionMap adds existing Sections as a child of this Section. Panics if a Section with the same name already exists
func ChildSectionMap(sections wargcore.SectionMap) SectionOpt {
	return func(app *wargcore.Section) {
		for name, value := range sections {
			ChildSection(name, value)(app)
		}
	}
}

// ChildCmd adds an existing ChildCmd as a child of this Section. Panics if a ChildCmd with the same name already exists
func ChildCmd(name string, value wargcore.Cmd) SectionOpt {
	return func(app *wargcore.Section) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// ChildCmdMap adds existing Commands as a child of this Section. Panics if a Command with the same name already exists
func ChildCmdMap(commands wargcore.CmdMap) SectionOpt {
	return func(app *wargcore.Section) {
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
func NewChildCmd(name string, helpShort string, action wargcore.Action, opts ...command.CommandOpt) SectionOpt {
	return ChildCmd(name, command.NewCmd(helpShort, action, opts...))
}

// SectionFooter adds an optional help string to this Section
func SectionFooter(footer string) SectionOpt {
	return func(cat *wargcore.Section) {
		cat.Footer = footer
	}
}

// SectionHelpLong adds an optional help string to this Section
func SectionHelpLong(helpLong string) SectionOpt {
	return func(cat *wargcore.Section) {
		cat.HelpLong = helpLong
	}
}
