package warg

import (
	"log"
	"sort"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/value"
)

// FlagOpt is a functional option for configuring a [Flag] during creation.
type FlagOpt func(*Flag)

// NewFlag creates a [Flag] with the given short help text, value constructor, and options.
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
	}
	for _, opt := range opts {
		opt(&flag)
	}
	return flag
}

// Alias sets a short alternative name for the flag (e.g., "-n" for "--name").
func Alias(alias string) FlagOpt {
	return func(f *Flag) {
		f.Alias = alias
	}
}

// ConfigPath sets the dot-separated path used to look up this flag's value in a config file
// (e.g., "database.host").
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

// FlagCompletions sets a custom [CompletionsFunc] for generating tab-completion candidates.
func FlagCompletions(CompletionsFunc CompletionsFunc) FlagOpt {
	return func(flag *Flag) {
		flag.Completions = CompletionsFunc
	}
}

// EnvVars sets environment variable names to read this flag's value from.
// The first existing variable wins; subsequent ones are ignored.
func EnvVars(name ...string) FlagOpt {
	return func(f *Flag) {
		f.EnvVars = name
	}
}

// Required marks the flag as mandatory. Parsing fails if a required flag is not set
// from any source (CLI, config, env var, or default).
func Required() FlagOpt {
	return func(f *Flag) {
		f.Required = true
	}
}

// UnsetSentinel defines a special value that, when passed on the command line, resets
// the flag to its empty state, undoing any prior prior value from arguments / config / env vars / defaults
// Conventionally set to "UNSET".
//
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

// FlagGroup sets the group name for organizing this flag in help output.
// Flags with the same group are displayed together. Empty string means ungrouped.
func FlagGroup(group string) FlagOpt {
	return func(f *Flag) {
		f.Group = group
	}
}

// FlagMap maps flag names (e.g., "--verbose") to [Flag] definitions.
type FlagMap map[string]Flag

// SortedNames returns the flag names in alphabetical order.
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

// FlagNameGroup pairs a group name with sorted flag names for grouped help output.
type FlagNameGroup struct {
	// Name is the group name. Empty string means ungrouped (displayed first).
	Name string
	// FlagNames are the sorted flag names in this group.
	FlagNames []string
}

// groupedNames returns flag names organized by Group, with ungrouped flags first,
// then groups in alphabetical order. Within each group, flags are sorted alphabetically.
func (fm *FlagMap) groupedNames() []FlagNameGroup {
	groups := make(map[string][]string)
	for name, flag := range *fm {
		groups[flag.Group] = append(groups[flag.Group], name)
	}
	for _, names := range groups {
		sort.Strings(names)
	}

	groupNames := make([]string, 0, len(groups))
	for g := range groups {
		if g != "" {
			groupNames = append(groupNames, g)
		}
	}
	sort.Strings(groupNames)

	var result []FlagNameGroup
	if names, ok := groups[""]; ok {
		result = append(result, FlagNameGroup{Name: "", FlagNames: names})
	}
	for _, g := range groupNames {
		result = append(result, FlagNameGroup{Name: g, FlagNames: groups[g]})
	}
	return result
}

// AddFlag inserts a flag into the map. Panics if a flag with the same name already exists.
func (fm FlagMap) AddFlag(name string, value Flag) {
	if _, alreadyThere := (fm)[name]; !alreadyThere {
		(fm)[name] = value
	} else {
		log.Panicf("flag already exists: %#v\n", name)
	}
}

// AddFlags merges another [FlagMap] into this one. Panics if any name already exists.
func (fm FlagMap) AddFlags(flagMap FlagMap) {
	for name, value := range flagMap {
		fm.AddFlag(name, value)
	}
}

// Flag defines a single command-line flag, including its type, help text, aliases,
// config/env bindings, and validation rules.
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

	// Group is an optional group name for organizing flags in help output.
	// Flags with an empty group are printed first, then groups are printed in alphabetical order.
	Group string

	// HelpShort is a message for the user on how to use this flag
	HelpShort string

	// Required means the user MUST fill this flag
	Required bool

	// When UnsetSentinal is passed as a flag value, Value is reset and SetBy is set to ""
	UnsetSentinel *string
}
