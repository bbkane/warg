package clide

import (
	"fmt"
	"log"
)

type CategoryMap = map[string]CategoryValue
type CommandMap = map[string]CommandValue
type FlagMap = map[string]FlagValue

type CategoryOpt = func(*CategoryValue)
type CommandOpt = func(*CommandValue)

type App struct {
	Name         string
	RootCategory CategoryValue
}

type CategoryValue struct {
	Flags      FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands   CommandMap
	Categories CategoryMap
}
type CommandValue struct {
	Flags FlagMap
}
type FlagValue struct {
	// Value holds what gets passed to the flag: --myflag value
	Value string
}

func AddCategory(name string, value CategoryValue) CategoryOpt {
	return func(app *CategoryValue) {
		if _, alreadyThere := app.Categories[name]; !alreadyThere {
			app.Categories[name] = value
		} else {
			log.Fatalf("category already exists: %#v\n", name)
		}
	}
}

func AddCommand(name string, value CommandValue) CategoryOpt {
	return func(app *CategoryValue) {
		if _, alreadyThere := app.Commands[name]; !alreadyThere {
			app.Commands[name] = value
		} else {
			log.Fatalf("command already exists: %#v\n", name)
		}
	}
}

func AddCategoryFlag(name string, value FlagValue) CategoryOpt {
	return func(app *CategoryValue) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}

	}
}

func AddCommandFlag(name string, value FlagValue) CommandOpt {
	return func(app *CommandValue) {
		if _, alreadyThere := app.Flags[name]; !alreadyThere {
			app.Flags[name] = value
		} else {
			log.Fatalf("flag already exists: %#v\n", name)
		}
	}
}

func WithCategory(name string, opts ...CategoryOpt) CategoryOpt {
	return AddCategory(name, NewCategory(opts...))
}

func WithCommand(name string, opts ...CommandOpt) CategoryOpt {
	return AddCommand(name, NewCommand(opts...))
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

func NewCommand(opts ...CommandOpt) CommandValue {
	category := CommandValue{
		Flags: make(map[string]FlagValue),
	}
	for _, opt := range opts {
		opt(&category)
	}
	return category
}

func (app *CategoryValue) Parse(args []string) ([]string, FlagMap, error) {

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
