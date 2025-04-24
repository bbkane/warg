package command

import (
	"errors"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/wargcore"
)

// A CommandOpt customizes a Command
type CommandOpt func(*wargcore.Command)

// DoNothing is a command action that simply returns an error.
// Useful for prototyping
func DoNothing(_ wargcore.Context) error {
	return errors.New("NOTE: replace this command.DoNothing call")
}

// New builds a Command
func New(helpShort string, action wargcore.Action, opts ...CommandOpt) wargcore.Command {
	command := wargcore.Command{
		HelpShort:            helpShort,
		Action:               action,
		CompletionCandidates: DefaultCompletionCandidates,
		Flags:                make(wargcore.FlagMap),
		Footer:               "",
		HelpLong:             "",
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

func DefaultCompletionCandidates(cmdCtx wargcore.Context) (*completion.Candidates, error) {
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
		// scalar flags set by passed arg can't be appended to or overridden, so don't suggest them
		val, isScalar := cmdCtx.ParseState.FlagValues[name].(value.ScalarValue)
		if isScalar && val.UpdatedBy() == value.UpdatedByFlag {
			continue
		}
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

func CompletionCandidates(completionCandidatesFunc wargcore.CompletionCandidates) CommandOpt {
	return func(flag *wargcore.Command) {
		flag.CompletionCandidates = completionCandidatesFunc
	}
}

// Flag adds an existing flag to a Command. It panics if a flag with the same name exists
func Flag(name string, value wargcore.Flag) CommandOpt {
	return func(com *wargcore.Command) {
		com.Flags.AddFlag(name, value)
	}
}

// FlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func FlagMap(flagMap wargcore.FlagMap) CommandOpt {
	return func(com *wargcore.Command) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func NewFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) CommandOpt {
	return Flag(name, flag.New(helpShort, empty, opts...))
}

// Footer adds an Help string to the command - useful from a help function
func Footer(footer string) CommandOpt {
	return func(cat *wargcore.Command) {
		cat.Footer = footer
	}
}

// HelpLong adds an Help string to the command - useful from a help function
func HelpLong(helpLong string) CommandOpt {
	return func(cat *wargcore.Command) {
		cat.HelpLong = helpLong
	}
}
