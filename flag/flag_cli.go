package flag

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

	// ConfigPath is the path from the config to the value the flag updates
	ConfigPath string

	// Envvars holds a list of environment variables to update this flag. Only the first one that exists will be used.
	EnvVars []string

	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor value.EmptyConstructor

	// HelpShort is a message for the user on how to use this flag
	HelpShort string

	// Required means the user MUST fill this flag
	Required bool

	// When UnsetSentinal is passed as a flag value, Value is reset and SetBy is set to ""
	UnsetSentinel string

	// -- the following are set when parsing

	// IsCommandFlag is set when parsing. Set to true if the flag was attached to a command (as opposed to being inherited from a section)
	IsCommandFlag bool

	// Value is set when parsing. Use SetBy != "" to determine whether a value was actually passed  instead of being empty
	Value value.Value
}
