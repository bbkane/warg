package cli

import (
	"context"
	"fmt"
	"os"

	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value"
)

// -- moved from app_parse_cli.go

// ParseOpts allows overriding the default inputs to the Parse function. Useful for tests. Create it using the [go.bbkane.com/warg/parseopt] package.
type ParseOpts struct {
	Args []string

	// Context for unstructured data. Useful for setting up mocks for tests (i.e., pass in in memory database and use it if it's here in the context)
	Context context.Context

	LookupEnv LookupEnv

	// Stderr will be passed to command.Context for user commands to print to.
	// This file is never closed by warg, so if setting to something other than stderr/stdout,
	// remember to close the file after running the command.
	// Useful for saving output for tests. Defaults to os.Stderr if not passed
	Stderr *os.File

	// Stdout will be passed to command.Context for user commands to print to.
	// This file is never closed by warg, so if setting to something other than stderr/stdout,
	// remember to close the file after running the command.
	// Useful for saving output for tests. Defaults to os.Stdout if not passed
	Stdout *os.File
}

type ParseOpt func(*ParseOpts)

func NewParseOpts(opts ...ParseOpt) ParseOpts {
	parseOptHolder := ParseOpts{
		Context:   context.Background(),
		Args:      os.Args,
		LookupEnv: os.LookupEnv,
		Stderr:    os.Stderr,
		Stdout:    os.Stdout,
	}

	for _, opt := range opts {
		opt(&parseOptHolder)
	}

	return parseOptHolder
}

// ParseResult holds the result of parsing the command line.
type ParseResult struct {
	Context Context
	// Action holds the passed command's action to execute.
	Action Action
}

// -- FlagValueMap

type FlagValueMap map[string]value.Value

func (m FlagValueMap) ToPassedFlags() PassedFlags {
	pf := make(PassedFlags)
	for name, v := range m {
		if v.UpdatedBy() != value.UpdatedByUnset {
			pf[string(name)] = v.Get()
		}
	}
	return pf
}

// -- ParseState

type ExpectingArg string

const (
	ExpectingArg_SectionOrCommand ExpectingArg = "ExpectingArg_SectionOrCommand"
	ExpectingArg_FlagNameOrEnd    ExpectingArg = "ExpectingArg_FlagNameOrEnd"
	ExpectingArg_FlagValue        ExpectingArg = "ExpectingArg_FlagValue"
)

// -- unsetFlagNameSet

type unsetFlagNameSet map[string]struct{}

func (u unsetFlagNameSet) Add(name string) {
	u[name] = struct{}{}
}

func (u unsetFlagNameSet) Delete(name string) {
	delete(u, name)
}

func (u unsetFlagNameSet) Contains(name string) bool {
	_, exists := u[name]
	return exists
}

type ParseState struct {
	SectionPath    []string
	CurrentSection *Section

	CurrentCommandName string
	CurrentCommand     *Command

	CurrentFlagName string
	CurrentFlag     *Flag
	FlagValues      FlagValueMap
	UnsetFlagNames  unsetFlagNameSet

	HelpPassed   bool
	ExpectingArg ExpectingArg
}

func (a *App) parseArgs(args []string) (ParseState, error) {
	pr := ParseState{
		SectionPath:    nil,
		CurrentSection: &a.RootSection,

		CurrentCommandName: "",
		CurrentCommand:     nil,

		CurrentFlagName: "",
		CurrentFlag:     nil,
		FlagValues:      make(FlagValueMap),
		UnsetFlagNames:  make(unsetFlagNameSet),

		HelpPassed:   false,
		ExpectingArg: ExpectingArg_SectionOrCommand,
	}

	aliasToFlagName := make(map[string]string)
	for flagName, fl := range a.GlobalFlags {
		if fl.Alias != "" {
			aliasToFlagName[string(fl.Alias)] = flagName
		}
	}

	// fill the FlagValues map with empty values from the app
	for flagName := range a.GlobalFlags {
		val := a.GlobalFlags[flagName].EmptyValueConstructor()
		pr.FlagValues[flagName] = val
	}

	for i, arg := range args {

		// --help <helptype> or --help must be the last thing passed and can appear at any state we aren't expecting a flag value
		if i >= len(args)-2 &&
			arg != "" && // just in case there's not help flag alias
			(arg == a.HelpFlagName || arg == a.GlobalFlags[a.HelpFlagName].Alias) &&
			pr.ExpectingArg != ExpectingArg_FlagValue {

			pr.HelpPassed = true
			// set the value of --help if an arg was passed, otherwise let it resolve with the rest of them...
			if i == len(args)-2 {
				err := pr.FlagValues[a.HelpFlagName].Update(args[i+1], value.UpdatedByFlag)
				if err != nil {
					return pr, fmt.Errorf("error updating help flag: %w", err)
				}
			}

			return pr, nil
		}

		switch pr.ExpectingArg {
		case ExpectingArg_SectionOrCommand:
			if childSection, exists := pr.CurrentSection.Sections[string(arg)]; exists {
				pr.CurrentSection = &childSection
				pr.SectionPath = append(pr.SectionPath, arg)
			} else if childCommand, exists := pr.CurrentSection.Commands[string(arg)]; exists {
				pr.CurrentCommand = &childCommand
				pr.CurrentCommandName = string(arg)

				// fill the FlagValues map with empty values from the command
				// All names in (command flag names, command flag aliases, global flag names, global flag aliases)
				// should be unique because app.Validate should have caught any conflicts
				for flagName, f := range pr.CurrentCommand.Flags {
					pr.FlagValues[flagName] = f.EmptyValueConstructor()

					if f.Alias != "" {
						aliasToFlagName[string(f.Alias)] = flagName
					}

				}
				pr.ExpectingArg = ExpectingArg_FlagNameOrEnd
			} else {
				return pr, fmt.Errorf("expecting section or command, got %s", arg)
			}

		case ExpectingArg_FlagNameOrEnd:
			flagName := string(arg)
			if actualFlagName, exists := aliasToFlagName[flagName]; exists {
				flagName = actualFlagName
			}
			fl := findFlag(flagName, a.GlobalFlags, pr.CurrentCommand.Flags)
			if fl == nil {
				// return pr, fmt.Errorf("expecting command flag name %v or app flag name %v, got %s", pr.CurrentCommand.ChildrenNames(), a.GlobalFlags.SortedNames(), arg)
				return pr, fmt.Errorf("expecting flag name, got %s", arg)
			}
			pr.CurrentFlagName = flagName
			pr.CurrentFlag = fl
			pr.ExpectingArg = ExpectingArg_FlagValue

		case ExpectingArg_FlagValue:
			// TODO: unset the flag if UnsetSentinel is passed. Search though global flags and command flags, reset the value to unset sentinal and store in the parseResult that it was unset so calls to resolveFlags won't set it...
			if arg == pr.CurrentFlag.UnsetSentinel {
				pr.FlagValues[pr.CurrentFlagName] = pr.CurrentFlag.EmptyValueConstructor()
				pr.UnsetFlagNames.Add(pr.CurrentFlagName)
			} else {
				err := pr.FlagValues[pr.CurrentFlagName].Update(arg, value.UpdatedByFlag)
				if err != nil {
					return pr, err
				}
				pr.UnsetFlagNames.Delete(pr.CurrentFlagName)
			}
			pr.ExpectingArg = ExpectingArg_FlagNameOrEnd

		default:
			panic("unexpected state: " + pr.ExpectingArg)
		}
	}
	return pr, nil
}

func findFlag(flagName string, globalFlags FlagMap, currentCommandFlags FlagMap) *Flag {
	if fl, exists := globalFlags[flagName]; exists {
		return &fl
	}
	if fl, exists := currentCommandFlags[flagName]; exists {
		return &fl
	}
	return nil
}

func resolveFlag2(
	flagName string,
	fl Flag,
	flagValues FlagValueMap, // this gets updated - all other params are readonly
	configReader config.Reader,
	lookupEnv LookupEnv,
	unsetFlagNames unsetFlagNameSet,
) error {

	// don't update if its been explicitly unset or already set
	if unsetFlagNames.Contains(flagName) || flagValues[flagName].UpdatedBy() != value.UpdatedByUnset {
		return nil
	}

	// config
	if fl.ConfigPath != "" && configReader != nil {
		fpr, err := configReader.Search(fl.ConfigPath)
		if err != nil {
			return err
		}
		if fpr != nil {
			if !fpr.IsAggregated {
				err := flagValues[flagName].ReplaceFromInterface(fpr.IFace, value.UpdatedByConfig)
				if err != nil {
					return fmt.Errorf(
						"could not replace container type value:\nval:\n%#v\nreplacement:\n%#v\nerr: %w",
						flagValues[flagName],
						fpr.IFace,
						err,
					)
				}
				return nil
			} else {
				v, ok := flagValues[flagName].(value.SliceValue)
				if !ok {
					return fmt.Errorf("could not update scalar value with aggregated value from config: name: %v, configPath: %v", flagName, fl.ConfigPath)
				}
				under, ok := fpr.IFace.([]interface{})
				if !ok {
					return fmt.Errorf("expected []interface{}, got: %#v", under)
				}
				for _, e := range under {
					err := v.AppendFromInterface(e, value.UpdatedByConfig)
					if err != nil {
						return fmt.Errorf("could not update container type value: err: %w", err)
					}
				}
				flagValues[flagName] = v
				return nil
			}
		}
	}

	// envvar
	for _, e := range fl.EnvVars {
		val, exists := lookupEnv(e)
		if exists {
			err := flagValues[flagName].Update(val, value.UpdatedByEnvVar)
			if err != nil {
				return fmt.Errorf("error updating flag %v from envvar %v: %w", flagName, val, err)
			}
			// Use first env var found
			return nil
		}
	}

	// default
	if flagValues[flagName].HasDefault() {
		err := flagValues[flagName].ReplaceFromDefault(value.UpdatedByDefault)
		if err != nil {
			return fmt.Errorf("error updating flag %v from default: %w", flagName, err)
		}
		return nil
	}
	return nil
}

// resolveFlags resolves the config flag first, and then uses its values to resolve the rest of the flags.
func (a *App) resolveFlags(currentCommand *Command, flagValues FlagValueMap, lookupEnv LookupEnv, unsetFlagNames unsetFlagNameSet) error {
	// resolve config flag first and try to get a reader
	var configReader config.Reader
	if a.ConfigFlagName != "" {
		err := resolveFlag2(
			a.ConfigFlagName, a.GlobalFlags[a.ConfigFlagName], flagValues, nil, lookupEnv, unsetFlagNames)
		if err != nil {
			return fmt.Errorf("resolveFlag error for flag %s: %w", a.ConfigFlagName, err)
		}
		if flagValues[a.ConfigFlagName].UpdatedBy() != value.UpdatedByUnset {
			configPath := flagValues[a.ConfigFlagName].Get().(path.Path)
			configPathStr, err := configPath.Expand()
			if err != nil {
				return fmt.Errorf("error expanding config path ( %s ) : %w", configPath, err)
			}
			configReader, err = a.NewConfigReader(configPathStr)
			if err != nil {
				return fmt.Errorf("error reading config path ( %s ) : %w", configPath, err)
			}

		}
	}

	// resolve app global flags
	for flagName, fl := range a.GlobalFlags {
		err := resolveFlag2(flagName, fl, flagValues, configReader, lookupEnv, unsetFlagNames)
		if err != nil {
			return fmt.Errorf("resolveFlag error for global flag %s: %w", flagName, err)
		}
	}

	// resolve current command flags
	if currentCommand != nil { // can be nil in the case of --help
		for flagName, fl := range currentCommand.Flags {
			err := resolveFlag2(flagName, fl, flagValues, configReader, lookupEnv, unsetFlagNames)
			if err != nil {
				return fmt.Errorf("resolveFlag error for command flag %s: %w", flagName, err)
			}
		}
	}

	return nil
}

func (app *App) Parse(opts ...ParseOpt) (*ParseResult, error) {

	parseOpts := NewParseOpts(opts...)

	parseState, err := app.parseArgs(parseOpts.Args[1:]) // TODO: make callers do [:1]
	if err != nil {
		return nil, fmt.Errorf("Parse args error: %w", err)
	}

	// --help means we don't need to do a lot of error checking
	if parseState.HelpPassed || parseState.ExpectingArg == ExpectingArg_SectionOrCommand {
		err = app.resolveFlags(parseState.CurrentCommand, parseState.FlagValues, parseOpts.LookupEnv, parseState.UnsetFlagNames)
		if err != nil {
			return nil, err
		}

		helpType := parseState.FlagValues[app.HelpFlagName].Get().(string)
		command := app.HelpCommands[helpType]
		pr := ParseResult{
			Context: Context{
				App:        app,
				Context:    parseOpts.Context,
				Flags:      parseState.FlagValues.ToPassedFlags(),
				ParseState: &parseState,
				Stderr:     parseOpts.Stderr,
				Stdout:     parseOpts.Stdout,
			},
			Action: command.Action,
		}
		return &pr, nil
	}

	// ok, we're running a real command, let's do the error checking
	if parseState.ExpectingArg != ExpectingArg_FlagNameOrEnd {
		return nil, fmt.Errorf("unexpected parse state: %s", parseState.ExpectingArg)
	}

	err = app.resolveFlags(parseState.CurrentCommand, parseState.FlagValues, parseOpts.LookupEnv, parseState.UnsetFlagNames)
	if err != nil {
		return nil, err
	}

	missingRequiredFlags := []string{}
	for flagName, flag := range app.GlobalFlags {
		if flag.Required && parseState.FlagValues[flagName].UpdatedBy() == value.UpdatedByUnset {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	for flagName, flag := range parseState.CurrentCommand.Flags {
		if flag.Required && parseState.FlagValues[flagName].UpdatedBy() == value.UpdatedByUnset {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	if len(missingRequiredFlags) > 0 {
		return nil, fmt.Errorf("missing but required flags: %s", missingRequiredFlags)
	}

	pr := ParseResult{
		Context: Context{
			App:        app,
			Context:    parseOpts.Context,
			Flags:      parseState.FlagValues.ToPassedFlags(),
			ParseState: &parseState,
			Stderr:     parseOpts.Stderr,
			Stdout:     parseOpts.Stdout,
		},
		Action: parseState.CurrentCommand.Action,
	}
	return &pr, nil
}
