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

	// DefaultValues will be shoved into Value if the app builder specifies it.
	// For scalar values, the last DefaultValues wins
	DefaultValues []string

	// Envvars holds a list of environment variables to update this flag. Only the first one that exists will be used.
	EnvVars []string

	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor value.EmptyConstructor

	// HelpShort is a message for the user on how to use this flag
	HelpShort HelpShort

	// Required means the user MUST fill this flag
	Required bool

	// -- the following are set when parsing

	// IsCommandFlag is set when parsing. Set to true if the flag was attached to a command (as opposed to being inherited from a section)
	IsCommandFlag bool

	// SetBy might be set when parsing. Possible values: appdefault, config, passedflag
	SetBy string

	// TypeDescription is set when parsing. Describes the type: int, string, ...
	TypeDescription string

	// TypeInfo is set when parsing. Describes the "shape" of the type
	TypeInfo value.TypeContainer

	// Value might be set when parsing. The interface returned by updating a flag
	Value value.Value
}

// New creates a Flag with options!
func New(helpShort HelpShort, empty value.EmptyConstructor, opts ...FlagOpt) Flag {
	flag := Flag{
		HelpShort:             helpShort,
		EmptyValueConstructor: empty,
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

// Default adds default values to a flag. The flag will be updated with each of the values when Resolve is called.
// Panics when multiple values are passed and the flags is scalar
func Default(values ...string) FlagOpt {
	return func(flag *Flag) {
		empty, err := flag.EmptyValueConstructor()
		if err != nil {
			log.Panicf("cannot create empty flag value when checking default: %v", flag)
		}
		if empty.TypeInfo() == value.TypeContainerScalar && len(values) != 1 {
			log.Panicf("a scalar flag should only have one default value: We don't know the name of the type, but here's the Help: %#v", flag.HelpShort)
		}
		flag.DefaultValues = values
	}
}

// EnvVars adds a list of environmental variables to search through to update this flag. The first one that exists will be used to update the flag. Further existing envvars will be ignored.
func EnvVars(name ...string) FlagOpt {
	return func(f *Flag) {
		f.EnvVars = name
	}
}

// Required means the user MUST fill this flag
func Required() FlagOpt {
	return func(f *Flag) {
		f.Required = true
	}
}
