package warg

import (
	"fmt"
	"runtime/debug"

	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/scalar"
)

// AppOpt let's you customize the app. Most AppOpts panic if incorrectly called
type AppOpt func(*cli.App)

// OverrideHelpFlag customizes your help flag. helpFlagName should point to a previously added global flag with the following properties:
//
//   - scalar string type
//   - choices that match the names in helpCommands
//   - default value set to one of the choices
//
// These properties are checked at runtime with app.Validate()
func OverrideHelpFlag(helpFlagName string, helpCommands cli.CommandMap) AppOpt {
	return func(a *cli.App) {
		a.HelpFlagName = helpFlagName
		a.HelpCommands = helpCommands
	}
}

// GlobalFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func GlobalFlag(name string, value cli.Flag) AppOpt {
	return func(com *cli.App) {
		com.GlobalFlags.AddFlag(name, value)
	}
}

// GlobalFlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func GlobalFlagMap(flagMap cli.FlagMap) AppOpt {
	return func(com *cli.App) {
		com.GlobalFlags.AddFlags(flagMap)
	}
}

// NewGlobalFlag adds a flag to the app. It panics if a flag with the same name exists
func NewGlobalFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) AppOpt {
	return GlobalFlag(name, flag.NewFlag(helpShort, empty, opts...))

}

// Use ConfigFlag in conjunction with flag.ConfigPath to allow users to override flag defaults with values from a config.
// This flag will be parsed and any resulting config will be read before other flag value sources.
func ConfigFlag(
	// TODO: put the new stuff at the front to be consistent with OverrideHelpFlag
	configFlagName string,
	// TODO: can I make this nicer?
	scalarOpts []scalar.ScalarOpt[path.Path],
	newConfigReader config.NewReader,
	helpShort string,
	flagOpts ...flag.FlagOpt,
) AppOpt {
	return func(app *cli.App) {
		app.ConfigFlagName = configFlagName
		app.NewConfigReader = newConfigReader
		// TODO: need to have value opts here
		configFlag := flag.NewFlag(helpShort, scalar.Path(scalarOpts...), flagOpts...)
		app.ConfigFlag = &configFlag
	}
}

// SkipValidation skips (most of) the app's internal consistency checks when the app is created.
// If used, make sure to call app.Validate() in a test!
func SkipValidation() AppOpt {
	return func(a *cli.App) {
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
func ColorFlagMap() cli.FlagMap {
	return cli.FlagMap{
		"--color": flag.NewFlag(
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
func VersionCommandMap() cli.CommandMap {
	return cli.CommandMap{
		"version": command.NewCommand(
			"Print version",
			func(ctx cli.Context) error {
				fmt.Fprintln(ctx.Stdout, ctx.App.Version)
				return nil
			},
		),
	}
}

// NewApp creates a warg app. name is used for help output only (though generally it should match the name of the compiled binary). version is the app version - if empty, warg will attempt to set it to the go module version, or "unknown" if that fails.
func NewApp(name string, version string, rootSection cli.SectionT, opts ...AppOpt) cli.App {
	app := cli.App{
		Name:            name,
		RootSection:     rootSection,
		ConfigFlagName:  "",
		NewConfigReader: nil,
		ConfigFlag:      nil,
		HelpFlagName:    "",
		HelpCommands:    make(cli.CommandMap),
		SkipValidation:  false,
		Version:         version,
		GlobalFlags:     make(cli.FlagMap),
	}
	for _, opt := range opts {
		opt(&app)
	}

	if app.HelpFlagName == "" {
		app.GlobalFlags.AddFlags(help.DefaultHelpFlagMap("default", help.DefaultHelpCommandMap().SortedNames()))
		OverrideHelpFlag(
			"--help",
			help.DefaultHelpCommandMap(),
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
