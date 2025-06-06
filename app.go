package warg

import (
	"fmt"
	"runtime/debug"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

// AppOpt let's you customize the app. Most AppOpts panic if incorrectly called
type AppOpt func(*wargcore.App)

// GlobalFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func GlobalFlag(name string, value wargcore.Flag) AppOpt {
	return func(com *wargcore.App) {
		com.GlobalFlags.AddFlag(name, value)
	}
}

// GlobalFlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func GlobalFlagMap(flagMap wargcore.FlagMap) AppOpt {
	return func(com *wargcore.App) {
		com.GlobalFlags.AddFlags(flagMap)
	}
}

// NewGlobalFlag adds a flag to the app. It panics if a flag with the same name exists
func NewGlobalFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) AppOpt {
	return GlobalFlag(name, flag.New(helpShort, empty, opts...))

}

// ConfigFlag adds a flag that will be used to read a config file. If the passed flagMap is nil, DefaultConfigFlagMap will be used. The flag will be added to the app's global flags. When parsed, the config flag will be parsed before other flags, any config file found will be read, and any values found will be used to update other flags. This allows users to override flag defaults with values from a config file.
func ConfigFlag(reader config.NewReader, flagMap wargcore.FlagMap) AppOpt {
	return func(app *wargcore.App) {
		if len(flagMap) != 1 {
			panic(fmt.Sprintf("ConfigFlagMap must have exactly one flag, got %d", len(flagMap)))
		}
		app.NewConfigReader = reader
		app.ConfigFlagName = flagMap.SortedNames()[0]
		app.GlobalFlags.AddFlags(flagMap)
	}
}

// HelpFlag customizes your help flag. This option is only needed if you're also writing a custom help function. helpFlags be either `nil` to autogenerate or a flag map with one flat that with the followng properties:
//
//   - scalar string type
//   - choices that match the names in helpCommands
//   - default value set to one of the choices
//
// These properties are checked at runtime with app.Validate().
func HelpFlag(helpCommands wargcore.CommandMap, helpFlags wargcore.FlagMap) AppOpt {
	return func(a *wargcore.App) {
		switch len(helpFlags) {
		case 0:
			helpFlags = help.DefaultHelpFlagMap("default", helpCommands.SortedNames())
		case 1:
			break
		default:
			panic(fmt.Sprintf("helpFlags must have 0 or 1 flags, got %d", len(helpFlags)))
		}

		a.HelpFlagName = helpFlags.SortedNames()[0]
		a.HelpCommands = helpCommands
		a.GlobalFlags.AddFlags(helpFlags)
	}
}

// SkipValidation skips (most of) the app's internal consistency checks when the app is created.
// If used, make sure to call app.Validate() in a test!
func SkipValidation() AppOpt {
	return func(a *wargcore.App) {
		a.SkipValidation = true
	}
}

func debugBuildInfoVersion() string {
	// If installed via `go install`, we'll be able to read runtime version info
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	// when run with `go run`, this will return "(devel)"
	return info.Main.Version
}

// ColorFlagMap returns a map with a single "--color" flag that can be used to control color output.
//
// Example:
//
//	warg.GlobalFlagMap(warg.ColorFlagMap())
func ColorFlagMap() wargcore.FlagMap {
	return wargcore.FlagMap{
		"--color": flag.New(
			"Use ANSI colors",
			scalar.String(
				scalar.Choices("true", "false", "auto"),
				scalar.Default("auto"),
			),
			flag.EnvVars("WARG_COLOR"),
		),
	}
}

// VersioncommandMap returns a map with a single "version" command that prints the app version.
//
// Example:
//
//	warg.GlobalFlagMap(warg.ColorFlagMap())
func VersionCommandMap() wargcore.CommandMap {
	return wargcore.CommandMap{
		"version": command.New(
			"Print version",
			func(ctx wargcore.Context) error {
				fmt.Fprintln(ctx.Stdout, ctx.App.Version)
				return nil
			},
		),
	}
}

// New creates a warg app. name is used for help output only (though generally it should match the name of the compiled binary). version is the app version - if empty, warg will attempt to set it to the go module version, or "unknown" if that fails.
func New(name string, version string, rootSection wargcore.Section, opts ...AppOpt) wargcore.App {
	app := wargcore.App{
		Name:            name,
		RootSection:     rootSection,
		ConfigFlagName:  "",
		NewConfigReader: nil,
		HelpFlagName:    "",
		HelpCommands:    make(wargcore.CommandMap),
		SkipValidation:  false,
		Version:         version,
		GlobalFlags:     make(wargcore.FlagMap),
	}
	for _, opt := range opts {
		opt(&app)
	}

	if app.HelpFlagName == "" {
		HelpFlag(
			help.DefaultHelpCommandMap(),
			help.DefaultHelpFlagMap("default", help.DefaultHelpCommandMap().SortedNames()),
		)(&app)
	}

	if app.Version == "" {
		app.Version = debugBuildInfoVersion()
	}

	// validate or not and return
	if app.SkipValidation {
		return app
	}

	err := app.Validate()
	if err != nil {
		panic(err)
	}
	return app
}
