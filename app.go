package clide

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type AppOpt = func(*App)

type App struct {
	// Help()
	name          string
	helpFlagNames []string
	// Version()
	version          string
	versionFlagNames []string
	// Categories
	rootCategory Category
}

func EnableHelpFlag(helpFlagNames []string, appName string) AppOpt {
	return func(app *App) {
		app.name = appName

		app.helpFlagNames = helpFlagNames
		for _, n := range helpFlagNames {
			if !strings.HasPrefix(n, "-") {
				log.Panicf("helpFlags should start with '-': %#v\n", n)
			}
		}
	}
}

func EnableVersionFlag(versionFlagNames []string, version string) AppOpt {
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
	currentCategory := &(app.rootCategory)
	var currentCommand *Command = nil
	allowedFlags := currentCategory.Flags
	allowedCommands := currentCategory.Commands
	allowedCategories := currentCategory.Categories
	for _, word := range gatherArgsResult.CommandPath {
		if command, exists := allowedCommands[word]; exists {
			currentCommand = &command
			currentCategory = nil
			pr.Action = command.Action
			allowedCommands = nil   // commands terminate
			allowedCategories = nil // categories terminiate
			for k, v := range command.Flags {
				// TODO: check if key exists already
				allowedFlags[k] = v
			}
		} else if category, exists := allowedCategories[word]; exists {
			currentCategory = &category
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
		flag.SetBy = "commandline"
		// I would think this woudn't be necessary...
		// I think because this isn't explicitly a pointer its passed by value? I'm too used to Python...
		// TODO: look into this more :)
		allowedFlags[name] = flag
	}
	// fmt.Printf("allowed flags: %#v\n", allowedFlags)

	// update unset flags backup values
	for name, flag := range allowedFlags {
		// update from default value
		if flag.SetBy == "" && flag.Default != nil {
			flag.Value = flag.Default
			flag.SetBy = "appdefault"
			allowedFlags[name] = flag
		}
	}

	if gatherArgsResult.HelpPassed {
		if currentCategory != nil && currentCommand == nil {
			pr.Action = func(_ ValueMap) error {
				f := bufio.NewWriter(os.Stdout)
				defer f.Flush()
				// let's assume that HelpLong doesn't exist
				fmt.Fprintf(f, "Current Category:\n")
				fmt.Fprintf(f, "  %s: %s\n", gatherArgsResult.CommandPath, currentCategory.HelpShort)
				fmt.Fprintf(f, "Subcategories:\n")
				// TODO: sort these :)
				for name, value := range currentCategory.Categories {
					fmt.Fprintf(f, "  %s: %s\n", name, value.HelpShort)
				}
				// TODO: sort these too :)
				fmt.Fprintf(f, "Commands:\n")
				for name, value := range currentCategory.Commands {
					fmt.Fprintf(f, "  %s: %s\n", name, value.HelpShort)
				}
				return nil
			}
		} else if currentCommand != nil && currentCategory == nil {
			pr.Action = func(_ ValueMap) error {
				// TODO
				fmt.Printf("TODO :)")
				return nil
			}
		} else {
			return nil, fmt.Errorf("Internal Error: invalid help state: currentCategory == %v, currentCommand == %v\n", currentCategory, currentCommand)
		}

		return pr, nil
	}

	// make some values!
	for name, flag := range allowedFlags {
		if flag.SetBy != "" {
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
