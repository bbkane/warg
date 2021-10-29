package flag

import (
	"log"
	"strings"

	v "github.com/bbkane/warg/value"
)

// FlagMap holds flags - used by Commands and Sections
type FlagMap = map[string]Flag

// FlagOpt customizes a Flag on creation
type FlagOpt = func(*Flag)

// PassedFlags holds a map of flag names to flag Values and is passed to a command's Action
type PassedFlags = map[string]interface{}

type Flag struct {
	// Alias is an alternative name for a flag, usually shorter :)
	Alias string

	// ConfigPath is the path from the config to the value the flag updates
	ConfigPath string

	// DefaultValues will be shoved into Value if the app builder specifies it.
	// For scalar values, the last DefaultValues wins
	DefaultValues []string

	// Envvars holds a list of environment variables to update this flag. Only the first one that exists will be used.
	EnvVars []string

	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor v.EmptyConstructor

	// Help is a message for the user on how to use this flag
	Help string

	// Required means the user MUST fill this flag
	Required bool

	// -- the following are set when parsing

	// IsCommandFlag is set when parsing. Set to true if the flag was attached to a command (as opposed to being inherited from a section)
	IsCommandFlag bool

	// SetBy might be set when parsing. Possible values: appdefault, config, passedflag
	SetBy string

	// TypeDescription is set when parsing. Describes the type: int, string, ...
	TypeDescription string

	// Value might be set when parsing. The interface returned by updating a flag
	Value v.Value
}

// New creates a Flag with options!
func New(helpShort string, empty v.EmptyConstructor, opts ...FlagOpt) Flag {
	flag := Flag{
		Help:                  helpShort,
		EmptyValueConstructor: empty,
	}
	for _, opt := range opts {
		opt(&flag)
	}
	return flag
}

// Alias is an alternative name for a flag, usually shorter :)
func Alias(alias string) FlagOpt {
	if !strings.HasPrefix(alias, "-") {
		log.Panicf("All aliases should start with '-': %v", alias)
	}
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
		if empty.TypeInfo() == v.TypeInfoScalar && len(values) != 1 {
			log.Panicf("a scalar flag should only have one default value: We don't know the name of the type, but here's the Help: %#v", flag.Help)
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
