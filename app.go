package warg

import (
	"fmt"
	"runtime/debug"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/section"
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
func HelpFlag(helpCommands wargcore.CmdMap, helpFlags wargcore.FlagMap) AppOpt {
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

// SkipAll skips adding:
//   - the default completion commands (<app> completion)
//   - the default color flag map (<app> --color)
//   - the default version command map (<app> version)
//   - the default validation checks
//
// This is inteded for tests where you just want to assert against a minimal application
func SkipAll() AppOpt {
	return func(a *wargcore.App) {
		a.SkipCompletionCommands = true
		a.SkipGlobalColorFlag = true
		a.SkipVersionCommand = true
		a.SkipValidation = true
	}
}

// SkipCompletionCommands skips adding the default completion commands (<app> completion).
func SkipCompletionCommands() AppOpt {
	return func(a *wargcore.App) {
		a.SkipCompletionCommands = true
	}
}

// SkipColorFlag skips adding the default color flag map (<app> --color).
func SkipGlobalColorFlag() AppOpt {
	return func(a *wargcore.App) {
		a.SkipGlobalColorFlag = true
	}
}

// SkipValidation skips (most of) the app's internal consistency checks when the app is created.
// If used, make sure to call app.Validate() in a test!
func SkipValidation() AppOpt {
	return func(a *wargcore.App) {
		a.SkipValidation = true
	}
}

// SkipVersionCommand skips adding the default version command (<app> version).
func SkipVersionCommand() AppOpt {
	return func(a *wargcore.App) {
		a.SkipVersionCommand = true
	}
}

// FindVersion returns the version of the app. If the version is already set (eg. via a build flag), it returns that. Otherwise, it tries to read the go module version from the runtime info, or returns "unknown" if that fails.
func FindVersion(version string) string {
	// if the version is already set (eg. via a build flag), return it
	if version != "" {
		return version
	}

	// If installed via `go install`, we'll be able to read runtime version info
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	// when run with `go run`, this will return "(devel)"
	return info.Main.Version
}

func CompletionsDirectories(ctx wargcore.Context) (*completion.Candidates, error) {
	return &completion.Candidates{
		Type:   completion.Type_Directories,
		Values: nil,
	}, nil
}

func CompletionsDirectoriesFiles(ctx wargcore.Context) (*completion.Candidates, error) {
	return &completion.Candidates{
		Type:   completion.Type_DirectoriesFiles,
		Values: nil,
	}, nil
}

func CompletionsNone(ctx wargcore.Context) (*completion.Candidates, error) {
	return &completion.Candidates{
		Type:   completion.Type_None,
		Values: nil,
	}, nil
}

func CompletionsValues(values []string) wargcore.CompletionsFunc {
	var vals []completion.Candidate
	for _, v := range values {
		vals = append(vals, completion.Candidate{Name: v, Description: ""})
	}
	return func(ctx wargcore.Context) (*completion.Candidates, error) {

		return &completion.Candidates{
			Type:   completion.Type_Values,
			Values: vals,
		}, nil
	}
}

func CompletionsValuesDescriptions(values []completion.Candidate) wargcore.CompletionsFunc {
	return func(ctx wargcore.Context) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_ValuesDescriptions,
			Values: values,
		}, nil
	}
}

// New creates a warg app. name is used for help output only (though generally it should match the name of the compiled binary). version is the app version - if empty, warg will attempt to set it to the go module version, or "unknown" if that fails.
func New(name string, version string, rootSection wargcore.Section, opts ...AppOpt) wargcore.App {
	app := wargcore.App{
		Name:                   name,
		RootSection:            rootSection,
		ConfigFlagName:         "",
		NewConfigReader:        nil,
		HelpFlagName:           "",
		HelpCommands:           make(wargcore.CmdMap),
		SkipCompletionCommands: false,
		SkipValidation:         false,
		SkipGlobalColorFlag:    false,
		SkipVersionCommand:     false,
		Version:                version,
		GlobalFlags:            make(wargcore.FlagMap),
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

	app.Version = FindVersion(app.Version)

	if !app.SkipGlobalColorFlag {
		GlobalFlagMap(wargcore.FlagMap{
			"--color": flag.New(
				"Use ANSI colors",
				scalar.String(
					scalar.Choices("true", "false", "auto"),
					scalar.Default("auto"),
				),
				flag.EnvVars("WARG_COLOR"),
			),
		})(&app)
	}

	if !app.SkipCompletionCommands {
		section.NewChildSection(
			"completion",
			"Print shell completion scripts",
			section.NewChildCmd(
				"zsh",
				"Print zsh completion script",
				func(ctx wargcore.Context) error {
					completion.ZshCompletionScriptWrite(ctx.Stdout, app.Name)
					return nil
				},
			),
		)(&app.RootSection)
	}

	if !app.SkipVersionCommand {
		section.NewChildCmd(
			"version",
			"Print version",
			func(ctx wargcore.Context) error {
				fmt.Fprintln(ctx.Stdout, ctx.App.Version)
				return nil
			},
		)(&app.RootSection)
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
