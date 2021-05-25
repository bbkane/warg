///usr/bin/true; exec /usr/bin/env go test ./...
///usr/bin/true; exec /usr/bin/env go run "$0" .
package main

import (
	"fmt"
	"log"
)

type FlagMap = map[string]FlagValue
type CommandMap = map[string]CommandValue
type CategoryMap = map[string]CategoryValue

type FlagValue struct {
	// Value holds what gets passed to the flag: --myflag value
	Value string
}

type CommandValue struct {
	Flags FlagMap
}

type CategoryValue struct {
	Flags      FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands   CommandMap
	Categories CategoryMap
}

type App struct {
	Name       string
	Flags      FlagMap
	Commands   CommandMap
	Categories CategoryMap
}

type AppOpt = func(*App)
type CategoryOpt = func(*CategoryValue)
type CommandOpt = func(*CommandValue)

func CommandFlag(name string, value FlagValue) CommandOpt {
	return func(app *CommandValue) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}
	}
}

func CategoryFlag(name string, value FlagValue) CategoryOpt {
	return func(app *CategoryValue) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}

	}
}

func AppFlag(name string, value FlagValue) AppOpt {
	return func(app *App) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}

	}
}

func CategoryCategory(name string, value CategoryValue) CategoryOpt {
	return func(app *CategoryValue) {
		if _, alreadyThere := app.Categories[name]; !alreadyThere {
			app.Categories[name] = value
		} else {
			log.Fatalf("category already exists: %#v\n", name)
		}
	}
}

func AppCategoryWith(name string, opts ...CategoryOpt) AppOpt {
	return AppCategory(name, NewCategory(opts...))
}

func AppCategory(name string, value CategoryValue) AppOpt {
	return func(app *App) {
		if _, alreadyThere := app.Categories[name]; !alreadyThere {
			app.Categories[name] = value
		} else {
			log.Fatalf("category already exists: %#v\n", name)
		}
	}
}

func CategoryCommand(name string, value CommandValue) CategoryOpt {
	return func(app *CategoryValue) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Fatalf("command already exists: %#v\n", name)
		}
	}
}

func AppCommandWith(name string, opts ...CommandOpt) AppOpt {
	return AppCommand(name, NewCommand(opts...))
}

func AppCommand(name string, value CommandValue) AppOpt {
	return func(app *App) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Fatalf("command already exists: %#v\n", name)
		}
	}
}

func NewCommand(opts ...CommandOpt) CommandValue {
	category := CommandValue{
		Flags: make(map[string]FlagValue),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func NewCategory(opts ...CategoryOpt) CategoryValue {
	category := CategoryValue{
		Flags:      make(map[string]FlagValue),
		Categories: make(map[string]CategoryValue),
		Commands:   make(map[string]CommandValue),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func NewApp(name string, opts ...AppOpt) App {
	app := App{
		Name:       name,
		Flags:      make(map[string]FlagValue),
		Categories: make(map[string]CategoryValue),
		Commands:   make(map[string]CommandValue),
	}
	for _, opt := range opts {
		opt(&app)
	}
	return app
}

func (app *App) Parse(args []string) ([]string, FlagMap, error) {

	// TODO: I'd like flags to be callable in any order after their command is called
	// so instead of reassigning allowedFlags, merge it with the new one
	allowedFlags := app.Flags
	allowedCommands := app.Commands
	allowedCategories := app.Categories
	passedFlags := make(FlagMap)
	passedCommand := make([]string, 0, len(args)-1)
	for i := 1; i < len(args); i = i + 1 {
		val := args[i]
		if _, ok := allowedFlags[val]; ok {
			passedFlags[val] = FlagValue{Value: args[i+1]} // TODO: what if someone passes a flag without a value
			i += 1
		} else if command, ok := allowedCommands[val]; ok {
			passedCommand = append(passedCommand, val)
			allowedFlags = command.Flags
			allowedCommands = nil
			allowedCategories = nil
		} else if category, ok := allowedCategories[val]; ok {
			passedCommand = append(passedCommand, val)
			allowedFlags = category.Flags
			allowedCommands = category.Commands
			allowedCategories = category.Categories
		} else {
			return nil, nil, fmt.Errorf("unexpected string: %#v\n", val)
		}
	}
	return passedCommand, passedFlags, nil
}
