package flag

import (
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/value"
)

// FlagOpt customizes a Flag on creation
type FlagOpt func(*cli.Flag)

// New creates a New with options!
func New(helpShort string, empty value.EmptyConstructor, opts ...FlagOpt) cli.Flag {
	flag := cli.Flag{
		Alias:                 "",
		CompletionCandidates:  nil,
		ConfigPath:            "",
		EmptyValueConstructor: empty,
		EnvVars:               nil,
		HelpShort:             helpShort,
		Required:              false,
		UnsetSentinel:         "",
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
	return func(f *cli.Flag) {
		f.Alias = alias
	}
}

// ConfigPath adds a configpath to a flag
func ConfigPath(path string) FlagOpt {
	return func(f *cli.Flag) {
		f.ConfigPath = path
	}
}

func CompletionCandidate(completionCandidatesFunc func(cli.Context) (*completion.Candidates, error)) FlagOpt {
	return func(flag *cli.Flag) {
		flag.CompletionCandidates = completionCandidatesFunc
	}
}

// EnvVars adds a list of environmental variables to search through to update this flag. The first one that exists will be used to update the flag. Further existing envvars will be ignored.
func EnvVars(name ...string) FlagOpt {
	return func(f *cli.Flag) {
		f.EnvVars = name
	}
}

// Required means the user MUST fill this flag
func Required() FlagOpt {
	return func(f *cli.Flag) {
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
	return func(f *cli.Flag) {
		f.UnsetSentinel = name
	}
}
