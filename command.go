package warg

import (
	"errors"
	"fmt"
	"strings"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/value"
)

// A CmdOpt customizes a Command
type CmdOpt func(*Cmd)

// UnimplementedCmd is a command action that simply returns an error.
// Useful for prototyping
func UnimplementedCmd(_ CmdContext) error {
	return errors.New("TODO: implement this command")
}

// NewCmd builds a Cmd
func NewCmd(helpShort string, action Action, opts ...CmdOpt) Cmd {
	command := Cmd{
		HelpShort:   helpShort,
		Action:      action,
		Completions: DefaultCmdCompletions,
		Flags:       make(FlagMap),
		Footer:      "",
		HelpLong:    "",
	}
	for _, opt := range opts {
		opt(&command)
	}
	return command
}

func DefaultCmdCompletions(cmdCtx CmdContext) (*completion.Candidates, error) {
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

func CmdCompletions(CompletionsFunc CompletionsFunc) CmdOpt {
	return func(flag *Cmd) {
		flag.Completions = CompletionsFunc
	}
}

// CmdFlag adds an existing flag to a Command. It panics if a flag with the same name exists
func CmdFlag(name string, value Flag) CmdOpt {
	return func(com *Cmd) {
		com.Flags.AddFlag(name, value)
	}
}

// CmdFlagMap adds existing flags to a Command. It panics if a flag with the same name exists
func CmdFlagMap(flagMap FlagMap) CmdOpt {
	return func(com *Cmd) {
		com.Flags.AddFlags(flagMap)
	}
}

// NewCmdFlag builds a flag and adds it to a Command. It panics if a flag with the same name exists
func NewCmdFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...FlagOpt) CmdOpt {
	return CmdFlag(name, NewFlag(helpShort, empty, opts...))
}

// CmdFooter adds an Help string to the command - useful from a help function
func CmdFooter(footer string) CmdOpt {
	return func(cat *Cmd) {
		cat.Footer = footer
	}
}

// CmdHelpLong adds an Help string to the command - useful from a help function
func CmdHelpLong(helpLong string) CmdOpt {
	return func(cat *Cmd) {
		cat.HelpLong = helpLong
	}
}
