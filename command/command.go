package command

import (
	"errors"

	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value"
)

// A CommandOpt customizes a Command
type CommandOpt func(*cli.Command)

// DoNothing is a command action that simply returns an error.
// Useful for prototyping
func DoNothing(_ cli.Context) error {
	return errors.New("NOTE: replace this command.DoNothing call")
}

// New builds a Command
func New(helpShort string, action cli.Action, opts ...CommandOpt) cli.Command {
	command := cli.Command{
		HelpShort:            helpShort,
		Action:               action,
		CompletionCandidates: DefaultCompletionCandidates,
		Flags:                make(cli.FlagMap),
		Footer:               "",
		HelpLong:             "",
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

func DefaultCompletionCandidates(cmdCtx cli.Context) (*completion.Candidates, error) {
	// TODO: flag name completion ideas that will actually use the full parse above
	//  - if a scalar flag has been passed by arg, don't suggest it again (as args override everything else)
	//  - if the flag is required and is not set, suggest it first
	//  - suggest command flags before global flags
	//  - let the flags define rank or priority for completion order
	candidates := &completion.Candidates{
		Type:   completion.Type_ValuesDescriptions,
		Values: []completion.Candidate{},
	}
	// command flags
	for _, name := range cmdCtx.ParseState.CurrentCommand.Flags.SortedNames() {
		candidates.Values = append(candidates.Values, completion.Candidate{
			Name:        string(name),
			Description: string(cmdCtx.ParseState.CurrentCommand.Flags[name].HelpShort),
		})
	}
	// global flags
	for _, name := range cmdCtx.App.GlobalFlags.SortedNames() {
		candidates.Values = append(candidates.Values, completion.Candidate{
			Name:        string(name),
			Description: string(cmdCtx.App.GlobalFlags[name].HelpShort),
		})
	}
	return candidates, nil
}

func CompletionCandidates(completionCandidatesFunc func(cli.Context) (*completion.Candidates, error)) CommandOpt {
	return func(flag *cli.Command) {
		flag.CompletionCandidates = completionCandidatesFunc
	}
}

// Flag adds an existing flag to a Command. It panics if a flag with the same name exists
func Flag(name string, value cli.Flag) CommandOpt {
	return func(com *cli.Command) {
		com.Flags.AddFlag(name, value)
	}
}

// FlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func FlagMap(flagMap cli.FlagMap) CommandOpt {
	return func(com *cli.Command) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func NewFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) CommandOpt {
	return Flag(name, flag.New(helpShort, empty, opts...))
}

// Footer adds an Help string to the command - useful from a help function
func Footer(footer string) CommandOpt {
	return func(cat *cli.Command) {
		cat.Footer = footer
	}
}

// HelpLong adds an Help string to the command - useful from a help function
func HelpLong(helpLong string) CommandOpt {
	return func(cat *cli.Command) {
		cat.HelpLong = helpLong
	}
}
