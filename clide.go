package clide

import (
	"fmt"
	"log"
	"strings"
)

type Action = func(ValueMap) error

type CategoryMap = map[string]Category
type CommandMap = map[string]Command
type FlagMap = map[string]Flag
type ValueMap = map[string]Value

type AppOpt = func(*App)
type CategoryOpt = func(*Category)
type CommandOpt = func(*Command)

type App struct {
	// Help()
	name          string
	description   string
	helpFlagNames []string
	// Version()
	version          string
	versionFlagNames []string
	// Categories
	rootCategory Category
}

type Category struct {
	Flags      FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands   CommandMap
	Categories CategoryMap
}
type Command struct {
	Action Action

	Flags FlagMap
}

type Flag struct {
	// Default will be shoved into Value if needed
	// can be nil
	// TODO: actually use this
	Default Value
	// IsSet should be set when the flag is set so defaults don't override something
	IsSet bool
	// Value holds what gets passed to the flag: --myflag value
	// and should be initialized to the empty value
	Value Value
}

func AddCategory(name string, value Category) CategoryOpt {
	return func(app *Category) {
		if _, alreadyThere := app.Categories[name]; !alreadyThere {
			app.Categories[name] = value
		} else {
			log.Fatalf("category already exists: %#v\n", name)
		}
	}
}

func AddCommand(name string, value Command) CategoryOpt {
	return func(app *Category) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Fatalf("command already exists: %#v\n", name)
		}
	}
}

func AddCategoryFlag(name string, value Flag) CategoryOpt {
	return func(app *Category) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}

	}
}

func AddCommandFlag(name string, value Flag) CommandOpt {
	return func(app *Command) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}
	}
}

func WithAction(action Action) CommandOpt {
	return func(cmd *Command) {
		cmd.Action = action
	}
}

func WithCategory(name string, opts ...CategoryOpt) CategoryOpt {
	return AddCategory(name, NewCategory(opts...))
}

func WithCommand(name string, opts ...CommandOpt) CategoryOpt {
	return AddCommand(name, NewCommand(opts...))
}

func AppHelp(helpFlagNames []string, appName string, appDescription string) AppOpt {
	return func(app *App) {
		app.name = appName
		app.description = appDescription
		app.helpFlagNames = helpFlagNames
		for _, n := range helpFlagNames {
			if !strings.HasPrefix(n, "-") {
				log.Panicf("helpFlags should start with '-': %#v\n", n)
			}
		}
	}
}

func AppVersion(versionFlagNames []string, version string) AppOpt {
	return func(app *App) {
		app.versionFlagNames = versionFlagNames
		app.version = version
	}
}

func AppRootCategory(opts ...CategoryOpt) AppOpt {
	return func(app *App) {
		app.rootCategory = NewCategory(opts...)
	}
}

func NewApp(opts ...AppOpt) App {
	app := App{}
	for _, opt := range opts {
		opt(&app)
	}
	// TODO: will it panic if we try to Parse an empty category?
	return app
}

func NewCategory(opts ...CategoryOpt) Category {
	category := Category{
		Flags:      make(map[string]Flag),
		Categories: make(map[string]Category),
		Commands:   make(map[string]Command),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func NewCommand(opts ...CommandOpt) Command {
	category := Command{
		Flags: make(map[string]Flag),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

type gatherArgsResult struct {
	// Appname holds os.Args[0]
	AppName string
	// CommandPath holds the path to the current command
	CommandPath []string
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
			// TODO: search for --help

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
				res.CommandPath = append(res.CommandPath, word)
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

func (app *App) Parse(osArgs []string) (*ParseResult, error) {
	gatherArgsResult, err := gatherArgs(osArgs, app.helpFlagNames, app.versionFlagNames)
	if err != nil {
		return nil, err
	}

	pr := &ParseResult{
		PassedCmd:   gatherArgsResult.CommandPath,
		PassedFlags: make(ValueMap),
		Action:      nil,
	}

	// special case versionFlag
	if gatherArgsResult.VersionPassed {
		pr.Action = func(_ map[string]Value) error {
			fmt.Print(app.version)
			return nil
		}
		return pr, nil
	}

	// validate passed command and get available flags
	current := app.rootCategory
	allowedFlags := current.Flags
	allowedCommands := current.Commands
	allowedCategories := current.Categories
	for _, word := range gatherArgsResult.CommandPath {
		if command, exists := allowedCommands[word]; exists {
			pr.Action = command.Action
			allowedCommands = nil   // commands terminate
			allowedCategories = nil // categories terminiate
			for k, v := range command.Flags {
				// TODO: check if key exists already
				allowedFlags[k] = v
			}
		} else if category, exists := allowedCategories[word]; exists {
			allowedCommands = category.Commands
			allowedCategories = category.Categories
			for k, v := range command.Flags {
				// TODO: check if key exists already
				allowedFlags[k] = v
			}
		} else {
			return nil, fmt.Errorf("unexpected string: %#v\n", word)
		}
	}

	// fmt.Printf("allowed flags: %#v\n", allowedFlags)

	// update flags with passed values and ensure that no extra flags were passed
	// TODO: ensure passed flags match available flags, only aggregrate flags passed multiple times, required flags make it
	for name, passed := range gatherArgsResult.FlagStrs {
		flag, exists := allowedFlags[name]
		if !exists {
			return nil, fmt.Errorf("Unrecognized flag: %#v\n", name)
		}
		// TODO: check for repeated flags that aren't supposed to be repeated
		for _, str := range passed {
			flag.Value.Update(str)
		}
		flag.IsSet = true
		// I would think this woudn't be necessary...
		// I think because this isn't explicitly a pointer its passed by value? I'm too used to Python...
		// TODO: look into this more :)
		allowedFlags[name] = flag
	}
	// fmt.Printf("allowed flags: %#v\n", allowedFlags)

	// update unset flags backup values
	for name, flag := range allowedFlags {
		// update from default value
		if flag.IsSet == false && flag.Default != nil {
			flag.Value = flag.Default
			flag.IsSet = true
			allowedFlags[name] = flag
		}
	}

	// TODO: set action to print --help if needed and return

	// make some values!
	for name, flag := range allowedFlags {
		if flag.IsSet == true {
			pr.PassedFlags[name] = flag.Value
		}
	}
	return pr, nil
}

type ParseResult struct {
	PassedCmd   []string
	PassedFlags ValueMap
	Action      Action
}
