package wargcore

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"slices"

	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value"
)

// An App contains your defined sections, commands, and flags
// Create a new App with New()
type App struct {
	// Config
	ConfigFlagName  string
	NewConfigReader config.NewReader

	// Help
	HelpFlagName string
	HelpCommands CommandMap

	GlobalFlags            FlagMap
	Name                   string
	RootSection            Section
	SkipGlobalColorFlag    bool
	SkipCompletionCommands bool
	SkipValidation         bool
	SkipVersionCommand     bool
	Version                string
}

// MustRun runs the app.
// Any flag parsing errors will be printed to stderr and os.Exit(64) (EX_USAGE) will be called.
// Any errors on an Action will be printed to stderr and os.Exit(1) will be called.
func (app *App) MustRun(opts ...ParseOpt) {
	if len(os.Args) >= 3 && os.Args[1] == "--completion-zsh" {
		// app --completion-zsh <args> . Note that <args> must be something, even if it's the empty string

		candidates, err := app.CompletionCandidates(opts...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		completion.ZshCompletionsWrite(os.Stdout, candidates)

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
type LookupEnv func(key string) (string, bool)

// LookupMap loooks up keys from a provided map. Useful to mock os.LookupEnv when parsing
func LookupMap(m map[string]string) LookupEnv {
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
	globalFlags FlagMap,
	comFlags FlagMap,
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
//   - the help flag is the right type
//   - Sections and commands don't start with "-" (needed for parsing)
//   - Flag names and aliases do start with "-" (needed for parsing)
//   - Flag names and aliases don't collide
func (app *App) Validate() error {

	// validate --help flag
	if app.HelpFlagName == "" {
		return fmt.Errorf("HelpFlagName must be set")
	}
	helpFlag, exists := app.GlobalFlags[app.HelpFlagName]
	if !exists {
		return fmt.Errorf("HelpFlagName not found in GlobalFlags: %v", app.HelpFlagName)
	}
	helpFlagValEmpty, ok := helpFlag.EmptyValueConstructor().(value.ScalarValue)
	if !ok {
		return fmt.Errorf("HelpFlagName must be a scalar: %v", app.HelpFlagName)
	}
	if _, ok := helpFlagValEmpty.Get().(string); !ok {
		return fmt.Errorf("HelpFlagName must be a string: %v", app.HelpFlagName)
	}
	if !helpFlagValEmpty.HasDefault() {
		return fmt.Errorf("HelpFlagName must have a default value: %v", app.HelpFlagName)
	}
	if !slices.Equal(helpFlagValEmpty.Choices(), app.HelpCommands.SortedNames()) {
		return fmt.Errorf("HelpFlagName choices must match HelpCommands: %v", app.HelpFlagName)
	}
	if !slices.Contains(helpFlagValEmpty.Choices(), helpFlagValEmpty.DefaultString()) {
		return fmt.Errorf("HelpFlagName default value (%v) must be in choices (%v): %v", helpFlagValEmpty.DefaultString(), helpFlagValEmpty.Choices(), app.HelpFlagName)
	}

	// validate --config flag
	if app.ConfigFlagName != "" {
		if app.NewConfigReader == nil {
			return fmt.Errorf("ConfigFlagName must have a NewConfigReader: %v", app.ConfigFlagName)
		}
		configFlag, exists := app.GlobalFlags[app.ConfigFlagName]
		if !exists {
			return fmt.Errorf("ConfigFlagName not found in GlobalFlags: %v", app.ConfigFlagName)
		}
		configFlagValEmpty, ok := configFlag.EmptyValueConstructor().(value.ScalarValue)
		if !ok {
			return fmt.Errorf("ConfigFlagName must be a scalar: %v", app.ConfigFlagName)
		}
		if _, ok := configFlagValEmpty.Get().(path.Path); !ok {
			return fmt.Errorf("ConfigFlagName must be a path: %v", app.ConfigFlagName)
		}
	}

	// TODO: check that the default value is in the choices and the choices match app help mappings and that the flag is a scalar

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

// CompletionCandidatesFunc is a function that returns completion candidates for a flag. See warg.CompletionCandidates[Type] for convenience functions to make this
type CompletionCandidatesFunc func(Context) (*completion.Candidates, error)

func (a *App) CompletionCandidates(opts ...ParseOpt) (*completion.Candidates, error) {
	parseOpts := NewParseOpts(opts...)

	// parseOpts.Args looks like: <exe> --completion-zsh <args>... <partialOrEmptyString>
	// the partial or empty string is passed to us from the completion script. Empty if the user just typed space and pressed tab, partial if the user pressed tab after typing part of something. zsh will filter that for us
	// so we need to remove the first two args and the last arg
	args := parseOpts.Args[2 : len(parseOpts.Args)-1]

	// I could to a full parse here, but that would be slower and more prone to failure than just parsing the args - we don't need a lot of info to complete section/command names
	parseState, err := a.parseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("unexpected parseArgs err: %w", err)
	}

	// special case if help is passed
	if parseState.HelpPassed {
		// if the value of the flag has been passed, don't suggest anything
		if parseState.FlagValues[a.HelpFlagName].UpdatedBy() == value.UpdatedByFlag {
			return &completion.Candidates{
				Type:   completion.Type_None,
				Values: nil,
			}, nil
		}

		// otherwise suggest the help commands as the values of the help flag
		res := &completion.Candidates{
			Type:   completion.Type_Values,
			Values: []completion.Candidate{},
		}
		for _, name := range a.HelpCommands.SortedNames() {
			res.Values = append(res.Values, completion.Candidate{
				Name:        string(name),
				Description: "",
			})
		}
		return res, nil
	}

	if parseState.ExpectingArg == ExpectingArg_SectionOrCommand {
		s := parseState.CurrentSection
		ret := completion.Candidates{
			Type:   completion.Type_ValuesDescriptions,
			Values: []completion.Candidate{},
		}
		for _, name := range s.Commands.SortedNames() {
			ret.Values = append(ret.Values, completion.Candidate{
				Name:        string(name),
				Description: string(s.Commands[name].HelpShort),
			})
		}
		for _, name := range s.Sections.SortedNames() {
			ret.Values = append(ret.Values, completion.Candidate{
				Name:        string(name),
				Description: string(s.Sections[name].HelpShort),
			})
		}
		ret.Values = append(ret.Values, completion.Candidate{
			Name:        a.HelpFlagName,
			Description: a.GlobalFlags[a.HelpFlagName].HelpShort,
		})
		return &ret, nil
	}

	// Finish the parse!
	err = a.resolveFlags(parseState.CurrentCommand, parseState.FlagValues, parseOpts.LookupEnv, parseState.UnsetFlagNames)
	if err != nil {
		return nil, fmt.Errorf("unexpected resolveFlags err: %w", err)
	}
	cmdContext := Context{
		App:        a,
		Context:    parseOpts.Context,
		Flags:      parseState.FlagValues.ToPassedFlags(),
		ParseState: &parseState,
		Stderr:     parseOpts.Stderr,
		Stdout:     parseOpts.Stdout,
	}

	switch parseState.ExpectingArg {
	case ExpectingArg_FlagNameOrEnd:
		return parseState.CurrentCommand.CompletionCandidates(cmdContext)
	case ExpectingArg_FlagValue:
		return parseState.CurrentFlag.CompletionCandidates(cmdContext)
	case ExpectingArg_SectionOrCommand:
		panic("unreachable state: ExpectingArg_SectionOrCommand")
	default:
		return nil, fmt.Errorf("unexpected ParseState: %v", parseState.ExpectingArg)
	}
}
