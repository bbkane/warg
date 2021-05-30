package clide

import (
	"fmt"
	"log"
)

type CategoryMap = map[string]Category
type CommandMap = map[string]Command
type FlagMap = map[string]Flag
type ValueMap = map[string]Value

type CategoryOpt = func(*Category)
type CommandOpt = func(*Command)

type App struct {
	Name         string
	RootCategory Category
}

type Category struct {
	Flags      FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	Commands   CommandMap
	Categories CategoryMap
}
type Command struct {
	Flags FlagMap
}
type Flag struct {
	// Value holds what gets passed to the flag: --myflag value
	// and should be initialized to the empty value
	Value Value
	// Default will be shoved into Value if needed
	// can be nil
	// TODO: actually use this
	Default Value
	// IsSet should be set when the flag is set so defaults don't override something
	IsSet bool
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

func WithCategory(name string, opts ...CategoryOpt) CategoryOpt {
	return AddCategory(name, NewCategory(opts...))
}

func WithCommand(name string, opts ...CommandOpt) CategoryOpt {
	return AddCommand(name, NewCommand(opts...))
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

func (app *Category) Parse(args []string) ([]string, ValueMap, error) {

	// TODO: I'd like flags to be callable in any order after their command is called
	// so instead of reassigning allowedFlags, merge it with the new one
	allowedFlags := app.Flags
	allowedCommands := app.Commands
	allowedCategories := app.Categories
	passedFlagValues := make(ValueMap)
	passedCommand := make([]string, 0, len(args)-1)
	for i := 1; i < len(args); i = i + 1 {
		str := args[i]
		if currFlag, ok := allowedFlags[str]; ok {
			passedFlagValues[str] = currFlag.Value
			valueToParse := args[i+1] // TODO: gracefully handle someone passing a flag without a value
			err := currFlag.Value.Update(valueToParse)
			if err != nil {
				return nil, nil, fmt.Errorf(
					"flag: %#v: flag parse error for value : %#v: %#v\n",
					str,
					valueToParse,
					err,
				)
			}
			i += 1
		} else if command, ok := allowedCommands[str]; ok {
			passedCommand = append(passedCommand, str)
			allowedFlags = command.Flags
			allowedCommands = nil
			allowedCategories = nil
		} else if category, ok := allowedCategories[str]; ok {
			passedCommand = append(passedCommand, str)
			allowedFlags = category.Flags
			allowedCommands = category.Commands
			allowedCategories = category.Categories
		} else {
			return nil, nil, fmt.Errorf("unexpected string: %#v\n", str)
		}
	}
	return passedCommand, passedFlagValues, nil
}
