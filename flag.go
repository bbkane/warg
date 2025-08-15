package warg

import (
	"log"
	"sort"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/value"
)

// FlagOpt customizes a Flag on creation
type FlagOpt func(*Flag)

// NewFlag creates a NewFlag with options!
func NewFlag(helpShort string, empty value.EmptyConstructor, opts ...FlagOpt) Flag {
	flag := Flag{
		Alias:                 "",
		Completions:           defaultFlagCompletions,
		ConfigPath:            "",
		EmptyValueConstructor: empty,
		EnvVars:               nil,
		HelpShort:             helpShort,
		Required:              false,
		UnsetSentinel:         nil,
		// Deprecated
		IsCommandFlag: false,
		Value:         nil,
	}
	for _, opt := range opts {
		opt(&flag)
	}
	return flag
}

// Alias is an alternative name for a flag, usually shorter :)
func Alias(alias string) FlagOpt {
	return func(f *Flag) {
		f.Alias = alias
	}
}

// ConfigPath adds a configpath to a flag
func ConfigPath(path string) FlagOpt {
	return func(f *Flag) {
		f.ConfigPath = path
	}
}

func defaultFlagCompletions(cmdCtx CmdContext) (*completion.Candidates, error) {
	choices := cmdCtx.ParseState.FlagValues[cmdCtx.ParseState.CurrentFlagName].Choices()
	if len(choices) > 0 {
		candidates := &completion.Candidates{
			Type:   completion.Type_Values,
			Values: []completion.Candidate{},
		}
		// pr.FlagValues is always filled with at least the empty values
		for _, name := range choices {
			candidates.Values = append(candidates.Values, completion.Candidate{
				Name:        name,
				Description: "",
			})
		}
		return candidates, nil
	}

	// special case: bools can only be true or false, so let's be helpful and suggest those
	if _, ok := cmdCtx.ParseState.FlagValues[cmdCtx.ParseState.CurrentFlagName].Get().(bool); ok {
		return &completion.Candidates{
			Type: completion.Type_Values,
			Values: []completion.Candidate{
				{Name: "true", Description: ""},
				{Name: "false", Description: ""},
			},
		}, nil
	}

	// default
	return &completion.Candidates{
		Type:   completion.Type_DirectoriesFiles,
		Values: nil,
	}, nil

}

func FlagCompletions(CompletionsFunc CompletionsFunc) FlagOpt {
	return func(flag *Flag) {
		flag.Completions = CompletionsFunc
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
		f.UnsetSentinel = &name
	}
}

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

	// Completions is a function that returns a list of completion candidates for this flag.
	// Note that some flags in the cli.Context Flags map may not be set, even if they're required.
	// TODO: get a comprehensive list of restrictions on the context.
	Completions CompletionsFunc

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
	UnsetSentinel *string

	// -- the following are set when parsing and they're all deprecated for my march to immutabality

	// IsCommandFlag is set when parsing. Set to true if the flag was attached to a command (as opposed to being inherited from a section)
	// Deprecated: Check if the name is in command.Flags instead
	IsCommandFlag bool

	// Value is set when parsing. Use SetBy != "" to determine whether a value was actually passed  instead of being empty
	// Deprecated: check the value from ParseState.FlagValues
	Value value.Value
}
