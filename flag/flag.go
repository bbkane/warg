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
	// DefaultValue will be shoved into Value if needed
	// can be nil
	DefaultValue v.Value
	Help         string
	// SetBy holds where a flag is initialized. Is empty if not initialized
	SetBy string
	// Value holds what gets passed to the flag: --myflag value
	// and should be initialized to the empty value
	Value v.Value
}

func NewFlag(helpShort string, empty v.Value, opts ...FlagOpt) Flag {
	flag := Flag{
		Help:  helpShort,
		Value: empty,
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

func Default(value v.Value) FlagOpt {
	return func(flag *Flag) {
		flag.DefaultValue = value
	}
}
