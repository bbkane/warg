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
	"go.bbkane.com/warg/section"
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

	// New Help()
	name         string
	helpFlagName flag.Name
	// Note that this can be ""
	helpFlagAlias flag.Name
	helpMappings  []help.HelpFlagMapping

	// HelpFile contains the file set in OverrideHelpFlag.
	// HelpFile is part of the public API to allow for easier testing.
	// HelpFile is never closed by warg, so if setting it to something other than stderr/stdout,
	// please remember to close HelpFile after using ParseResult.Action (which writes to HelpFile).
	HelpFile *os.File

	// rootSection holds the good stuff!
	rootSection section.SectionT

	skipValidation bool
}

// OverrideHelpFlag customizes your --help. If you write a custom --help function, you'll want to add it to your app here!
func OverrideHelpFlag(
	mappings []help.HelpFlagMapping,
	helpFile *os.File,
	flagName flag.Name,
	flagHelp flag.HelpShort,
	flagOpts ...flag.FlagOpt,
) AppOpt {
	return func(a *App) {

		if !strings.HasPrefix(string(flagName), "-") {
			log.Panicf("flagName should start with '-': %#v\n", flagName)
		}

		if _, alreadyThere := a.rootSection.Flags[flagName]; alreadyThere {
			log.Panicf("flag already exists: %#v\n", flagName)
		}
		helpValues := make([]string, len(mappings))
		for i := range mappings {
			helpValues[i] = mappings[i].Name
		}

		helpFlag := flag.New(
			flagHelp,
			scalar.String(
				scalar.Choices(helpValues...),
			),
			flagOpts...,
		)

		if len(helpFlag.DefaultValues) == 0 {
			log.Panic("--help flag must have a default. use flag.Default(...) to set one")
		}

		a.rootSection.Flags[flagName] = helpFlag
		// This is used in parsing, so no need to strongly type it
		a.helpFlagName = flagName
		a.helpFlagAlias = helpFlag.Alias
		a.helpMappings = mappings
		a.HelpFile = helpFile

	}
}

// Use ConfigFlag in conjunction with flag.ConfigPath to allow users to override flag defaults with values from a config.
func ConfigFlag(
	// TODO: put the new stuff at the front to be consistent with OverrideHelpFlag
	configFlagName flag.Name,
	newConfigReader config.NewReader,
	helpShort flag.HelpShort,
	flagOpts ...flag.FlagOpt,
) AppOpt {
	return func(app *App) {
		app.configFlagName = configFlagName
		app.newConfigReader = newConfigReader
		// TODO: need to have value opts here
		configFlag := flag.New(helpShort, scalar.Path(), flagOpts...)
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

// AddVersionCommand adds a "version" command to the root section.
// Pass an empty string to use .Main.Version from `debug.ReadBuildInfo`,
// which returns "(devel)" when using `go run`
func AddVersionCommand(version string) AppOpt {
	return func(a *App) {
		action := func(ctx command.Context) error {
			if version != "" {
				fmt.Println(version)
				return nil
			}
			// If installed via `go install`, we'll be able to read runtime version info
			info, ok := debug.ReadBuildInfo()
			if !ok {
				return errors.New("unable to read build info")
			}
			// when run with `go run`, this will return "(devel)"
			fmt.Println(info.Main.Version)
			return nil
		}
		section.Command("version", "Print version", action)(&a.rootSection)
	}
}

// New builds a new App!
func New(name string, rootSection section.SectionT, opts ...AppOpt) App {
	app := App{
		name:        name,
		rootSection: rootSection,
	}
	for _, opt := range opts {
		opt(&app)
	}

	if app.helpFlagName == "" {
		OverrideHelpFlag(
			help.BuiltinHelpFlagMappings(),
			os.Stdout,
			"--help",
			"Print help",
			flag.Alias("-h"),
			flag.Default("default"),
		)(&app)
	}

	if !app.skipValidation {
		err := app.Validate()
		if err != nil {
			panic(err)
		}
	}

	return app
}

// MustRun runs the app.
// Any flag parsing errors will be printed to stderr and os.Exit(64) (EX_USAGE) will be called.
// Any errors on an Action will be printed to stderr and os.Exit(1) will be called.
func (app *App) MustRun(osArgs []string, osLookupEnv LookupFunc) {
	pr, err := app.Parse(osArgs, osLookupEnv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// https://unix.stackexchange.com/a/254747/185953
		os.Exit(64)
	}
	err = pr.Action(pr.Context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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

type flagNameSet map[flag.Name]struct{}

// addFlags adds a flag's name and alias to the set. Returns an error
// if the name OR alias already exists
func (fs flagNameSet) addFlags(fm flag.FlagMap) error {
	for flagName := range fm {
		_, exists := fs[flagName]
		if exists {
			return fmt.Errorf("flag or alias name exists twice: %v", flagName)
		}
		fs[flagName] = struct{}{}

		alias := fm[flagName].Alias
		if alias != "" {
			_, exists := fs[alias]
			if exists {
				return fmt.Errorf("flag or alias name exists twice: %v", alias)
			}
			fs[alias] = struct{}{}
		}
	}
	return nil
}

func validateFlags(
	inheritedFlags flag.FlagMap,
	secFlags flag.FlagMap,
	comFlags flag.FlagMap,
) error {
	nameSet := make(flagNameSet)
	var err error

	err = nameSet.addFlags(inheritedFlags)
	if err != nil {
		return err
	}

	err = nameSet.addFlags(secFlags)
	if err != nil {
		return err
	}

	err = nameSet.addFlags(comFlags)
	if err != nil {
		return err
	}

	// fmt.Printf("%#v\n", nameSet)

	for name := range nameSet {
		if !strings.HasPrefix(string(name), "-") {
			return fmt.Errorf("flag and alias names must start with '-': %#v", name)
		}
	}

	return nil

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

		// No need to check section flags here; all sections will end in a command
		// and we can check there

		for name, com := range flatSec.Sec.Commands {

			// Commands must not start wtih "-"
			if strings.HasPrefix(string(name), "-") {
				return fmt.Errorf("command names must not start with '-': %#v", name)
			}

			err := validateFlags(flatSec.InheritedFlags, flatSec.Sec.Flags, com.Flags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
