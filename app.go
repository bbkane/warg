package warg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"

	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

type AppOpt = func(*App)

type ConfigMap = map[string]interface{}

// Unmarshaller turns a string into a map so we can index into it!
// Useful for configs who will read a file to get it
type Unmarshaller = func(string) (ConfigMap, error)

// JSONUnmarshaller tries to turn a filepath into a map[string]interface .
// It does NOT error if the file can not be read. If it did, then users would
// be forced to have a config before the app would work. TODO: is this the best method?
// It DOES error if the file can't be unmarshalled
// Note that all numbers in JSON are floats. So, no int flags if you use this encoder.
// Cast to an int after parsing if you like instead
func JSONUnmarshaller(filePath string) (map[string]interface{}, error) {
	// TODO: expand homedir?
	var m map[string]interface{}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// the file not existing is ok
		return m, nil
	}

	err = json.Unmarshal(content, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type App struct {
	// Config()
	configFlagName     string
	configUnmarshaller Unmarshaller
	configFlag         *f.Flag
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
	unmarshaller Unmarshaller,
	helpShort string,
	flagOpts ...f.FlagOpt,
) AppOpt {
	return func(app *App) {
		app.configFlagName = configFlagName
		app.configUnmarshaller = unmarshaller
		configFlag := f.NewFlag(helpShort, v.StringEmpty(), flagOpts...)
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
		app.helpFlagNames = []string{"--help", "-h"}
		app.sectionHelp = DefaultSectionHelp
		app.commandHelp = DefaultCommandHelp
	}
	// Version
	if len(app.versionFlagNames) == 0 {
		app.versionFlagNames = []string{"--version"}
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

// resolveFLag updates a flag's value from the command line, and then from the
// default value. flag should not be nil. deletes from flagStrs
func resolveFlag(flag *f.Flag, name string, flagStrs map[string][]string, configMap ConfigMap) error {
	// update from command line
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

	// update from config
	if flag.SetBy == "" && configMap != nil && flag.ConfigFromInterface != nil {
		i, exists, err := followPath(configMap, flag.ConfigPath)
		if err != nil {
			return err
		}
		if exists {
			v, err := flag.ConfigFromInterface(i)
			if err != nil {
				return err
			}
			flag.Value = v
			flag.SetBy = "config"
		}
	}

	// update from default
	if flag.SetBy == "" && len(flag.DefaultValues) > 0 {
		for _, v := range flag.DefaultValues {
			flag.Value.Update(v)
		}
		flag.SetBy = "appdefault"
	}

	return nil
}

// followPath takes a map and a path with elements separated by dots
// and retrieves the interface at the end of it. If the interface
// doesn't exist, then the bool value is false
func followPath(m ConfigMap, path string) (interface{}, bool, error) {
	pathSlice := strings.Split(path, ".")
	lastIndex := len(pathSlice) - 1
	var err error
	// step down the path
	for _, step := range pathSlice[:lastIndex] {
		nextIface, exists := m[step]
		if !exists {
			return nil, false, nil
		}
		nextMap, isMap := nextIface.(map[string]interface{})
		if !isMap {
			err = fmt.Errorf(
				"error: expected map[string]interface{} at %#v: got %#v",
				step,
				nextIface,
			)
			return nil, false, err
		}
		m = nextMap
	}

	step := pathSlice[lastIndex]
	val, exists := m[step]
	if !exists {
		return nil, false, err
	}

	return val, true, nil
}

func (app *App) Parse(osArgs []string) (*ParseResult, error) {
	gar, err := gatherArgs(osArgs, app.helpFlagNames, app.versionFlagNames)
	if err != nil {
		return nil, err
	}

	// special case versionFlag and exit early
	if gar.VersionPassed {
		pr := ParseResult{
			Action: func(_ map[string]v.Value) error {
				fmt.Print(app.version)
				return nil
			},
		}
		return &pr, nil
	}

	ftar, err := fitToApp(app.rootSection, gar.Path, gar.FlagStrs)
	if err != nil {
		return nil, err
	}

	// update the config flag :)
	var configMap ConfigMap
	if app.configFlag != nil {
		// we're gonna make a config map out of this if everything goes well
		// so pass nil for that now
		err = resolveFlag(app.configFlag, app.configFlagName, gar.FlagStrs, nil)
		if err != nil {
			return nil, err
		}
		// TODO: don't panic if not not a string. return an error :)

		configMap, err = app.configUnmarshaller(app.configFlag.Value.Get().(string))
		if err != nil {
			return nil, err
		}
	}

	// We need to loop over a map by value, so we can't modify it
	// in place :/
	for name, flag := range ftar.AllowedFlags {

		err = resolveFlag(&flag, name, gar.FlagStrs, configMap)
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
			vm := make(v.ValueMap)
			for name, flag := range ftar.AllowedFlags {
				if flag.SetBy != "" {
					vm[name] = flag.Value
				}
			}

			pr := ParseResult{
				PasssedPath: gar.Path,
				PassedFlags: vm,
				Action:      ftar.Action,
			}
			return &pr, nil
		}
	} else {
		return nil, fmt.Errorf("internal Error: invalid parse state: currentCategory == %v, currentCommand == %v", ftar.Section, ftar.Command)
	}
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
	return func(vm v.ValueMap) error {
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
	return func(vm v.ValueMap) error {
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
	PassedFlags v.ValueMap
	Action      c.Action
}
