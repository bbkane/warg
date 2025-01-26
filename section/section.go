package section

import (
	"log"
	"sort"

	"go.bbkane.com/warg/command"
)

// Name of the section
type Name string

// HelpShort is a required short description of the section
type HelpShort string

// SectionMap holds Sections - used by other Sections
type SectionMap map[Name]SectionT

func (fm SectionMap) Empty() bool {
	return len(fm) == 0
}

func (fm SectionMap) SortedNames() []Name {
	keys := make([]Name, 0, len(fm))
	for k := range fm {
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
		Sections:  make(SectionMap),
		Commands:  make(command.CommandMap),
		HelpLong:  "",
		Footer:    "",
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

// ExistingSections adds existing Sections underneath this Section. Panics if a Section with the same name already exists
func ExistingSections(sections SectionMap) SectionOpt {
	return func(app *SectionT) {
		for name, value := range sections {
			ExistingSection(name, value)(app)
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

// ExistingCommands adds existing Commands underneath this Section. Panics if a Command with the same name already exists
func ExistingCommands(commands command.CommandMap) SectionOpt {
	return func(app *SectionT) {
		for name, value := range commands {
			ExistingCommand(name, value)(app)
		}
	}
}

// Section creates a Section and adds it underneath this Section. Panics if a Section with the same name already exists
func Section(name Name, helpShort HelpShort, opts ...SectionOpt) SectionOpt {
	return ExistingSection(name, New(helpShort, opts...))
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

// FlatSection represents a section and relevant parent information
type FlatSection struct {

	// Path to this section
	Path []Name
	// Sec is this section
	Sec SectionT
}

// Breadthfirst returns a SectionIterator that yields sections sorted alphabetically breadth-first by path.
// Yielded sections should never be modified - they can share references to the same inherited flags
// SectionIterator's Next() method panics if two sections in the path have flags with the same name.
// Breadthfirst is used by app.Validate and help.AllCommandCommandHelp/help.AllCommandSectionHelp
func (sec *SectionT) BreadthFirst(path []Name) SectionIterator {

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
		childPath := make([]Name, len(current.Path)+1)
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
