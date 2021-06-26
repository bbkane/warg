package warg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

type AppOpt = func(*App)

// Unmarshaller turns a string into a map so we can index into it!
// Useful for configs who will read a file to get it
type Unmarshaller = func(string) (map[string]interface{}, error)

func JSONUnmarshaller(filePath string) (map[string]interface{}, error) {
	// TODO: expand homedir?
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
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
	// Version()
	version          string
	versionFlagNames []string
	// rootSection holds the good stuff!
	rootSection s.Section
}

func OverrideHelp(helpFlagNames []string) AppOpt {
	return func(app *App) {

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

func AddRootSection(rootSection s.Section) AppOpt {
	return func(app *App) {
		app.rootSection = rootSection
	}
}

func WithRootSection(helpShort string, opts ...s.SectionOpt) AppOpt {
	return func(app *App) {
		app.rootSection = s.NewSection(helpShort, opts...)
	}
}

func Config(
	configFlagName string,
	unmarshaller Unmarshaller,
	helpShort string,
	flagOpts ...f.FlagOpt,
) AppOpt {
	return func(app *App) {
		app.configFlagName = configFlagName
		app.configUnmarshaller = unmarshaller
		configFlag := f.NewFlag(helpShort, v.NewEmptyStringValue(), flagOpts...)
		app.configFlag = &configFlag
	}
}

func New(name string, version string, opts ...AppOpt) App {
	app := App{
		name:    name,
		version: version,
	}
	for _, opt := range opts {
		opt(&app)
	}
	// stitch up some "optional" parameters I'm expecting
	// RootSection
	if app.rootSection.Commands == nil {
		app.rootSection = s.NewSection("")
	}
	// Config - if passed, add to flags
	if app.configFlag != nil {
		app.rootSection.Flags[app.configFlagName] = *app.configFlag
	}

	// Help
	if len(app.helpFlagNames) == 0 {
		app.helpFlagNames = []string{"--help", "-h"}
		// TODO: custom help functions
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
			return nil, fmt.Errorf("Internal Error: not expecting state: %#v\n", expecting)
		}
	}
	if expecting == expectingFlagValue {
		return nil, fmt.Errorf("Flag passed without value. All flags must have one value passed. Flags can be repeated to accumulate values. Example: --flag value")
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
			return nil, fmt.Errorf("unexpected string: %#v\n", word)
		}
	}
	return &ftar, nil
}

func (app *App) Parse(osArgs []string) (*ParseResult, error) {
	gar, err := gatherArgs(osArgs, app.helpFlagNames, app.versionFlagNames)
	if err != nil {
		return nil, err
	}

	pr := &ParseResult{
		PasssedPath: gar.Path,
		PassedFlags: make(v.ValueMap),
		Action:      nil,
	}

	// special case versionFlag and exit early
	if gar.VersionPassed {
		pr.Action = func(_ map[string]v.Value) error {
			fmt.Print(app.version)
			return nil
		}
		return pr, nil
	}

	ftar, err := fitToApp(app.rootSection, gar.Path, gar.FlagStrs)
	if err != nil {
		return nil, err
	}

	pr.Action = ftar.Action

	for name, flag := range ftar.AllowedFlags {

		// update from command line
		strValues, exists := gar.FlagStrs[name]
		if exists {
			for _, v := range strValues {
				flag.Value.Update(v)
			}
			flag.SetBy = "commandline"
			// if they aren't all used
			delete(gar.FlagStrs, name)
		}

		// TODO: update from config

		// update from default
		if flag.SetBy == "" && flag.Default != nil {
			flag.Value = flag.Default
			flag.SetBy = "appdefault"
		}
		// I think this is legit :)
		ftar.AllowedFlags[name] = flag
	}

	// check for passed flags that arent' allowed
	if len(gar.FlagStrs) != 0 {
		return nil, fmt.Errorf("Unrecognized flags: %v\n", gar.FlagStrs)
	}

	if gar.HelpPassed {
		if ftar.Section != nil && ftar.Command == nil {
			pr.Action = DefaultCategoryHelp(app.name, gar.Path, *ftar.Section)
		} else if ftar.Command != nil && ftar.Section == nil {
			pr.Action = func(_ v.ValueMap) error {
				// TODO
				fmt.Printf("TODO :)")
				return nil
			}
		} else {
			return nil, fmt.Errorf("Internal Error: invalid help state: currentCategory == %v, currentCommand == %v\n", ftar.Section, ftar.Command)
		}

		return pr, nil
	}

	// make some values!
	for name, flag := range ftar.AllowedFlags {
		if flag.SetBy != "" {
			pr.PassedFlags[name] = flag.Value
		}
	}
	return pr, nil
}

func DefaultCategoryHelp(
	appName string,
	path []string,
	currentCategory s.Section,
) c.Action {
	return func(vm v.ValueMap) error {
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		// let's assume that HelpLong doesn't exist
		fmt.Fprintf(f, "Current Category:\n")
		totalPath := appName + " " + strings.Join(path, " ")
		fmt.Fprintf(f, "  %s: %s\n", totalPath, currentCategory.HelpShort)
		fmt.Fprintf(f, "Subcategories:\n")
		// TODO: sort these :)
		for name, value := range currentCategory.Sections {
			fmt.Fprintf(f, "  %s: %s\n", name, value.HelpShort)
		}
		// TODO: sort these too :)
		fmt.Fprintf(f, "Commands:\n")
		for name, value := range currentCategory.Commands {
			fmt.Fprintf(f, "  %s: %s\n", name, value.HelpShort)
		}
		return nil
	}
}

type ParseResult struct {
	PasssedPath []string
	PassedFlags v.ValueMap
	Action      c.Action
}
