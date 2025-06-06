package wargcore

import (
	"context"
	"os"
	"sort"
)

// PassedFlags holds a map of flag names to flag Values
type PassedFlags map[string]interface{} // This can just stay a string for the convenience of the user.

// Context holds everything a command needs.
type Context struct {
	App   *App
	Flags PassedFlags

	ParseState *ParseState

	// Context to smuggle user-defined state (i.e., not flags) into an Action. I use this for mocks when testing
	Context context.Context

	Stderr *os.File
	Stdout *os.File
}

// An Action is run as the result of a command
type Action func(Context) error

// A CommandMap holds Commands and is used by Sections
type CommandMap map[string]Command

func (fm CommandMap) Empty() bool {
	return len(fm) == 0
}

func (fm CommandMap) SortedNames() []string {
	keys := make([]string, 0, len(fm))
	for k := range fm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	return keys
}

// A Command will run code for you!
// The name of a Command should probably be a verb - add , edit, run, ...
// A Command should not be constructed directly. Use Command / New / ExistingCommand functions
type Command struct {
	// Action to run when command is invoked
	Action Action

	// CompletionCandidates is a function that returns a list of completion candidates for this commmand.
	// Note that some flags in the cli.Context Flags map may not be set, even if they're required.
	// TODO: get a comprehensive list of restrictions on the context.
	CompletionCandidates CompletionCandidates

	// Parsed Flags
	Flags FlagMap

	// Footer is yet another optional longer description.
	Footer string

	// HelpLong is an optional longer description
	HelpLong string

	// HelpShort is a required one-line description
	HelpShort string
}
