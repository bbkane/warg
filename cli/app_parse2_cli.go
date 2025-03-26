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

type ParseOptHolder struct {
	Args []string

	// Context for unstructured data. Useful for setting up mocks for tests (i.e., pass in in memory database and use it if it's here in the context)
	Context context.Context

	LookupFunc LookupFunc

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

type ParseOpt func(*ParseOptHolder)

func AddContext(ctx context.Context) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Context = ctx
	}
}

func OverrideArgs(args []string) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Args = args
	}
}

func OverrideLookupFunc(lookup LookupFunc) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.LookupFunc = lookup
	}
}

func OverrideStderr(stderr *os.File) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Stderr = stderr
	}
}

func OverrideStdout(stdout *os.File) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Stdout = stdout
	}
}

func NewParseOptHolder(opts ...ParseOpt) ParseOptHolder {
	parseOptHolder := ParseOptHolder{
		Context:    nil,
		Args:       nil,
		LookupFunc: nil,
		Stderr:     nil,
		Stdout:     nil,
	}

	for _, opt := range opts {
		opt(&parseOptHolder)
	}

	if parseOptHolder.Args == nil {
		OverrideArgs(os.Args)(&parseOptHolder)
	}

	if parseOptHolder.Context == nil {
		AddContext(context.Background())(&parseOptHolder)
	}

	if parseOptHolder.LookupFunc == nil {
		OverrideLookupFunc(os.LookupEnv)(&parseOptHolder)
	}

	if parseOptHolder.Stderr == nil {
		OverrideStderr(os.Stderr)(&parseOptHolder)
	}

	if parseOptHolder.Stdout == nil {
		OverrideStdout(os.Stdout)(&parseOptHolder)
	}

	return parseOptHolder
}

// ParseResult holds the result of parsing the command line.
type ParseResult struct {
	Context Context
	// Action holds the passed command's action to execute.
	Action Action
}

// -- FlagValue

type FlagValue struct {
	SetBy string
	Value value.Value
}

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

type ParseState string

const (
	Parse_ExpectingSectionOrCommand ParseState = "Parse_ExpectingSectionOrCommand"
	Parse_ExpectingFlagNameOrEnd    ParseState = "Parse_ExpectingFlagNameOrEnd"
	Parse_ExpectingFlagValue        ParseState = "Parse_ExpectingFlagValue"
)

// -- unsetFlagNameSet

type UnsetFlagNameSet map[string]struct{}

func (u UnsetFlagNameSet) Add(name string) {
	u[name] = struct{}{}
}

func (u UnsetFlagNameSet) Delete(name string) {
	delete(u, name)
}

func (u UnsetFlagNameSet) Contains(name string) bool {
	_, exists := u[name]
	return exists
}

type ParseResult2 struct {
	SectionPath    []string
	CurrentSection *SectionT

	CurrentCommandName string
	CurrentCommand     *Command

	CurrentFlagName string
	CurrentFlag     *Flag
	FlagValues      FlagValueMap
	UnsetFlagNames  UnsetFlagNameSet

	HelpPassed bool
	State      ParseState
}

func (a *App) parseArgs(args []string) (ParseResult2, error) {
	pr := ParseResult2{
		SectionPath:    nil,
		CurrentSection: &a.RootSection,

		CurrentCommandName: "",
		CurrentCommand:     nil,

		CurrentFlagName: "",
		CurrentFlag:     nil,
		FlagValues:      make(FlagValueMap),
		UnsetFlagNames:  make(UnsetFlagNameSet),

		HelpPassed: false,
		State:      Parse_ExpectingSectionOrCommand,
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
			pr.State != Parse_ExpectingFlagValue {

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

		switch pr.State {
		case Parse_ExpectingSectionOrCommand:
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
				pr.State = Parse_ExpectingFlagNameOrEnd
			} else {
				return pr, fmt.Errorf("expecting section or command, got %s", arg)
			}

		case Parse_ExpectingFlagNameOrEnd:
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
			pr.State = Parse_ExpectingFlagValue

		case Parse_ExpectingFlagValue:
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
			pr.State = Parse_ExpectingFlagNameOrEnd

		default:
			panic("unexpected state: " + pr.State)
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
	lookupEnv LookupFunc,
	unsetFlagNames UnsetFlagNameSet,
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
func (a *App) resolveFlags(currentCommand *Command, flagValues FlagValueMap, lookupEnv LookupFunc, unsetFlagNames UnsetFlagNameSet) error {
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
			return fmt.Errorf("resolveFlag error for flag %s: %w", flagName, err)
		}
	}

	// resolve current command flags
	if currentCommand != nil { // can be nil in the case of --help
		for flagName, fl := range currentCommand.Flags {
			err := resolveFlag2(flagName, fl, flagValues, configReader, lookupEnv, unsetFlagNames)
			if err != nil {
				return fmt.Errorf("resolveFlag error for flag %s: %w", flagName, err)
			}
		}
	}

	return nil
}

func (app *App) Parse(opts ...ParseOpt) (*ParseResult, error) {

	parseOpts := NewParseOptHolder(opts...)

	// --config flag...
	// original Parse treats it specially
	// Parse2 expects it to be in app.GlobalFlags
	// TODO: rework the config flag handling. I'd prefer everything to be immutable before calling parse
	if app.ConfigFlag != nil {
		app.GlobalFlags[app.ConfigFlagName] = *app.ConfigFlag
	}

	pr2, err := app.parseArgs(parseOpts.Args[1:]) // TODO: make callers do [:1]
	if err != nil {
		return nil, fmt.Errorf("Parse args error: %w", err)
	}

	// --help means we don't need to do a lot of error checking
	if pr2.HelpPassed || pr2.State == Parse_ExpectingSectionOrCommand {
		err = app.resolveFlags(pr2.CurrentCommand, pr2.FlagValues, parseOpts.LookupFunc, pr2.UnsetFlagNames)
		if err != nil {
			return nil, err
		}

		helpType := pr2.FlagValues[app.HelpFlagName].Get().(string)
		command := app.HelpCommands[helpType]
		pr := ParseResult{
			Context: Context{
				App:         app,
				Context:     parseOpts.Context,
				Flags:       pr2.FlagValues.ToPassedFlags(),
				ParseResult: &pr2,
				Stderr:      parseOpts.Stderr,
				Stdout:      parseOpts.Stdout,
			},
			Action: command.Action,
		}
		return &pr, nil
	}

	// ok, we're running a real command, let's do the error checking
	if pr2.State != Parse_ExpectingFlagNameOrEnd {
		return nil, fmt.Errorf("unexpected parse state: %s", pr2.State)
	}

	err = app.resolveFlags(pr2.CurrentCommand, pr2.FlagValues, parseOpts.LookupFunc, pr2.UnsetFlagNames)
	if err != nil {
		return nil, err
	}

	missingRequiredFlags := []string{}
	for flagName, flag := range app.GlobalFlags {
		if flag.Required && pr2.FlagValues[flagName].UpdatedBy() == value.UpdatedByUnset {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	for flagName, flag := range pr2.CurrentCommand.Flags {
		if flag.Required && pr2.FlagValues[flagName].UpdatedBy() == value.UpdatedByUnset {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	if len(missingRequiredFlags) > 0 {
		return nil, fmt.Errorf("missing but required flags: %s", missingRequiredFlags)
	}

	pr := ParseResult{
		Context: Context{
			App:         app,
			Context:     parseOpts.Context,
			Flags:       pr2.FlagValues.ToPassedFlags(),
			ParseResult: &pr2,
			Stderr:      parseOpts.Stderr,
			Stdout:      parseOpts.Stdout,
		},
		Action: pr2.CurrentCommand.Action,
	}
	return &pr, nil
}
