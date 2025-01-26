// Declaratively create heirarchical command line apps.
package warg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/scalar"
)

// AppOpt let's you customize the app. Most AppOpts panic if incorrectly called
type AppOpt func(*App)

// An App contains your defined sections, commands, and flags
// Create a new App with New()
type App struct {
	// Config()
	configFlagName  flag.Name
	newConfigReader config.NewReader
	configFlag      *flag.Flag

	globalFlags flag.FlagMap

	// New Help()
	name         string
	helpFlagName flag.Name
	// Note that this can be ""
	helpFlagAlias flag.Name
	helpMappings  []help.HelpFlagMapping

	// rootSection holds the good stuff!
	rootSection section.SectionT

	skipValidation bool

	version string
}

// OverrideHelpFlag customizes your --help. If you write a custom --help function, you'll want to add it to your app here!
func OverrideHelpFlag(
	mappings []help.HelpFlagMapping,
	defaultChoice string,
	flagName flag.Name,
	flagHelp flag.HelpShort,
	flagOpts ...flag.FlagOpt,
) AppOpt {
	return func(a *App) {

		if !strings.HasPrefix(string(flagName), "-") {
			log.Panicf("flagName should start with '-': %#v\n", flagName)
		}

		if _, alreadyThere := a.globalFlags[flagName]; alreadyThere {
			log.Panicf("flag already exists: %#v\n", flagName)
		}

		defaultFound := false
		helpValues := make([]string, len(mappings))
		for i := range mappings {
			helpValues[i] = mappings[i].Name
			if helpValues[i] == defaultChoice {
				defaultFound = true
			}
		}

		if !defaultFound {
			panic(fmt.Sprintf("default (%#v) not found in helpValues (%#v)", defaultChoice, helpValues))
		}

		helpFlag := flag.New(
			flagHelp,
			scalar.String(
				scalar.Choices(helpValues...),
				scalar.Default(defaultChoice),
			),
			flagOpts...,
		)

		a.globalFlags[flagName] = helpFlag
		// This is used in parsing, so no need to strongly type it
		a.helpFlagName = flagName
		a.helpFlagAlias = helpFlag.Alias
		a.helpMappings = mappings

	}
}

// ExistingGlobalFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func ExistingGlobalFlag(name flag.Name, value flag.Flag) AppOpt {
	return func(com *App) {
		com.globalFlags.AddFlag(name, value)
	}
}

// GlobalFlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func GlobalFlagMap(flagMap flag.FlagMap) AppOpt {
	return func(com *App) {
		com.globalFlags.AddFlags(flagMap)
	}
}

// GlobalFlag adds a flag to the app. It panics if a flag with the same name exists
func GlobalFlag(name flag.Name, helpShort flag.HelpShort, empty value.EmptyConstructor, opts ...flag.FlagOpt) AppOpt {
	return ExistingGlobalFlag(name, flag.New(helpShort, empty, opts...))

}

// Use ConfigFlag in conjunction with flag.ConfigPath to allow users to override flag defaults with values from a config.
// This flag will be parsed and any resulting config will be read before other flag value sources.
func ConfigFlag(
	// TODO: put the new stuff at the front to be consistent with OverrideHelpFlag
	configFlagName flag.Name,
	// TODO: can I make this nicer?
	scalarOpts []scalar.ScalarOpt[path.Path],
	newConfigReader config.NewReader,
	helpShort flag.HelpShort,
	flagOpts ...flag.FlagOpt,
) AppOpt {
	return func(app *App) {
		app.configFlagName = configFlagName
		app.newConfigReader = newConfigReader
		// TODO: need to have value opts here
		configFlag := flag.New(helpShort, scalar.Path(scalarOpts...), flagOpts...)
		app.configFlag = &configFlag
	}
}

// SkipValidation skips (most of) the app's internal consistency checks when the app is created.
// If used, make sure to call app.Validate() in a test!
func SkipValidation() AppOpt {
	return func(a *App) {
		a.skipValidation = true
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
func ColorFlagMap() flag.FlagMap {
	return flag.FlagMap{
		"--color": flag.New(
			"Use ANSI colors",
			scalar.String(
				scalar.Choices("true", "false", "auto"),
				scalar.Default("auto"),
			),
		),
	}
}

// VersioncommandMap returns a map with a single "version" command that prints the app version.
//
// Example:
//
//	warg.GlobalFlagMap(warg.ColorFlagMap())
func VersionCommandMap() command.CommandMap {
	return command.CommandMap{
		"version": command.New(
			"Print version",
			func(ctx command.Context) error {
				fmt.Fprintln(ctx.Stdout, ctx.Version)
				return nil
			},
		),
	}
}

// New creates a warg app. name is used for help output only (though generally it should match the name of the compiled binary). version is the app version - if empty, warg will attempt to set it to the go module version, or "unknown" if that fails.
func New(name string, version string, rootSection section.SectionT, opts ...AppOpt) App {
	app := App{
		name:            name,
		rootSection:     rootSection,
		configFlagName:  "",
		newConfigReader: nil,
		configFlag:      nil,
		helpFlagName:    "",
		helpFlagAlias:   "",
		helpMappings:    nil,
		skipValidation:  false,
		version:         version,
		globalFlags:     make(flag.FlagMap),
	}
	for _, opt := range opts {
		opt(&app)
	}

	if app.helpFlagName == "" {
		OverrideHelpFlag(
			help.BuiltinHelpFlagMappings(),
			"default",
			"--help",
			"Print help",
			flag.Alias("-h"),
		)(&app)
	}

	if app.version == "" {
		app.version = debugBuildInfoVersion()
	}

	// validate or not and return
	if app.skipValidation {
		return app
	}

	err := app.Validate()
	if err != nil {
		panic(err)
	}
	return app
}

// MustRun runs the app.
// Any flag parsing errors will be printed to stderr and os.Exit(64) (EX_USAGE) will be called.
// Any errors on an Action will be printed to stderr and os.Exit(1) will be called.
func (app *App) MustRun(opts ...ParseOpt) {
	pr, err := app.Parse(opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// https://unix.stackexchange.com/a/254747/185953
		os.Exit(64)
	}
	err = pr.Action(pr.Context)
	if err != nil {
		fmt.Fprintln(pr.Context.Stderr, err)
		os.Exit(1)
	}
}

// Look up keys (meant for environment variable parsing) - fulfillable with os.LookupEnv or warg.LookupMap(map)
type LookupFunc func(key string) (string, bool)

// LookupMap loooks up keys from a provided map. Useful to mock os.LookupEnv when parsing
func LookupMap(m map[string]string) LookupFunc {
	return func(key string) (string, bool) {
		val, exists := m[key]
		return val, exists
	}
}

// validateFlags2 checks that global and command flag names and aliases start with "-" and are unique.
// It does not need to check the following scenarios:
//
//   - global flag names don't collide with global flag names (app will panic when adding the second global flag) - TOOD: ensure there's a test for this
//   - command flag names in the same command don't collide with each other (app will panic when adding the second command flag) TODO: ensure there's a test for this
//   - command flag names/aliases don't collide with command flag names/aliases in other commands (since only one command will be run, this is not a problem)
func validateFlags2(
	globalFlags flag.FlagMap,
	comFlags flag.FlagMap,
) error {
	nameCount := make(map[flag.Name]int)
	for name, fl := range globalFlags {
		nameCount[name]++
		if fl.Alias != "" {
			nameCount[fl.Alias]++
		}
	}
	for name, fl := range comFlags {
		nameCount[name]++
		if fl.Alias != "" {
			nameCount[fl.Alias]++
		}
	}
	var errs []error
	for name, count := range nameCount {
		if !strings.HasPrefix(string(name), "-") {
			errs = append(errs, fmt.Errorf("flag and alias names must start with '-': %#v", name))
		}
		if count > 1 {
			errs = append(errs, fmt.Errorf("flag or alias name exists %d times: %v", count, name))
		}
	}
	return errors.Join(errs...)
}

// Validate checks app for creation errors. It checks:
//
// - Sections and commands don't start with "-" (needed for parsing)
//
// - Flag names and aliases do start with "-" (needed for parsing)
//
// - Flag names and aliases don't collide
func (app *App) Validate() error {
	// NOTE: we need to be able to validate before we parse, and we may not know the app name
	// till after prsing so set the root path to "root"
	rootPath := []section.Name{section.Name(app.name)}
	it := app.rootSection.BreadthFirst(rootPath)

	for it.HasNext() {
		flatSec := it.Next()

		// Sections don't start with "-"
		secName := flatSec.Path[len(flatSec.Path)-1]
		if strings.HasPrefix(string(secName), "-") {
			return fmt.Errorf("section names must not start with '-': %#v", secName)
		}

		// Sections must not be leaf nodes
		if flatSec.Sec.Sections.Empty() && flatSec.Sec.Commands.Empty() {
			return fmt.Errorf("sections must have either child sections or child commands: %#v", secName)
		}

		{
			// child section names should not clash with child command names
			nameCount := make(map[string]int)
			for name := range flatSec.Sec.Commands {
				nameCount[string(name)]++
			}
			for name := range flatSec.Sec.Sections {
				nameCount[string(name)]++
			}
			errs := []error{}
			for name, count := range nameCount {
				if count > 1 {
					errs = append(errs, fmt.Errorf("command and section name clash: %s", name))
				}
			}
			err := errors.Join(errs...)
			if err != nil {
				return fmt.Errorf("name collision: %w", err)
			}
		}

		for name, com := range flatSec.Sec.Commands {

			// Commands must not start wtih "-"
			if strings.HasPrefix(string(name), "-") {
				return fmt.Errorf("command names must not start with '-': %#v", name)
			}

			err := validateFlags2(app.globalFlags, com.Flags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
