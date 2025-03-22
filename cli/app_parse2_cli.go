package cli

import (
	"fmt"

	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value"
)

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
			(string(arg) == a.HelpFlagName || string(arg) == a.HelpFlagAlias) &&
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

func (a *App) Parse2(args []string, lookupEnv LookupFunc) (*ParseResult2, error) {

	// --config flag...
	// original Parse treats it specially
	// Parse2 expects it to be in app.GlobalFlags
	// TODO: rework the config flag handling
	if a.ConfigFlag != nil {
		a.GlobalFlags[a.ConfigFlagName] = *a.ConfigFlag
	}

	pr, err := a.parseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("Parse args error: %w", err)
	}

	// If we're in a section, just print the help
	if pr.State == Parse_ExpectingSectionOrCommand {
		pr.HelpPassed = true
	}

	// --help means we don't need to do a lot of error checking
	if pr.HelpPassed {
		err = a.resolveFlags(pr.CurrentCommand, pr.FlagValues, lookupEnv, pr.UnsetFlagNames)
		if err != nil {
			return nil, err
		}
		return &pr, nil
	}

	// ok, we're running a real command, let's do the error checking
	if pr.State != Parse_ExpectingFlagNameOrEnd {
		return nil, fmt.Errorf("unexpected parse state: %s", pr.State)
	}

	err = a.resolveFlags(pr.CurrentCommand, pr.FlagValues, lookupEnv, pr.UnsetFlagNames)
	if err != nil {
		return nil, err
	}

	missingRequiredFlags := []string{}
	for flagName, flag := range a.GlobalFlags {
		if flag.Required && pr.FlagValues[flagName].UpdatedBy() == value.UpdatedByUnset {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	for flagName, flag := range pr.CurrentCommand.Flags {
		if flag.Required && pr.FlagValues[flagName].UpdatedBy() == value.UpdatedByUnset {
			missingRequiredFlags = append(missingRequiredFlags, string(flagName))
		}
	}

	if len(missingRequiredFlags) > 0 {
		return nil, fmt.Errorf("missing but required flags: %s", missingRequiredFlags)
	}

	return &pr, nil

}

func (app *App) parseWithOptHolder2(parseOptHolder ParseOptHolder) (*ParseResult, error) {

	pr2, err := app.Parse2(parseOptHolder.Args[1:], parseOptHolder.LookupFunc)
	if err != nil {
		return nil, fmt.Errorf("Parse err: %w", err)
	}

	// TODO: handle aliases and sentinel values later

	// section or help passed
	if pr2.State == Parse_ExpectingSectionOrCommand || pr2.HelpPassed {
		helpType := pr2.FlagValues[app.HelpFlagName].Get().(string)
		// helpType := ftarAllowedFlags[string(app.HelpFlagName)].Value.Get().(string)
		for _, e := range app.HelpMappings {
			if e.Name == helpType {
				command := HelpToCommand(e.CommandHelp, e.SectionHelp)
				pr := ParseResult{
					Context: Context{
						App:         app,
						Context:     parseOptHolder.Context,
						Flags:       pr2.FlagValues.ToPassedFlags(),
						ParseResult: pr2,
						Stderr:      parseOptHolder.Stderr,
						Stdout:      parseOptHolder.Stdout,
					},
					Action: command.Action,
				}
				return &pr, nil
			}
		}
		return nil, fmt.Errorf("could not find help: %v, in help mappings: %v", helpType, app.HelpMappings)
	}
	if pr2.State == Parse_ExpectingFlagNameOrEnd {
		pr := ParseResult{
			Context: Context{
				App:         app,
				Context:     parseOptHolder.Context,
				Flags:       pr2.FlagValues.ToPassedFlags(),
				ParseResult: pr2,
				Stderr:      parseOptHolder.Stderr,
				Stdout:      parseOptHolder.Stdout,
			},
			Action: pr2.CurrentCommand.Action,
		}
		return &pr, nil
	}
	return nil, fmt.Errorf("internal Error: invalid parse state == %v: currentSection == %v, currentCommand == %v", pr2.State, pr2.SectionPath, pr2.CurrentCommandName)
}
