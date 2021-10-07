package warg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	c "github.com/bbkane/warg/command"
	"github.com/bbkane/warg/configreader"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

type AppOpt = func(*App)

type App struct {
	// Config()
	configFlagName  string
	newConfigReader configreader.NewConfigReader
	configFlag      *f.Flag
	// Help()
	name          string
	helpFlagNames []string
	sectionHelp   SectionHelp
	commandHelp   CommandHelp
	// Version()
	version          string
	versionFlagNames []string
	// rootSection holds the good stuff!
	rootSection s.Section
}

func OverrideHelp(helpFlagNames []string, sectionHelp SectionHelp, commandHelp CommandHelp) AppOpt {
	return func(app *App) {
		app.sectionHelp = sectionHelp
		app.commandHelp = commandHelp
		app.helpFlagNames = helpFlagNames
		for _, n := range helpFlagNames {
			if !strings.HasPrefix(n, "-") {
				log.Panicf("helpFlags should start with '-': %#v\n", n)
			}
		}
	}
}

func OverrideVersion(versionFlagNames []string) AppOpt {
	return func(app *App) {
		app.versionFlagNames = versionFlagNames
	}
}

func ConfigFlag(
	configFlagName string,
	newConfigReader configreader.NewConfigReader,
	helpShort string,
	flagOpts ...f.FlagOpt,
) AppOpt {
	return func(app *App) {
		app.configFlagName = configFlagName
		app.newConfigReader = newConfigReader
		configFlag := f.NewFlag(helpShort, v.StringEmpty, flagOpts...)
		app.configFlag = &configFlag
	}
}

func New(name string, version string, rootSection s.Section, opts ...AppOpt) App {
	app := App{
		name:        name,
		rootSection: rootSection,
		version:     version,
	}
	for _, opt := range opts {
		opt(&app)
	}

	// Help
	if len(app.helpFlagNames) == 0 {
		OverrideHelp(
			[]string{"-h", "--help"},
			DefaultSectionHelp,
			DefaultCommandHelp,
		)(&app)
	}
	// Version
	if len(app.versionFlagNames) == 0 {
		OverrideVersion(
			[]string{"--version"},
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
	FlagStrs      map[string][]string
	VersionPassed bool
	HelpPassed    bool
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
// TODO: test this! Also, --help and --version do NOT require values
func gatherArgs(osArgs []string, helpFlagNames []string, versionFlagNames []string) (*gatherArgsResult, error) {
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
			if containsString(versionFlagNames, word) {
				res.VersionPassed = true
				// No need to do any more processing. Let's get out of here
				// NOTE: as is, this means that any number of categories can be passed. Not sure if I care...
				return res, nil
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
		return nil, fmt.Errorf("flag passed without value. All flags must have one value passed. Flags can be repeated to accumulate values. Example: --flag value")
	}
	return res, nil
}

type fitToAppResult struct {
	Section      *s.Section
	Command      *c.Command
	Action       c.Action
	AllowedFlags f.FlagMap
}

// fitToApp takes the command entered by a user and maps it to a command in the tree
func fitToApp(rootSection s.Section, path []string, flagStrs map[string][]string) (*fitToAppResult, error) {
	// validate passed command and get available flags
	ftar := fitToAppResult{
		Section:      &rootSection,
		AllowedFlags: rootSection.Flags,
		Command:      nil, // we start with a section, not a command
	}
	allowedCommands := rootSection.Commands
	allowedCategories := rootSection.Sections
	for _, word := range path {
		if command, exists := allowedCommands[word]; exists {
			ftar.Command = &command
			ftar.Section = nil
			ftar.Action = command.Action
			allowedCommands = nil   // commands terminate
			allowedCategories = nil // categories terminiate
			for k, v := range command.Flags {
				// TODO: check if key exists already
				ftar.AllowedFlags[k] = v
			}
		} else if category, exists := allowedCategories[word]; exists {
			ftar.Section = &category
			allowedCommands = category.Commands
			allowedCategories = category.Sections
			for k, v := range command.Flags {
				// TODO: check if key exists already
				ftar.AllowedFlags[k] = v
			}
		} else {
			return nil, fmt.Errorf("unexpected string: %#v", word)
		}
	}
	return &ftar, nil
}

func (app *App) Parse(osArgs []string) (*ParseResult, error) {
	gar, err := gatherArgs(osArgs, app.helpFlagNames, app.versionFlagNames)
	if err != nil {
		return nil, err
	}

	// special case versionFlag and exit early
	if gar.VersionPassed {
		pr := ParseResult{
			Action: func(_ f.FlagValues) error {
				fmt.Println(app.version)
				return nil
			},
		}
		return &pr, nil
	}

	ftar, err := fitToApp(app.rootSection, gar.Path, gar.FlagStrs)
	if err != nil {
		return nil, err
	}

	var configReader configreader.ConfigReader
	// get the value of a potential passed --config flag
	if app.configFlag != nil {
		// we're gonna make a config map out of this if everything goes well
		// so pass nil for that now
		err = app.configFlag.Resolve(app.configFlagName, gar.FlagStrs, nil)
		if err != nil {
			return nil, err
		}
		// TODO: don't panic if not not a string. return an error :)
		configReader, err = app.newConfigReader(app.configFlag.Value.Get().(string))
		if err != nil {
			return nil, err
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
			Action: app.sectionHelp(app.name, gar.Path, *ftar.Section, ftar.AllowedFlags),
		}
		return &pr, nil
	} else if ftar.Section == nil && ftar.Command != nil {
		if gar.HelpPassed {
			pr := ParseResult{
				Action: app.commandHelp(app.name, gar.Path, *ftar.Command, ftar.AllowedFlags),
			}
			return &pr, nil
		} else {
			// TODO: change this
			fvs := make(f.FlagValues)
			for name, flag := range ftar.AllowedFlags {
				if flag.SetBy != "" {
					fvs[name] = flag.Value.Get()
				}
			}

			pr := ParseResult{
				PasssedPath: gar.Path,
				PassedFlags: fvs,
				Action:      ftar.Action,
			}
			return &pr, nil
		}
	} else {
		return nil, fmt.Errorf("internal Error: invalid parse state: currentCategory == %v, currentCommand == %v", ftar.Section, ftar.Command)
	}
}

func (app *App) Run(osArgs []string) error {
	pr, err := app.Parse(osArgs)
	if err != nil {
		return err
	}
	err = pr.Action(pr.PassedFlags)
	if err != nil {
		return err
	}
	return err
}

// TODO: actually put this in :)
type CommandHelp = func(appName string, path []string, cur c.Command, flagMap f.FlagMap) c.Action

type SectionHelp = func(appName string, path []string, cur s.Section, flagMap f.FlagMap) c.Action

func DefaultCommandHelp(
	appName string,
	path []string,
	cur c.Command,
	flagMap f.FlagMap,
) c.Action {
	return func(_ f.FlagValues) error {
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()

		// Print top help section
		if cur.HelpLong == "" {
			fmt.Fprintf(f, "%s\n", cur.Help)
		} else {
			fmt.Fprintf(f, "%s\n", cur.Help)
		}

		fmt.Fprintln(f)

		fmt.Fprintf(f, "Flags:\n")
		fmt.Fprintln(f)
		{
			keys := make([]string, 0, len(flagMap))
			for k := range flagMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				flag := flagMap[k]
				fmt.Fprintf(f, "  %s : %s\n", k, flag.Help)
				if flag.ConfigPath != "" {
					fmt.Fprintf(f, "    configpath : %s\n", flag.ConfigPath)
				}
				if flag.SetBy != "" {
					fmt.Fprintf(f, "    value : %s\n", flag.Value)
					fmt.Fprintf(f, "    setby : %s\n", flag.SetBy)
				}
				fmt.Fprintln(f)
			}
		}
		return nil
	}
}

func DefaultSectionHelp(
	appName string,
	path []string,
	cur s.Section,
	flagMap f.FlagMap,
) c.Action {
	return func(_ f.FlagValues) error {
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()

		// Print top help section
		if cur.HelpLong == "" {
			fmt.Fprintf(f, "%s\n", cur.Help)
		} else {
			fmt.Fprintf(f, "%s\n", cur.Help)
		}

		fmt.Fprintln(f)

		// Print sections
		fmt.Fprintf(f, "Sections:\n")
		{
			keys := make([]string, 0, len(cur.Sections))
			for k := range cur.Sections {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(f, "  %s : %s\n", k, cur.Sections[k].Help)
			}
		}

		fmt.Fprintln(f)

		// Print commands
		fmt.Fprintf(f, "Commands:\n")
		{
			keys := make([]string, 0, len(cur.Commands))
			for k := range cur.Commands {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(f, "  %s : %s\n", k, cur.Commands[k].Help)
			}
		}

		// TODO: print examples once we have them :)
		return nil
	}
}

type ParseResult struct {
	PasssedPath []string
	PassedFlags f.FlagValues
	Action      c.Action
}
