package command

import (
	"errors"
	"fmt"
	"strings"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/wargcore"
)

// A CommandOpt customizes a Command
type CommandOpt func(*wargcore.Cmd)

// DoNothing is a command action that simply returns an error.
// Useful for prototyping
func DoNothing(_ wargcore.Context) error {
	return errors.New("NOTE: replace this command.DoNothing call")
}

// NewCmd builds a Cmd
func NewCmd(helpShort string, action wargcore.Action, opts ...CommandOpt) wargcore.Cmd {
	command := wargcore.Cmd{
		HelpShort:   helpShort,
		Action:      action,
		Completions: DefaultCmdCompletions,
		Flags:       make(wargcore.FlagMap),
		Footer:      "",
		HelpLong:    "",
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

func DefaultCmdCompletions(cmdCtx wargcore.Context) (*completion.Candidates, error) {
	// FZF (or maybe zsh) auto-sorts by alphabetical order, so no need to get fancy with the following ideas
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
		var valStr string
		// TODO: does it matter if valstring is a large list?
		if cmdCtx.ParseState.FlagValues[name].UpdatedBy() != value.UpdatedByUnset {
			valStr = fmt.Sprint(cmdCtx.ParseState.FlagValues[name].Get())
			valStr = strings.ReplaceAll(valStr, "\n", " ")
			valStr = " (" + valStr + ")"
		}

		candidates.Values = append(candidates.Values, completion.Candidate{
			Name:        string(name),
			Description: string(cmdCtx.ParseState.CurrentCommand.Flags[name].HelpShort) + valStr,
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

func CmdCompletions(CompletionsFunc wargcore.CompletionsFunc) CommandOpt {
	return func(flag *wargcore.Cmd) {
		flag.Completions = CompletionsFunc
	}
}

// ChildFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func ChildFlag(name string, value wargcore.Flag) CommandOpt {
	return func(com *wargcore.Cmd) {
		com.Flags.AddFlag(name, value)
	}
}

// ChildFlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func ChildFlagMap(flagMap wargcore.FlagMap) CommandOpt {
	return func(com *wargcore.Cmd) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewChildFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func NewChildFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...flag.FlagOpt) CommandOpt {
	return ChildFlag(name, flag.NewFlag(helpShort, empty, opts...))
}

// CmdFooter adds an Help string to the command - useful from a help function
func CmdFooter(footer string) CommandOpt {
	return func(cat *wargcore.Cmd) {
		cat.Footer = footer
	}
}

// CmdHelpLong adds an Help string to the command - useful from a help function
func CmdHelpLong(helpLong string) CommandOpt {
	return func(cat *wargcore.Cmd) {
		cat.HelpLong = helpLong
	}
}
