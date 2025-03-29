package cli

import (
	"sort"
)

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
// They should have noun names.
type Section struct {
	// Commands holds the Commands under this Section
	Commands CommandMap
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
