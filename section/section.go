package section

import (
	"log"

	"go.bbkane.com/warg"
)

// Section adds an existing Section underneath this Section. Panics if a Section with the same name already exists
func Section(name string, value warg.SectionT) warg.SectionOpt {
	return func(app *warg.SectionT) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// SectionMap adds existing Sections underneath this Section. Panics if a Section with the same name already exists
func SectionMap(sections warg.SectionMapT) warg.SectionOpt {
	return func(app *warg.SectionT) {
		for name, value := range sections {
			Section(name, value)(app)
		}
	}
}

// Command adds an existing Command underneath this Section. Panics if a Command with the same name already exists
func Command(name string, value warg.Command) warg.SectionOpt {
	return func(app *warg.SectionT) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// CommandMap adds existing Commands underneath this Section. Panics if a Command with the same name already exists
func CommandMap(commands warg.CommandMap) warg.SectionOpt {
	return func(app *warg.SectionT) {
		for name, value := range commands {
			Command(name, value)(app)
		}
	}
}

// NewSection creates a NewSection and adds it underneath this NewSection. Panics if a NewSection with the same name already exists
func NewSection(name string, helpShort string, opts ...warg.SectionOpt) warg.SectionOpt {
	return Section(name, warg.NewSection(helpShort, opts...))
}

// NewCommand creates a NewCommand and adds it underneath this Section. Panics if a NewCommand with the same name already exists
func NewCommand(name string, helpShort string, action warg.Action, opts ...warg.CommandOpt) warg.SectionOpt {
	return Command(name, warg.NewCommand(helpShort, action, opts...))
}

// Footer adds an optional help string to this Section
func Footer(footer string) warg.SectionOpt {
	return func(cat *warg.SectionT) {
		cat.Footer = footer
	}
}

// HelpLong adds an optional help string to this Section
func HelpLong(helpLong string) warg.SectionOpt {
	return func(cat *warg.SectionT) {
		cat.HelpLong = helpLong
	}
}
