package flag

import (
	"log"
	"sort"

	"go.bbkane.com/warg/value"
)

// Name of a flag
type Name string

// HelpShort is a description of what this flag does.
type HelpShort string

// FlagMap holds flags - used by Commands and Sections
type FlagMap map[Name]Flag

func (fm *FlagMap) SortedNames() []Name {
	keys := make([]Name, 0, len(*fm))
	for k := range *fm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	return keys
}

// AddFlag adds a new flag and panics if it already exists
func (fm FlagMap) AddFlag(name Name, value Flag) {
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

// FlagOpt customizes a Flag on creation
type FlagOpt func(*Flag)

type Flag struct {
	// Alias is an alternative name for a flag, usually shorter :)
	Alias Name

	// ConfigPath is the path from the config to the value the flag updates
	ConfigPath string

	// Envvars holds a list of environment variables to update this flag. Only the first one that exists will be used.
	EnvVars []string

	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor value.EmptyConstructor

	// HelpShort is a message for the user on how to use this flag
	HelpShort HelpShort

	// Required means the user MUST fill this flag
	Required bool

	// When UnsetSentinal is passed as a flag value, Value is reset and SetBy is set to ""
	UnsetSentinel string

	// -- the following are set when parsing

	// IsCommandFlag is set when parsing. Set to true if the flag was attached to a command (as opposed to being inherited from a section)
	IsCommandFlag bool

	// SetBy might be set when parsing. Possible values: appdefault, config, passedflag
	SetBy string

	// Value is set when parsing. Use SetBy != "" to determine whether a value was actually passed  instead of being empty
	Value value.Value
}

// New creates a Flag with options!
func New(helpShort HelpShort, empty value.EmptyConstructor, opts ...FlagOpt) Flag {
	flag := Flag{
		HelpShort:             helpShort,
		EmptyValueConstructor: empty,
		Alias:                 "",
		ConfigPath:            "",
		EnvVars:               nil,
		Required:              false,
		IsCommandFlag:         false,
		SetBy:                 "",
		UnsetSentinel:         "",
		Value:                 nil,
	}
	for _, opt := range opts {
		opt(&flag)
	}
	return flag
}

// Alias is an alternative name for a flag, usually shorter :)
func Alias(alias Name) FlagOpt {
	return func(f *Flag) {
		f.Alias = alias
	}
}

// ConfigPath adds a configpath to a flag
func ConfigPath(path string) FlagOpt {
	return func(flag *Flag) {
		flag.ConfigPath = path
	}
}

// EnvVars adds a list of environmental variables to search through to update this flag. The first one that exists will be used to update the flag. Further existing envvars will be ignored.
func EnvVars(name ...string) FlagOpt {
	return func(f *Flag) {
		f.EnvVars = name
	}
}

// UnsetSentinel is a bit of an advanced feature meant to allow overriding a
// default, config, or environmental value with a command line flag.
// When UnsetSentinel is passed as a flag value, Value is reset and SetBy is set to "".
// It it recommended to set `name` to "UNSET" for consistency among warg apps.
// Scalar example:
//
//	app --flag UNSET  // undoes anything that sets --flag
//
// Slice example:
//
//	app --flag a --flag b --flag UNSET --flag c --flag d // ends up with []string{"c", "d"}
func UnsetSentinel(name string) FlagOpt {
	return func(f *Flag) {
		f.UnsetSentinel = name
	}
}

// Required means the user MUST fill this flag
func Required() FlagOpt {
	return func(f *Flag) {
		f.Required = true
	}
}
