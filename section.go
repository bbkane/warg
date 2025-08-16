package warg

import (
	"log"
	"sort"
)

// SectionOpt customizes a Section on creation
type SectionOpt func(*Section)

// NewSection creates a standalone [Section]. All section options are in the [go.bbkane.com/warg/section] package
func NewSection(helpShort string, opts ...SectionOpt) Section {
	section := Section{
		HelpShort: helpShort,
		Sections:  make(SectionMap),
		Cmds:      make(CmdMap),
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
		if _, alreadyThere := app.Cmds[name]; !alreadyThere {
			app.Cmds[name] = value
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

// SectionMap holds Sections - used by other Sections
type SectionMap map[string]Section

func (fm SectionMap) Empty() bool {
	return len(fm) == 0
}

func (fm SectionMap) SortedNames() []string {
	keys := make([]string, 0, len(fm))
	for k := range fm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	return keys
}

// Sections are like "folders" for Commmands.
// They should usually have noun names.
// Sections should not be be created directly, but with the APIs in [go.bbkane.com/warg/section].
type Section struct {
	// Cmds holds the Cmds under this Section
	Cmds CmdMap
	// Sections holds the Sections under this Section
	Sections SectionMap
	// HelpShort is a required one-line descripiton of this section
	HelpShort string
	// HelpLong is an optional longer description of this section
	HelpLong string
	// Footer is yet another optional longer description.
	Footer string
}

// FlatSection represents a section and relevant parent information
type FlatSection struct {

	// Path to this section
	Path []string
	// Sec is this section
	Sec Section
}

// Breadthfirst returns a SectionIterator that yields sections sorted alphabetically breadth-first by path.
// Yielded sections should never be modified - they can share references to the same inherited flags
// SectionIterator's Next() method panics if two sections in the path have flags with the same name.
// Breadthfirst is used by app.Validate and help.AllCommandCommandHelp/help.AllCommandSectionHelp
func (sec *Section) BreadthFirst(path []string) SectionIterator {

	queue := make([]FlatSection, 0, 1)
	queue = append(queue, FlatSection{
		Path: path,
		Sec:  *sec,
	})

	return SectionIterator{
		queue: queue,
	}
}

// SectionIterator is used in BreadthFirst. See BreadthFirst docs
type SectionIterator struct {
	queue []FlatSection
}

// HasNext is used in BreadthFirst. See BreadthFirst docs
func (s *SectionIterator) Next() FlatSection {
	current := s.queue[0]
	s.queue = s.queue[1:]

	// Add child sections to queue
	for _, childName := range current.Sec.Sections.SortedNames() {

		// child.Path = current.Path + child.name
		childPath := make([]string, len(current.Path)+1)
		copy(childPath, current.Path)
		childPath[len(childPath)-1] = childName

		s.queue = append(s.queue, FlatSection{
			Path: childPath,
			Sec:  current.Sec.Sections[childName],
		})
	}

	return current
}

// HasNext is used in BreadthFirst. See BreadthFirst docs
func (s *SectionIterator) HasNext() bool {
	return len(s.queue) > 0
}
