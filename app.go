// Declaratively create heirarchical command line apps.
package warg

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	c "github.com/bbkane/warg/command"
	"github.com/bbkane/warg/configreader"
	f "github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/help"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

// AppOpt let's you customize the app. It panics if there is an error
type AppOpt = func(*App)

// An App contains your defined sections, commands, and flags
// Create a new App with New()
type App struct {
	// Config()
	configFlagName  string
	newConfigReader configreader.NewConfigReader
	configFlag      *f.Flag
	// Help()
	name          string
	helpFlagNames []string
	helpWriter    io.Writer
	sectionHelp   help.SectionHelp
	commandHelp   help.CommandHelp
	// rootSection holds the good stuff!
	rootSection s.Section
}

// OverrideHelp will let you provide own help function.
func OverrideHelp(w io.Writer, helpFlagNames []string, sectionHelp help.SectionHelp, commandHelp help.CommandHelp) AppOpt {
	return func(app *App) {
		app.sectionHelp = sectionHelp
		app.commandHelp = commandHelp
		app.helpFlagNames = helpFlagNames
		app.helpWriter = w
		for _, n := range helpFlagNames {
			if !strings.HasPrefix(n, "-") {
				log.Panicf("helpFlags should start with '-': %#v\n", n)
			}
		}
	}
}

// ConfigFlag lets you customize your config flag. Especially useful for changing the config reader (for example to choose whether to use a JSON or YAML structured config)
func ConfigFlag(
	configFlagName string,
	newConfigReader configreader.NewConfigReader,
	helpShort string,
	flagOpts ...f.FlagOpt,
) AppOpt {
	return func(app *App) {
		app.configFlagName = configFlagName
		app.newConfigReader = newConfigReader
		configFlag := f.New(helpShort, v.Path, flagOpts...)
		app.configFlag = &configFlag
	}
}

// New builds a new App!
func New(name string, rootSection s.Section, opts ...AppOpt) App {
	app := App{
		name:        name,
		rootSection: rootSection,
	}
	for _, opt := range opts {
		opt(&app)
	}

	// Help
	if len(app.helpFlagNames) == 0 {
		OverrideHelp(
			os.Stderr,
			[]string{"-h", "--help"},
			help.DefaultSectionHelp,
			help.DefaultCommandHelp,
		)(&app)
	}

	return app
}

type gatherArgsResult struct {
	// Appname holds os.Args[0]
	AppName string
	// Path holds the path to the current command/section
	Path []string
	// FlagStrings is a map of all flags to their values
	FlagStrs   map[string][]string
	HelpPassed bool
}

func containsString(haystack []string, needle string) bool {
	for _, w := range haystack {
		if w == needle {
			return true
		}
	}
	return false
}

// gatherArgs "parses" os.Argv into commands and flags. It's a 'lowering' function,
// simplifying os.Args as much as possible before needing knowledge of this particular app
// --help does NOT require a value
func gatherArgs(osArgs []string, helpFlagNames []string) (*gatherArgsResult, error) {
	res := &gatherArgsResult{
		FlagStrs: make(map[string][]string),
	}
	res.AppName = osArgs[0]

	// let's declare some states with an "enum"...
	expectingAnything := "expectingAnything"
	expectingFlagValue := "expectingFlagValue"
	// currentFlagName is only valid when expectingFlagValue
	// I miss ADTs in go
	var currentFlagName string

	// set up initial conditions
	currentFlagName = ""
	expecting := expectingAnything
	for _, word := range osArgs[1:] {
		switch expecting {
		case expectingAnything:
			if containsString(helpFlagNames, word) {
				res.HelpPassed = true
				continue
			}
			if strings.HasPrefix(word, "-") {
				currentFlagName = word
				expecting = expectingFlagValue
			} else {
				// command case
				res.Path = append(res.Path, word)
			}
		case expectingFlagValue:
			res.FlagStrs[currentFlagName] = append(res.FlagStrs[currentFlagName], word)
			expecting = expectingAnything
		default:
			return nil, fmt.Errorf("internal Error: not expecting state: %#v", expecting)
		}
	}
	if expecting == expectingFlagValue {
		return nil, fmt.Errorf("flag passed without value( %#v) . All flags must have one value passed. Flags can be repeated to accumulate values. Example: --level 9000", currentFlagName)
	}
	return res, nil
}

// fitToAppResult holds the result of fitToApp
// Exactly one of Section or Command should hold something. The other should be nil
type fitToAppResult struct {
	Section      *s.Section
	Command      *c.Command
	Action       c.Action
	AllowedFlags f.FlagMap
}

// fitToApp takes the command entered by a user and uses it to "walk" down the apps command tree
func fitToApp(rootSection s.Section, path []string) (*fitToAppResult, error) {
	// validate passed command and get available flags
	ftar := fitToAppResult{
		Section:      &rootSection,
		Command:      nil, // we start with a section, not a command
		AllowedFlags: rootSection.Flags,
		Action:       nil,
	}
	childCommands := rootSection.Commands
	childSections := rootSection.Sections
	for _, word := range path {
		if command, exists := childCommands[word]; exists {
			ftar.Command = &command
			ftar.Section = nil
			ftar.Action = command.Action
			// once we're in a commmand, we should be at the end of the path
			// commands have no child commands or child sections
			childCommands = nil
			childSections = nil
			for k, v := range command.Flags {
				// TODO: check if key exists already
				v.IsCommandFlag = true
				ftar.AllowedFlags[k] = v
			}
		} else if section, exists := childSections[word]; exists {
			ftar.Section = &section
			childCommands = section.Commands
			childSections = section.Sections
			for k, v := range command.Flags {
				// TODO: check if key exists already
				ftar.AllowedFlags[k] = v
			}
		} else {
			retErr := fmt.Errorf("expected command or section, but got %#v, try --help", word)
			return nil, retErr
		}
	}
	return &ftar, nil
}

// ParseResult holds the result of parsing the command line.
type ParseResult struct {
	// Path to the command invoked. Does not include executable name (os.Args[0])
	Path []string
	// PassedFlags holds the set flags!
	PassedFlags f.PassedFlags
	// Action holds the passed command's action to execute.
	Action c.Action
}

// Parse parses the args, but does not execute anything.
func (app *App) Parse(osArgs []string, lookup f.LookupFunc) (*ParseResult, error) {
	gar, err := gatherArgs(osArgs, app.helpFlagNames)
	if err != nil {
		return nil, err
	}

	ftar, err := fitToApp(app.rootSection, gar.Path)
	if err != nil {
		return nil, err
	}

	// fill the flags
	var configReader configreader.ConfigReader
	// get the value of a potential passed --config flag
	if app.configFlag != nil {
		// we're gonna make a config map out of this if everything goes well
		// so pass nil for that now
		err = app.configFlag.Resolve(app.configFlagName, gar.FlagStrs, nil)
		if err != nil {
			return nil, err
		}
		// NOTE: this *should* always be a string
		configPath := app.configFlag.Value.Get().(string)
		configReader, err = app.newConfigReader(configPath)
		if err != nil {
			return nil, fmt.Errorf("error reading config path ( %s ) : %w", configPath, err)
		}
	}

	// Loop over allowed flags for the passed command and try to resolve them
	for name, flag := range ftar.AllowedFlags {

		err = flag.Resolve(name, gar.FlagStrs, configReader)
		if err != nil {
			return nil, err
		}

		ftar.AllowedFlags[name] = flag
	}

	// add the config flag so both help and actions can see it
	if app.configFlag != nil {
		ftar.AllowedFlags[app.configFlagName] = *app.configFlag
	}

	// check for passed flags that arent' allowed
	if len(gar.FlagStrs) != 0 {
		return nil, fmt.Errorf("unrecognized flags: %v", gar.FlagStrs)
	}

	// TODO: check that all required flags are resolved! Not sure I have required flags yet :)

	// OK! Let's make the ParseResult for each case and gtfo
	if ftar.Section != nil && ftar.Command == nil {
		// no legit actions, just print the help
		pr := ParseResult{
			Action: app.sectionHelp(app.helpWriter, *ftar.Section, help.HelpInfo{AppName: app.name, Path: gar.Path, AvailableFlags: ftar.AllowedFlags, RootSection: app.rootSection}),
			// Action: app.sectionHelp(app.helpWriter, app.name, gar.Path, *ftar.Section, ftar.AllowedFlags),
		}
		return &pr, nil
	} else if ftar.Section == nil && ftar.Command != nil {
		if gar.HelpPassed {
			pr := ParseResult{
				Action: app.commandHelp(app.helpWriter, *ftar.Command, help.HelpInfo{AppName: app.name, Path: gar.Path, AvailableFlags: ftar.AllowedFlags, RootSection: app.rootSection}),

				// Action: app.commandHelp(app.helpWriter, app.name, gar.Path, *ftar.Command, ftar.AllowedFlags),
			}
			return &pr, nil
		} else {
			// TODO: change this
			fvs := make(f.PassedFlags)
			for name, flag := range ftar.AllowedFlags {
				if flag.SetBy != "" {
					fvs[name] = flag.Value.Get()
				}
			}

			pr := ParseResult{
				Path:        gar.Path,
				PassedFlags: fvs,
				Action:      ftar.Action,
			}
			return &pr, nil
		}
	} else {
		return nil, fmt.Errorf("internal Error: invalid parse state: currentSection == %v, currentCommand == %v", ftar.Section, ftar.Command)
	}
}

// Run parses the args, runs the action for the command passed,
// and returns any errors encountered.
func (app *App) Run(osArgs []string, lookup f.LookupFunc) error {
	pr, err := app.Parse(osArgs, lookup)
	if err != nil {
		return err
	}
	err = pr.Action(pr.PassedFlags)
	if err != nil {
		return err
	}
	return nil
}

// MustRun runs the app.
// If there's an error, it will be printed to stderr and os.Exit(1)
// will be called
func (app *App) MustRun(osArgs []string, lookup f.LookupFunc) {
	err := app.Run(osArgs, lookup)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func DictLookup(m map[string]string) f.LookupFunc {
	return func(key string) (string, bool) {
		val, exists := m[key]
		return val, exists
	}
}
