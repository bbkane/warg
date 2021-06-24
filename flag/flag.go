package flag

import (
	v "github.com/bbkane/warg/value"
)

type FlagMap = map[string]Flag
type FlagOpt = func(*Flag)

type Flag struct {
	// Default will be shoved into Value if needed
	// can be nil
	Default   v.Value
	HelpLong  string
	HelpShort string
	// SetBy holds where a flag is initialized. Is empty if not initialized
	SetBy string
	// Value holds what gets passed to the flag: --myflag value
	// and should be initialized to the empty value
	Value v.Value
}

func NewFlag(empty v.Value, opts ...FlagOpt) Flag {
	flag := Flag{}
	flag.Value = empty
	for _, opt := range opts {
		opt(&flag)
	}
	return flag
}

func WithDefault(value v.Value) FlagOpt {
	return func(flag *Flag) {
		flag.Default = value
	}
}

func WithFlagHelpLong(helpLong string) FlagOpt {
	return func(cat *Flag) {
		cat.HelpLong = helpLong
	}
}

func WithFlagHelpShort(helpShort string) FlagOpt {
	return func(cat *Flag) {
		cat.HelpShort = helpShort
	}
}
