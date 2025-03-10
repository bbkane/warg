// Declaratively create heirarchical command line apps.
package warg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"slices"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/section"
)

// An App contains your defined sections, commands, and flags
// Create a new App with New()
type App struct {
	// Config()
	ConfigFlagName  string
	NewConfigReader config.NewReader
	ConfigFlag      *flag.Flag

	GlobalFlags flag.FlagMap

	// New Help()
	Name         string
	HelpFlagName string
	// Note that this can be ""
	HelpFlagAlias string
	HelpMappings  []help.HelpFlagMapping

	// RootSection holds the good stuff!
	RootSection section.SectionT

	SkipValidation bool

	Version string
}

// MustRun runs the app.
// Any flag parsing errors will be printed to stderr and os.Exit(64) (EX_USAGE) will be called.
// Any errors on an Action will be printed to stderr and os.Exit(1) will be called.
func (app *App) MustRun(opts ...ParseOpt) {
	// TODO: make this better
	if slices.Equal(os.Args, []string{os.Args[0], "--completion-script-zsh"}) {
		// app --completion-script-zsh
		completion.WriteCompletionScriptZsh(os.Stdout, app.Name)
	} else if len(os.Args) >= 3 && os.Args[1] == "--completion-zsh" {
		// app --completion-zsh <args> . Note that <args> must be something, even if it's the empty string

		// chop off the last arg since it's either:
		//  - the empty string (if the user just typed space)
		//  - a partial string (if the user pressed tab after typing part of something)
		toComplete := os.Args[2 : len(os.Args)-1]
		candidates, err := app.CompletionCandidates(toComplete)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// TODO: print errors to stderr maybe? For now just silently fail...
		if err == nil {
			fmt.Println(candidates.Type)
			for _, candidate := range candidates.Values {
				fmt.Println(candidate.Name)
				fmt.Println(candidate.Name + " - " + candidate.Description)
			}
		}
	} else {
		pr, err := app.Parse(opts...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			// https://unix.stackexchange.com/a/254747/185953
			os.Exit(64)
		}
		err = pr.Action(pr.Context)
		if err != nil {
			fmt.Fprintln(pr.Context.Stderr, err)
			os.Exit(1)
		}

	}

}

// Look up keys (meant for environment variable parsing) - fulfillable with os.LookupEnv or warg.LookupMap(map)
type LookupFunc func(key string) (string, bool)

// LookupMap loooks up keys from a provided map. Useful to mock os.LookupEnv when parsing
func LookupMap(m map[string]string) LookupFunc {
	return func(key string) (string, bool) {
		val, exists := m[key]
		return val, exists
	}
}

// validateFlags2 checks that global and command flag names and aliases start with "-" and are unique.
// It does not need to check the following scenarios:
//
//   - global flag names don't collide with global flag names (app will panic when adding the second global flag) - TOOD: ensure there's a test for this
//   - command flag names in the same command don't collide with each other (app will panic when adding the second command flag) TODO: ensure there's a test for this
//   - command flag names/aliases don't collide with command flag names/aliases in other commands (since only one command will be run, this is not a problem)
func validateFlags2(
	globalFlags flag.FlagMap,
	comFlags flag.FlagMap,
) error {
	nameCount := make(map[string]int)
	for name, fl := range globalFlags {
		nameCount[name]++
		if fl.Alias != "" {
			nameCount[fl.Alias]++
		}
	}
	for name, fl := range comFlags {
		nameCount[name]++
		if fl.Alias != "" {
			nameCount[fl.Alias]++
		}
	}
	var errs []error
	for name, count := range nameCount {
		if !strings.HasPrefix(string(name), "-") {
			errs = append(errs, fmt.Errorf("flag and alias names must start with '-': %#v", name))
		}
		if count > 1 {
			errs = append(errs, fmt.Errorf("flag or alias name exists %d times: %v", count, name))
		}
	}
	return errors.Join(errs...)
}

// Validate checks app for creation errors. It checks:
//
// - Sections and commands don't start with "-" (needed for parsing)
//
// - Flag names and aliases do start with "-" (needed for parsing)
//
// - Flag names and aliases don't collide
func (app *App) Validate() error {
	// NOTE: we need to be able to validate before we parse, and we may not know the app name
	// till after prsing so set the root path to "root"
	rootPath := []string{string(app.Name)}
	it := app.RootSection.BreadthFirst(rootPath)

	for it.HasNext() {
		flatSec := it.Next()

		// Sections don't start with "-"
		secName := flatSec.Path[len(flatSec.Path)-1]
		if strings.HasPrefix(string(secName), "-") {
			return fmt.Errorf("section names must not start with '-': %#v", secName)
		}

		// Sections must not be leaf nodes
		if flatSec.Sec.Sections.Empty() && flatSec.Sec.Commands.Empty() {
			return fmt.Errorf("sections must have either child sections or child commands: %#v", secName)
		}

		{
			// child section names should not clash with child command names
			nameCount := make(map[string]int)
			for name := range flatSec.Sec.Commands {
				nameCount[string(name)]++
			}
			for name := range flatSec.Sec.Sections {
				nameCount[string(name)]++
			}
			errs := []error{}
			for name, count := range nameCount {
				if count > 1 {
					errs = append(errs, fmt.Errorf("command and section name clash: %s", name))
				}
			}
			err := errors.Join(errs...)
			if err != nil {
				return fmt.Errorf("name collision: %w", err)
			}
		}

		for name, com := range flatSec.Sec.Commands {

			// Commands must not start wtih "-"
			if strings.HasPrefix(string(name), "-") {
				return fmt.Errorf("command names must not start with '-': %#v", name)
			}

			err := validateFlags2(app.GlobalFlags, com.Flags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *App) CompletionCandidates(args []string) (*completion.CompletionCandidates, error) {
	pr, err := a.parseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("unexpected parseArgs err: %w", err)
	}
	switch pr.State {
	case Parse_ExpectingSectionOrCommand:
		candidates, err := pr.CurrentSection.CompletionCandidates()
		if err != nil {
			return nil, fmt.Errorf("Parse_ExpectingSectionOrCommand CompletionCandidates err: %w", err)
		}
		return &candidates, nil
	case Parse_ExpectingFlagNameOrEnd:
		// TODO: if a scalar flag has been passsed, don't suggest it again
		// TODO: get a better order for the flags. For example, envelope needs to db first (unless it's resolved) so further flags can use that. Add an order or "depends on" param?
		candidates := &completion.CompletionCandidates{
			Type:   completion.CompletionType_ValueDescription,
			Values: []completion.CompletionCandidate{},
		}
		// command flags
		for _, name := range pr.CurrentCommand.Flags.SortedNames() {
			candidates.Values = append(candidates.Values, completion.CompletionCandidate{
				Name:        string(name),
				Description: string(pr.CurrentCommand.Flags[name].HelpShort),
			})
		}
		// global flags
		for _, name := range a.GlobalFlags.SortedNames() {
			candidates.Values = append(candidates.Values, completion.CompletionCandidate{
				Name:        string(name),
				Description: string(a.GlobalFlags[name].HelpShort),
			})
		}
		return candidates, nil
	case Parse_ExpectingFlagValue:
		// TODO: allow flags to look at the values of other flags before offering options.
		// This will require some package "flattening" as ParseResult is defined in App, which also import Flag. So... flag shouldn't import the app code...
		// For now, only suggest the flags choices
		candidates := &completion.CompletionCandidates{
			Type:   completion.CompletionType_ValueDescription,
			Values: []completion.CompletionCandidate{},
		}
		// pr.FlagValues is always filled with at least the empty values
		for _, name := range pr.FlagValues[pr.CurrentFlagName].Choices() {
			candidates.Values = append(candidates.Values, completion.CompletionCandidate{
				Name:        name,
				Description: "NO DESCRIPTION",
			})
		}
		return candidates, nil
	default:
		return nil, fmt.Errorf("unexpected ParseState: %v", pr.State)
	}
}
