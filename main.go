///usr/bin/true; exec /usr/bin/env go test ./...
///usr/bin/true; exec /usr/bin/env go run "$0" .
package main

import (
	"fmt"
)

type FlagMap = map[string]Flag
type CommandMap = map[string]Command
type CategoryMap = map[string]Category

type Flag struct {
	// Value holds what gets passed to the flag: --myflag value
	Value string
}

type Command struct {
	Flags FlagMap
}

type Category struct {
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

type AppOpt = func(*App) error

func AppFlag(name string, value Flag) AppOpt {
	return func(app *App) error {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
			return nil
		} else {
			return fmt.Errorf("flag already entered: %#v\n", name)
		}

	}
}

func NewApp(name string, opts ...AppOpt) (*App, error) {
	rootCmd := App{
		Name:       name,
		Flags:      make(map[string]Flag),
		Categories: make(map[string]Category),
		Commands:   make(map[string]Command),
	}
	for _, opt := range opts {
		err := opt(&rootCmd)
		if err != nil {
			return nil, err
		}
	}
	return &rootCmd, nil
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
			passedFlags[val] = Flag{Value: args[i+1]} // TODO: what if someone passes a flag without a value
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
