package warg

import (
	"log"
	"sort"
)

// SectionOpt is a functional option for configuring a [Section] during creation.
type SectionOpt func(*Section)

// NewSection creates a standalone [Section] that groups commands and child sections.
// Attach it to a parent with [SubSection] or pass it directly to [New] as the root section.
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

// SubSection attaches an existing [Section] as a child. Panics if a section with the same name exists.
func SubSection(name string, value Section) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Sections[name]; !alreadyThere {
			app.Sections[name] = value
		} else {
			log.Panicf("section already exists: %#v\n", name)
		}
	}
}

// SubSectionMap attaches multiple existing sections as children. Panics if any name already exists.
func SubSectionMap(sections SectionMap) SectionOpt {
	return func(app *Section) {
		for name, value := range sections {
			SubSection(name, value)(app)
		}
	}
}

// SubCmd attaches an existing [Cmd] as a child of this section. Panics if a command with the same name exists.
func SubCmd(name string, value Cmd) SectionOpt {
	return func(app *Section) {
		if _, alreadyThere := app.Cmds[name]; !alreadyThere {
			app.Cmds[name] = value
		} else {
			log.Panicf("command already exists: %#v\n", name)
		}
	}
}

// SubCmdMap attaches multiple existing commands as children. Panics if any name already exists.
func SubCmdMap(commands CmdMap) SectionOpt {
	return func(app *Section) {
		for name, value := range commands {
			SubCmd(name, value)(app)
		}
	}
}

// NewSubSection creates a new child [Section] with the given name and options.
// Panics if a section with the same name already exists.
func NewSubSection(name string, helpShort string, opts ...SectionOpt) SectionOpt {
	return SubSection(name, NewSection(helpShort, opts...))
}

// NewSubCmd creates a new [Cmd] and attaches it as a child of this section.
// Panics if a command with the same name already exists.
func NewSubCmd(name string, helpShort string, action Action, opts ...CmdOpt) SectionOpt {
	return SubCmd(name, NewCmd(helpShort, action, opts...))
}

// SectionFooter sets an optional footer text displayed at the end of help output for this section.
func SectionFooter(footer string) SectionOpt {
	return func(cat *Section) {
		cat.Footer = footer
	}
}

// SectionHelpLong sets an optional extended description for this section, shown in detailed help.
func SectionHelpLong(helpLong string) SectionOpt {
	return func(cat *Section) {
		cat.HelpLong = helpLong
	}
}

// SectionMap maps section names to [Section] instances.
type SectionMap map[string]Section

// Empty reports whether the map contains no sections.
func (fm SectionMap) Empty() bool {
	return len(fm) == 0
}

// SortedNames returns the section names in alphabetical order.
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

// Section groups related commands and child sections, forming the hierarchical
// structure of a CLI app. Section names should be nouns (e.g., "config", "users").
// Do not construct directly; use [NewSection] or [NewSubSection].
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

// flatSection represents a section and relevant parent information
type flatSection struct {

	// Path to this section
	Path []string
	// Sec is this section
	Sec Section
}

// Breadthfirst returns a SectionIterator that yields sections sorted alphabetically breadth-first by path.
// Yielded sections should never be modified - they can share references to the same inherited flags
// SectionIterator's Next() method panics if two sections in the path have flags with the same name.
// Breadthfirst is used by app.Validate and help.AllCommandCommandHelp/help.AllCommandSectionHelp
func (sec *Section) breadthFirst(path []string) sectionIterator {

	queue := make([]flatSection, 0, 1)
	queue = append(queue, flatSection{
		Path: path,
		Sec:  *sec,
	})

	return sectionIterator{
		queue: queue,
	}
}

// depthFirstSections returns all sections in depth-first pre-order: the current section first,
// then each child section (sorted alphabetically) and its descendants recursively.
// This ensures a section's own commands appear before its siblings' commands in help output.
// See https://github.com/bbkane/warg/issues/74
func depthFirstSections(sec Section, path []string) []flatSection {
	result := []flatSection{{Path: path, Sec: sec}}
	for _, childName := range sec.Sections.SortedNames() {
		childPath := append(append([]string(nil), path...), childName)
		result = append(result, depthFirstSections(sec.Sections[childName], childPath)...)
	}
	return result
}

// sectionIterator is used in BreadthFirst. See BreadthFirst docs
type sectionIterator struct {
	queue []flatSection
}

// HasNext is used in BreadthFirst. See BreadthFirst docs
func (s *sectionIterator) Next() flatSection {
	current := s.queue[0]
	s.queue = s.queue[1:]

	// Add child sections to queue
	for _, childName := range current.Sec.Sections.SortedNames() {

		// child.Path = current.Path + child.name
		childPath := make([]string, len(current.Path)+1)
		copy(childPath, current.Path)
		childPath[len(childPath)-1] = childName

		s.queue = append(s.queue, flatSection{
			Path: childPath,
			Sec:  current.Sec.Sections[childName],
		})
	}

	return current
}

// HasNext is used in BreadthFirst. See BreadthFirst docs
func (s *sectionIterator) HasNext() bool {
	return len(s.queue) > 0
}
