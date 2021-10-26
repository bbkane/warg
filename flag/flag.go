package flag

import (
	"fmt"
	"log"

	"github.com/bbkane/warg/configreader"
	v "github.com/bbkane/warg/value"
)

// FlagMap holds flags - used by Commands and Sections
type FlagMap = map[string]Flag

// FlagOpt customizes a Flag on creation
type FlagOpt = func(*Flag)

// Look up keys (meant for environment variable parsing) - fulfillable with os.LookupEnv or warg.DictLookup(map)
type LookupFunc = func(key string) (string, bool)

// PassedFlags holds a map of flag names to flag Values and is passed to a command's Action
type PassedFlags = map[string]interface{}

type Flag struct {
	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor v.EmptyConstructor
	// ConfigPath is the path from the config to the value the flag updates
	ConfigPath string
	// DefaultValues will be shoved into Value if the app builder specifies it.
	// For scalar values, the last DefaultValues wins
	DefaultValues []string
	// Help is a message for the user on how to use this flag
	Help string

	// IsCommandFlag is set when parsing. Set to true if the flag was attached to a command (as opposed to being inherited from a section)
	IsCommandFlag bool
	// TypeDescription is set when parsing. Describes the type: int, string, ...
	TypeDescription string
	// SetBy might be set when parsing. Possible values: appdefault, config, passedflag
	SetBy string
	// Value might be set when parsing. The interface returned by updating a flag
	Value v.Value
}

// Resolve updates a flag's value from the command line, and then from the
// default value. flag should not be nil. deletes from flagStrs
func (flag *Flag) Resolve(
	name string,
	flagStrs map[string][]string,
	configReader configreader.ConfigReader,
) error {

	val, err := flag.EmptyValueConstructor()
	if err != nil {
		return fmt.Errorf("flag error: %v: %w", name, err)
	}
	flag.Value = val
	flag.TypeDescription = val.Description()

	// update from command line
	{
		strValues, exists := flagStrs[name]
		// the setby check for the first case is needed to
		// idempotently resolve flags (like the config flag for example)
		if flag.SetBy == "" && exists {

			if val.TypeInfo() == v.TypeInfoScalar && len(strValues) > 1 {
				return fmt.Errorf("flag error: %v: flag passed multiple times, it's value (type %v), can only be updated once", name, flag.TypeDescription)
			}

			for _, v := range strValues {
				flag.Value.Update(v)
			}
			flag.SetBy = "passedflag"
			// later we'll ensure that these aren't all used
			delete(flagStrs, name)
		}
	}

	// update from config
	{
		if flag.SetBy == "" && configReader != nil {
			fpr, err := configReader.Search(flag.ConfigPath)
			if err != nil {
				return err
			}
			if fpr.Exists {
				if !fpr.IsAggregated {
					err := flag.Value.ReplaceFromInterface(fpr.IFace)
					if err != nil {
						return err
					}
					flag.SetBy = "config"
				} else {
					under, ok := fpr.IFace.([]interface{})
					if !ok {
						return fmt.Errorf("expected []interface{}, got: %#v", under)
					}
					for _, e := range under {
						err = flag.Value.UpdateFromInterface(e)
						if err != nil {
							return fmt.Errorf("could not update container type value: %w", err)
						}
					}
					flag.SetBy = "config"
				}
			}
		}
	}

	// update from default
	{
		if flag.SetBy == "" && len(flag.DefaultValues) > 0 {
			for _, v := range flag.DefaultValues {
				flag.Value.Update(v)
			}
			flag.SetBy = "appdefault"
		}
	}

	return nil
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
