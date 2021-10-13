package flag

import (
	"fmt"

	"github.com/bbkane/warg/configreader"
	v "github.com/bbkane/warg/value"
)

type FlagMap = map[string]Flag
type FlagOpt = func(*Flag)

type FlagValues = map[string]interface{}

type Flag struct {
	ConfigPath string
	// DefaultValues will be shoved into Value if the app builder specifies it.
	// For scalar values, the last DefaultValues wins
	DefaultValues []string
	Help          string
	// SetBy possible values: appdefault, config, passedflag
	SetBy string
	// Value holds what gets passed to the flag: --myflag value
	// and should be initialized to the empty value
	// TODO: make this private? TODO: Update docs once this is successfully
	// an output instead of an input
	Value v.Value

	// EmptyConstructor tells flag how to make a value
	EmptyValueConstructor v.EmptyConstructor
}

// resolveFLag updates a flag's value from the command line, and then from the
// default value. flag should not be nil. deletes from flagStrs
func (flag *Flag) Resolve(name string, flagStrs map[string][]string, configReader configreader.ConfigReader) error {

	v, err := flag.EmptyValueConstructor()
	if err != nil {
		return fmt.Errorf("flag error: %v: %w", name, err)
	}
	flag.Value = v

	// update from command line
	{
		strValues, exists := flagStrs[name]
		// the setby check for the first case is needed to
		// idempotently resolve flags (like the config flag for example)
		if flag.SetBy == "" && exists {
			for _, v := range strValues {
				// TODO: make sure we don't update over flags meant to be set once
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
				// TODO: don't update flags more than once if they're not supposed to be
				flag.Value.Update(v)
			}
			flag.SetBy = "appdefault"
		}
	}

	return nil
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

func ConfigPath(path string) FlagOpt {
	return func(flag *Flag) {
		flag.ConfigPath = path
	}
}

func Default(values ...string) FlagOpt {
	return func(flag *Flag) {
		flag.DefaultValues = values
	}
}
