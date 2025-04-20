package cli

import (
	"log"
	"sort"

	"go.bbkane.com/warg/value"
)

// FlagMap holds flags - used by Commands and Sections
type FlagMap map[string]Flag

func (fm *FlagMap) SortedNames() []string {
	keys := make([]string, 0, len(*fm))
	for k := range *fm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	return keys
}

// AddFlag adds a new flag and panics if it already exists
func (fm FlagMap) AddFlag(name string, value Flag) {
	if _, alreadyThere := (fm)[name]; !alreadyThere {
		(fm)[name] = value
	} else {
		log.Panicf("flag already exists: %#v\n", name)
	}
}

// AddFlags adds another FlagMap to this one and  and panics if a flag name already exists
func (fm FlagMap) AddFlags(flagMap FlagMap) {
	for name, value := range flagMap {
		fm.AddFlag(name, value)
	}
}

type Flag struct {
	// Alias is an alternative name for a flag, usually shorter :)
	Alias string

	// CompletionCandidates is a function that returns a list of completion candidates for this flag.
	// Note that some flags in the cli.Context Flags map may not be set, even if they're required.
	// TODO: get a comprehensive list of restrictions on the context.
	CompletionCandidates CompletionCandidates

	// ConfigPath is the path from the config to the value the flag updates
	ConfigPath string

	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor value.EmptyConstructor

	// Envvars holds a list of environment variables to update this flag. Only the first one that exists will be used.
	EnvVars []string

	// HelpShort is a message for the user on how to use this flag
	HelpShort string

	// Required means the user MUST fill this flag
	Required bool

	// When UnsetSentinal is passed as a flag value, Value is reset and SetBy is set to ""
	UnsetSentinel string

	// -- the following are set when parsing and they're all deprecated for my march to immutabality

	// IsCommandFlag is set when parsing. Set to true if the flag was attached to a command (as opposed to being inherited from a section)
	// Deprecated: Check if the name is in command.Flags instead
	IsCommandFlag bool

	// Value is set when parsing. Use SetBy != "" to determine whether a value was actually passed  instead of being empty
	// Deprecated: check the value from ParseState.FlagValues
	Value value.Value
}
