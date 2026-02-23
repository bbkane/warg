package warg

import (
	"errors"
	"fmt"
	"os"

	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/metadata"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/set"
	"go.bbkane.com/warg/value"
)

// -- moved from app_parse_cli.go

// ParseOpts allows overriding the default inputs to the Parse function. Useful for tests. Create it using the [go.bbkane.com/warg/parseopt] package.
type ParseOpts struct {
	// Args []string

	// ParseMetadata for unstructured data. Useful for setting up mocks for tests (i.e., pass in in memory database and use it if it's here in the context)
	ParseMetadata metadata.Metadata

	LookupEnv LookupEnv

	// Stderr will be passed to [CmdContext] for user commands to print to.
	// This file is never closed by warg, so if setting to something other than stderr/stdout,
	// remember to close the file after running the command.
	// Useful for saving output for tests. Defaults to os.Stderr if not passed
	Stderr *os.File

	// Stdin will be passed to [CmdContext] for user commands to read from.
	// This file is never closed by warg, so if setting to something other than stdin/stdout,
	// remember to close the file after running the command.
	// Useful for saving input for tests. Defaults to os.Stdin if not passed
	Stdin *os.File

	// Stdout will be passed to [CmdContext] for user commands to print to.
	// This file is never closed by warg, so if setting to something other than stderr/stdout,
	// remember to close the file after running the command.
	// Useful for saving output for tests. Defaults to os.Stdout if not passed
	Stdout *os.File
}

type ParseOpt func(*ParseOpts)

func NewParseOpts(opts ...ParseOpt) ParseOpts {
	parseOptHolder := ParseOpts{
		ParseMetadata: metadata.Empty(),
		LookupEnv:     os.LookupEnv,
		Stderr:        os.Stderr,
		Stdin:         os.Stdin,
		Stdout:        os.Stdout,
	}

	for _, opt := range opts {
		opt(&parseOptHolder)
	}

	return parseOptHolder
}

// ParseResult holds the result of parsing the command line.
type ParseResult struct {
	Context CmdContext
	// Action holds the passed command's action to execute.
	Action Action
}

// -- FlagValueMap

// ValueMap holds flag values. If produced as part of [ParseState], it will be fully resolved (i.e., config/env/defaults applied if possible).
type ValueMap map[string]value.Value

func (m ValueMap) ToPassedFlags() PassedFlags {
	pf := make(PassedFlags)
	for name, v := range m {
		if v.UpdatedBy() != value.UpdatedByUnset {
			pf[string(name)] = v.Get()
		}
	}
	return pf
}

// IsSet returns true if the flag with the given name has been set to a non-empty value (i.e., not its empty constructor value). Assumes the flag exists in the map.
func (m ValueMap) IsSet(flagName string) bool {
	return m[flagName].UpdatedBy() != value.UpdatedByUnset
}

// -- ParseState

// ParseArgState represents the current "thing" we want from the args. It transitions as we parse each incoming argument and match it to the expected application structure
type ParseArgState string

const (
	ParseArgState_WantSectionOrCmd  ParseArgState = "ParseArgState_WantSectionOrCmd"
	ParseArgState_WantFlagNameOrEnd ParseArgState = "ParseArgState_WantFlagNameOrEnd"
	ParseArgState_WantFlagValue     ParseArgState = "ParseArgState_WantFlagValue"
)

// ParseState holds the current state of parsing the command line arguments, as well as fully resolving all flag values (including from config/env/defaults).
//
// See ParseArgState for which fields are valid:
//
//   - [ParseArgState_WantSectionOrCmd]: only CurrentSection, SectionPath are valid
//   - [ParseArgState_WantFlagNameOrEnd], [ParseArgState_WantFlagValue]: all fields valid!
type ParseState struct {
	ParseArgState ParseArgState

	SectionPath    []string
	CurrentSection *Section

	CurrentCmdName          string
	CurrentCmd              *Cmd
	CurrentCmdForwardedArgs []string

	CurrentFlagName string
	CurrentFlag     *Flag

	// FlagValues holds all flag values, including global and command flags, keyed by flag name. It is always non-nil, and is filled with empty values for global flags at the start of parsing, and for command flags when a command is selected (state != [ParseArgState_WantSectionOrCmd]). These flags are updated with non-empty values as flags are resolved.
	FlagValues     ValueMap
	UnsetFlagNames set.Set[string]

	HelpPassed bool
}

// parseArgs parses the args into a ParseState. It does not resolve flag values from config/env/defaults, just from the command line. It should always be followed by a call to before returning from a public API so callers see fully resolved values.
func (app *App) parseArgs(args []string) (ParseState, error) {
	pr := ParseState{
		ParseArgState: ParseArgState_WantSectionOrCmd,

		SectionPath:    nil,
		CurrentSection: &app.RootSection,

		CurrentCmdName:          "",
		CurrentCmd:              nil,
		CurrentCmdForwardedArgs: nil,

		CurrentFlagName: "",
		CurrentFlag:     nil,
		FlagValues:      make(ValueMap),
		UnsetFlagNames:  set.New[string](),

		HelpPassed: false,
	}

	aliasToFlagName := make(map[string]string)
	for flagName, fl := range app.GlobalFlags {
		if fl.Alias != "" {
			aliasToFlagName[string(fl.Alias)] = flagName
		}
	}

	// fill the FlagValues map with empty values from the app
	for flagName := range app.GlobalFlags {
		val := app.GlobalFlags[flagName].EmptyValueConstructor()
		pr.FlagValues[flagName] = val
	}

	for i, arg := range args {

		// --help <helptype> or --help must be the last thing passed and can appear at any state we aren't expecting a flag value
		if i >= len(args)-2 &&
			arg != "" && // just in case there's not help flag alias
			(arg == app.HelpFlagName || arg == app.GlobalFlags[app.HelpFlagName].Alias) &&
			pr.ParseArgState != ParseArgState_WantFlagValue {

			pr.HelpPassed = true
			// set the value of --help if an arg was passed, otherwise let it resolve with the rest of them...
			if i == len(args)-2 {
				err := pr.FlagValues[app.HelpFlagName].Update(args[i+1], value.UpdatedByFlag)
				if err != nil {
					return pr, fmt.Errorf("error updating help flag: %w", err)
				}
			}

			return pr, nil
		}

		switch pr.ParseArgState {
		case ParseArgState_WantSectionOrCmd:
			if childSection, exists := pr.CurrentSection.Sections[string(arg)]; exists {
				pr.CurrentSection = &childSection
				pr.SectionPath = append(pr.SectionPath, arg)
			} else if childCommand, exists := pr.CurrentSection.Cmds[string(arg)]; exists {
				pr.CurrentCmd = &childCommand
				pr.CurrentCmdName = string(arg)

				// fill the FlagValues map with empty values from the command
				// All names in (command flag names, command flag aliases, global flag names, global flag aliases)
				// should be unique because app.Validate should have caught any conflicts
				for flagName, f := range pr.CurrentCmd.Flags {
					pr.FlagValues[flagName] = f.EmptyValueConstructor()

					if f.Alias != "" {
						aliasToFlagName[string(f.Alias)] = flagName
					}

				}
				pr.ParseArgState = ParseArgState_WantFlagNameOrEnd
			} else {
				return pr, fmt.Errorf("expecting section or command, got %s", arg)
			}

		case ParseArgState_WantFlagNameOrEnd:
			flagName := arg

			// check if we need to handle forwarded args
			if flagName == "--" && pr.CurrentCmd.AllowForwardedArgs {
				if i >= len(args)-1 {
					return pr, errors.New("expecting forwarded args after --")
				}
				// all remaining args are forwarded args
				pr.CurrentCmdForwardedArgs = append(pr.CurrentCmdForwardedArgs, args[i+1:]...)
				return pr, nil
			}

			if actualFlagName, exists := aliasToFlagName[flagName]; exists {
				flagName = actualFlagName
			}
			fl := findFlag(flagName, app.GlobalFlags, pr.CurrentCmd.Flags)
			if fl == nil {
				return pr, fmt.Errorf("expecting flag name, got %s", arg)
			}
			pr.CurrentFlagName = flagName
			pr.CurrentFlag = fl
			pr.ParseArgState = ParseArgState_WantFlagValue

		case ParseArgState_WantFlagValue:
			// if the flag has an unset sentinel and the user passed it, unset the flag
			// NOTE: UnsetSentinel must be a pointer to a string, because sometimes the user may pass an empty string
			if pr.CurrentFlag.UnsetSentinel != nil && arg == *pr.CurrentFlag.UnsetSentinel {
				pr.FlagValues[pr.CurrentFlagName] = pr.CurrentFlag.EmptyValueConstructor()
				pr.UnsetFlagNames.Add(pr.CurrentFlagName)
			} else {
				err := pr.FlagValues[pr.CurrentFlagName].Update(arg, value.UpdatedByFlag)
				if err != nil {
					return pr, err
				}
				pr.UnsetFlagNames.Delete(pr.CurrentFlagName)
			}
			pr.ParseArgState = ParseArgState_WantFlagNameOrEnd

		default:
			panic("unexpected state: " + pr.ParseArgState)
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

func resolveFlag(
	flagName string,
	fl Flag,
	flagValues ValueMap, // this gets updated - all other params are readonly
	configReader config.Reader,
	lookupEnv LookupEnv,
	unsetFlagNames set.Set[string],
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
func (app *App) resolveFlags(currentCmd *Cmd, flagValues ValueMap, lookupEnv LookupEnv, unsetFlagNames set.Set[string]) error {
	// resolve config flag first and try to get a reader
	var configReader config.Reader
	if app.ConfigFlagName != "" {
		err := resolveFlag(
			app.ConfigFlagName, app.GlobalFlags[app.ConfigFlagName], flagValues, nil, lookupEnv, unsetFlagNames)
		if err != nil {
			return fmt.Errorf("resolveFlag error for flag %s: %w", app.ConfigFlagName, err)
		}
		if flagValues[app.ConfigFlagName].UpdatedBy() != value.UpdatedByUnset {
			configPath := flagValues[app.ConfigFlagName].Get().(path.Path)
			configPathStr, err := configPath.Expand()
			if err != nil {
				return fmt.Errorf("error expanding config path ( %s ) : %w", configPath, err)
			}
			configReader, err = app.NewConfigReader(configPathStr)
			if err != nil {
				return fmt.Errorf("error reading config path ( %s ) : %w", configPath, err)
			}

		}
	}

	// resolve app global flags
	for flagName, fl := range app.GlobalFlags {
		err := resolveFlag(flagName, fl, flagValues, configReader, lookupEnv, unsetFlagNames)
		if err != nil {
			return fmt.Errorf("resolveFlag error for global flag %s: %w", flagName, err)
		}
	}

	// resolve current command flags
	if currentCmd != nil { // can be nil in the case of --help
		for flagName, fl := range currentCmd.Flags {
			err := resolveFlag(flagName, fl, flagValues, configReader, lookupEnv, unsetFlagNames)
			if err != nil {
				return fmt.Errorf("resolveFlag error for command flag %s: %w", flagName, err)
			}
		}
	}

	return nil
}

// Parse parses command line arguments, environment variables, and configuration files to produce a [ParseResult]. expects ParseOpts.Args to be like os.Args (i.e., first arg is app name). It returns an error if parsing fails or required flags are missing.
func (app *App) Parse(args []string, opts ...ParseOpt) (*ParseResult, error) {

	parseOpts := NewParseOpts(opts...)

	parseState, err := app.parseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("Parse args error: %w", err)
	}

	// --help means we don't need to do a lot of error checking
	if parseState.HelpPassed || parseState.ParseArgState == ParseArgState_WantSectionOrCmd {
		err = app.resolveFlags(parseState.CurrentCmd, parseState.FlagValues, parseOpts.LookupEnv, parseState.UnsetFlagNames)
		if err != nil {
			return nil, err
		}

		helpType := parseState.FlagValues[app.HelpFlagName].Get().(string)
		command := app.HelpCmds[helpType]
		pr := ParseResult{
			Context: CmdContext{
				App:           app,
				ParseMetadata: parseOpts.ParseMetadata,
				Flags:         parseState.FlagValues.ToPassedFlags(),
				ForwardedArgs: parseState.CurrentCmdForwardedArgs,
				ParseState:    &parseState,
				Stderr:        parseOpts.Stderr,
				Stdin:         parseOpts.Stdin,
				Stdout:        parseOpts.Stdout,
			},
			Action: command.Action,
		}
		return &pr, nil
	}

	// ok, we're running a real command, let's do the error checking
	if parseState.ParseArgState != ParseArgState_WantFlagNameOrEnd {
		return nil, fmt.Errorf("unexpected parse state: %s", parseState.ParseArgState)
	}

	err = app.resolveFlags(parseState.CurrentCmd, parseState.FlagValues, parseOpts.LookupEnv, parseState.UnsetFlagNames)
	if err != nil {
		return nil, err
	}

	missingRequiredFlags := []string{}
	for flagName, flag := range app.GlobalFlags {
		if flag.Required && !parseState.FlagValues.IsSet(flagName) {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	for flagName, flag := range parseState.CurrentCmd.Flags {
		if flag.Required && !parseState.FlagValues.IsSet(flagName) {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	if len(missingRequiredFlags) > 0 {
		return nil, fmt.Errorf("missing but required flags: %s", missingRequiredFlags)
	}

	pr := ParseResult{
		Context: CmdContext{
			App:           app,
			ParseMetadata: parseOpts.ParseMetadata,
			Flags:         parseState.FlagValues.ToPassedFlags(),
			ForwardedArgs: parseState.CurrentCmdForwardedArgs,
			ParseState:    &parseState,
			Stderr:        parseOpts.Stderr,
			Stdin:         parseOpts.Stdin,
			Stdout:        parseOpts.Stdout,
		},
		Action: parseState.CurrentCmd.Action,
	}
	return &pr, nil
}
