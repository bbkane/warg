package flag

import "go.bbkane.com/warg"

// Alias is an alternative name for a flag, usually shorter :)
func Alias(alias string) warg.FlagOpt {
	return func(f *warg.Flag) {
		f.Alias = alias
	}
}

// ConfigPath adds a configpath to a flag
func ConfigPath(path string) warg.FlagOpt {
	return func(flag *warg.Flag) {
		flag.ConfigPath = path
	}
}

// EnvVars adds a list of environmental variables to search through to update this flag. The first one that exists will be used to update the flag. Further existing envvars will be ignored.
func EnvVars(name ...string) warg.FlagOpt {
	return func(f *warg.Flag) {
		f.EnvVars = name
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
func UnsetSentinel(name string) warg.FlagOpt {
	return func(f *warg.Flag) {
		f.UnsetSentinel = name
	}
}

// Required means the user MUST fill this flag
func Required() warg.FlagOpt {
	return func(f *warg.Flag) {
		f.Required = true
	}
}
