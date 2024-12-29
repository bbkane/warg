package warg

import (
	"fmt"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

// -- FlagValue

type FlagValue struct {
	SetBy string
	Value value.Value
}

type FlagValueMap map[flag.Name]value.Value

func (m FlagValueMap) ToPassedFlags() command.PassedFlags {
	pf := make(command.PassedFlags)
	for name, v := range m {
		if v.UpdatedBy() != value.UpdatedByUnset {
			pf[string(name)] = v
		}
	}
	return pf
}

type ParseResult2 struct {
	SectionPath    []string
	CurrentSection *section.SectionT

	CurrentCommandName command.Name
	CurrentCommand     *command.Command

	CurrentFlagName flag.Name
	CurrentFlag     *flag.Flag

	FlagValues FlagValueMap
	State      ParseState
	HelpPassed bool
}

type ParseState string

const (
	Parse_ExpectingSectionOrCommand ParseState = "Parse_ExpectingSectionOrCommand"
	Parse_ExpectingFlagNameOrEnd    ParseState = "Parse_ExpectingFlagNameOrEnd"
	Parse_ExpectingFlagValue        ParseState = "Parse_ExpectingFlagValue"
)

func (a *App) parseArgs(args []string) (ParseResult2, error) {
	pr := ParseResult2{
		SectionPath:    []string{},
		CurrentSection: &a.rootSection,

		CurrentCommandName: "",
		CurrentCommand:     nil,

		CurrentFlagName: "",
		CurrentFlag:     nil,
		FlagValues:      make(FlagValueMap),

		HelpPassed: false,

		State: Parse_ExpectingSectionOrCommand,
	}

	// fill the FlagValues map with empty values from the app
	for flagName := range a.globalFlags {
		val := a.globalFlags[flagName].EmptyValueConstructor()
		pr.FlagValues[flagName] = val
	}

	for i, arg := range args {

		// --help <helptype> or --help must be the last thing passed and can appear at any state we aren't expecting a flag value
		if i >= len(args)-2 &&
			flag.Name(arg) == a.helpFlagName &&
			pr.State != Parse_ExpectingFlagValue {

			pr.HelpPassed = true
			// set the value of --help if an arg was passed, otherwise let it resolve with the rest of them...
			if i == len(args)-2 {
				err := pr.FlagValues[a.helpFlagName].Update(args[i+1], value.UpdatedByFlag)
				if err != nil {
					return pr, fmt.Errorf("error updating help flag: %w", err)
				}
			}

			return pr, nil
		}

		switch pr.State {
		case Parse_ExpectingSectionOrCommand:
			if childSection, exists := pr.CurrentSection.Sections[section.Name(arg)]; exists {
				pr.CurrentSection = &childSection
				pr.SectionPath = append(pr.SectionPath, arg)
			} else if childCommand, exists := pr.CurrentSection.Commands[command.Name(arg)]; exists {
				pr.CurrentCommand = &childCommand
				pr.CurrentCommandName = command.Name(arg)

				for flagName, f := range pr.CurrentCommand.Flags {
					_, exists := pr.FlagValues[flagName]
					if exists {
						// NOTE: move this check to app construction
						panic("app flags and command flags cannot share a name: " + flagName)
					}
					pr.FlagValues[flagName] = f.EmptyValueConstructor()
				}

				pr.State = Parse_ExpectingFlagNameOrEnd
			} else {
				return pr, fmt.Errorf("expecting section or command, got %s", arg)
			}

		case Parse_ExpectingFlagNameOrEnd:
			// TODO: handle aliases of flags
			if flagFromArg, exists := a.globalFlags[flag.Name(arg)]; exists {
				pr.CurrentFlagName = flag.Name(arg)
				pr.CurrentFlag = &flagFromArg
				pr.State = Parse_ExpectingFlagValue
			} else if flagFromArg, exists := pr.CurrentCommand.Flags[flag.Name(arg)]; exists {
				pr.CurrentFlagName = flag.Name(arg)
				pr.CurrentFlag = &flagFromArg
				pr.State = Parse_ExpectingFlagValue
			} else {
				// return pr, fmt.Errorf("expecting command flag name %v or app flag name %v, got %s", pr.CurrentCommand.ChildrenNames(), a.GlobalFlags.SortedNames(), arg)
				return pr, fmt.Errorf("expecting flag name, got %s", arg)
			}

		case Parse_ExpectingFlagValue:
			err := pr.FlagValues[pr.CurrentFlagName].Update(arg, value.UpdatedByFlag)
			if err != nil {
				return pr, err
			}
			pr.State = Parse_ExpectingFlagNameOrEnd

		default:
			panic("unexpected state: " + pr.State)
		}
	}
	return pr, nil
}

func resolveFlag2(
	flagName flag.Name,
	fl flag.Flag,
	flagValues FlagValueMap, // this gets updated - all other params are readonly
	configReader config.Reader,
	lookupEnv LookupFunc,
) (bool, error) {

	// maybe it's set by args already
	isSet := flagValues[flagName].UpdatedBy() != value.UpdatedByUnset
	if isSet {
		return true, nil
	}

	// config
	if fl.ConfigPath != "" && configReader != nil {
		fpr, err := configReader.Search(fl.ConfigPath)
		if err != nil {
			return false, err
		}
		if fpr != nil {
			if !fpr.IsAggregated {
				err := flagValues[flagName].ReplaceFromInterface(fpr.IFace, value.UpdatedByConfig)
				if err != nil {
					return false, fmt.Errorf(
						"could not replace container type value:\nval:\n%#v\nreplacement:\n%#v\nerr: %w",
						flagValues[flagName],
						fpr.IFace,
						err,
					)
				}
			} else {
				v, ok := flagValues[flagName].(value.SliceValue)
				if !ok {
					return false, fmt.Errorf("could not update scalar value with aggregated value from config: name: %v, configPath: %v", flagName, fl.ConfigPath)
				}
				under, ok := fpr.IFace.([]interface{})
				if !ok {
					return false, fmt.Errorf("expected []interface{}, got: %#v", under)
				}
				for _, e := range under {
					err := v.AppendFromInterface(e, value.UpdatedByConfig)
					if err != nil {
						return false, fmt.Errorf("could not update container type value: err: %w", err)
					}
				}
				flagValues[flagName] = v
			}
		}
		return true, nil
	}

	// envvar
	for _, e := range fl.EnvVars {
		val, exists := lookupEnv(e)
		if exists {
			err := flagValues[flagName].Update(val, value.UpdatedByEnvVar)
			if err != nil {
				return false, fmt.Errorf("error updating flag %v from envvar %v: %w", flagName, val, err)
			}
			// Use first env var found
			return true, nil
		}
	}

	// default
	if fl.Value.HasDefault() {
		flagValues[flagName].ReplaceFromDefault(value.UpdatedByDefault)
		return true, nil
	}
	return false, nil
}

func (a *App) resolveFlags(currentCommand *command.Command, flagValues FlagValueMap, lookupEnv LookupFunc) error {
	// resolve config flag first and try to get a reader
	var configReader config.Reader
	if a.configFlagName != "" {
		resolved, err := resolveFlag2(
			a.configFlagName, a.globalFlags[a.configFlagName], flagValues, nil, lookupEnv)
		if err != nil {
			return fmt.Errorf("resolveFlag error for flag %s: %w", a.configFlagName, err)
		}
		if resolved {
			configPath := flagValues[flag.Name(a.configFlag.ConfigPath)].Get().(path.Path)
			configPathStr, err := configPath.Expand()
			if err != nil {
				return fmt.Errorf("error expanding config path ( %s ) : %w", configPath, err)
			}
			configReader, err = a.newConfigReader(configPathStr)
			if err != nil {
				return fmt.Errorf("error reading config path ( %s ) : %w", configPath, err)
			}

		}
	}

	// resolve app global flags
	for flagName, fl := range a.globalFlags {
		_, err := resolveFlag2(flagName, fl, flagValues, configReader, lookupEnv)
		if err != nil {
			return fmt.Errorf("resolveFlag error for flag %s: %w", flagName, err)
		}
	}

	// resolve current command flags
	if currentCommand != nil { // can be nil in the case of --help
		for flagName, fl := range currentCommand.Flags {
			_, err := resolveFlag2(flagName, fl, flagValues, configReader, lookupEnv)
			if err != nil {
				return fmt.Errorf("resolveFlag error for flag %s: %w", flagName, err)
			}
		}
	}

	return nil
}

// next steps:
// - port gzc Parse -> Parse2
// - make warg.Parse call Parse2 instead of doing the parsing
// - make all the tests pass (unsetsentinel, etc...)
// - test against CLI apps
// - release version
// - delete old parsing code
// - update warg.Parse's signature and tests
// - actually add tab completion (need to stringify values so they can be suggested as flag values)
