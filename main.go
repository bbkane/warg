///usr/bin/true; exec /usr/bin/env go run "$0" .
package main

import (
	"fmt"
)

type FlagMap = map[string]Flag
type LeafCommandMap = map[string]LeafCommand
type SubCommandMap = map[string]SubCommand

type Flag struct {
	// Value holds what gets passed to the flag: --myflag value
	Value string
}

type LeafCommand struct {
	Flags FlagMap
}

type SubCommand struct {
	Flags        FlagMap // Do subcommands need flags? leaf commands are the ones that do work....
	LeafCommands LeafCommandMap
	SubCommands  SubCommandMap
}

type RootCommand struct {
	Value        string
	Flags        FlagMap
	LeafCommands LeafCommandMap
	SubCommands  SubCommandMap
}

func (command *RootCommand) Parse(args []string) ([]string, FlagMap, error) {

	// TODO: I'd like flags to be callable in any order after their command is called
	// so instead of reassigning allowedFlags, merge it with the new one
	allowedFlags := command.Flags
	allowedLeafCommands := command.LeafCommands
	allowedSubCommands := command.SubCommands
	passedFlags := make(FlagMap)
	passedCommand := make([]string, 0, len(args)-1)
	for i := 1; i < len(args); i = i + 1 {
		val := args[i]
		if _, ok := allowedFlags[val]; ok {
			passedFlags[val] = Flag{Value: args[i+1]} // TODO: what if someone passes a flag without a value
			i += 1
		} else if leafCommand, ok := allowedLeafCommands[val]; ok {
			passedCommand = append(passedCommand, val)
			allowedFlags = leafCommand.Flags
			allowedLeafCommands = nil
			allowedSubCommands = nil
		} else if subCommand, ok := allowedSubCommands[val]; ok {
			passedCommand = append(passedCommand, val)
			allowedFlags = subCommand.Flags
			allowedLeafCommands = subCommand.LeafCommands
			allowedSubCommands = subCommand.SubCommands
		} else {
			return nil, nil, fmt.Errorf("unexpected string: %#v\n", val)
		}
	}
	return passedCommand, passedFlags, nil
}

func main() {

	command := RootCommand{
		Value: "rc",
		Flags: FlagMap{
			"--rcf1": Flag{},
		},
		LeafCommands: LeafCommandMap{},
		SubCommands: SubCommandMap{
			"sc1": SubCommand{
				Flags: FlagMap{},
				LeafCommands: LeafCommandMap{
					"lc1": LeafCommand{
						Flags: FlagMap{
							"--lc1f1": Flag{},
						},
					},
				},
			},
		},
	}

	args := []string{"rc", "sc1", "lc1", "--lc1f1", "flagarg"}
	// args = []string{"rc", "--unexpected", "sc1", "lc1", "--lc1f1", "flagarg"}

	passedCommand, passedFlags, err := command.Parse(args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", passedCommand)
	fmt.Printf("%#v\n", passedFlags)
}
