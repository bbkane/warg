package flag

import (
	v "github.com/bbkane/warg/value"
)

type FlagMap = map[string]Flag
type FlagOpt = func(*Flag)

type Flag struct {

	// TODO: make these private. resolveFlag should probably be a method on flag
	ConfigFromInterface v.FromInterface
	ConfigPath          string
	// DefaultValues will be shoved into Value if the app builder specifies it.
	// For scalar values, the last DefaultValues wins
	DefaultValues []string
	Help          string
	// SetBy holds where a flag is initialized. Is empty if not initialized
	SetBy string
	// Value holds what gets passed to the flag: --myflag value
	// and should be initialized to the empty value
	// TODO: make this private? TODO: Update docs once this is successfully
	// an output instead of an input
	Value v.Value

	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor v.EmptyConstructor
}

func NewFlag(helpShort string, empty v.EmptyConstructor, opts ...FlagOpt) Flag {
	flag := Flag{
		Help:                  helpShort,
		EmptyValueConstructor: empty,
	}
	for _, opt := range opts {
		opt(&flag)
	}
	return flag
}

func ConfigPath(path string, valueFromInterface v.FromInterface) FlagOpt {
	return func(flag *Flag) {
		flag.ConfigPath = path
		flag.ConfigFromInterface = valueFromInterface
	}
}

func Default(values ...string) FlagOpt {
	return func(flag *Flag) {
		flag.DefaultValues = values
	}
}
