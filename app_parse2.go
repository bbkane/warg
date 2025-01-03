package warg

import (
	"fmt"
	"slices"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/common"
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

type UnsetFlagNameSet map[flag.Name]struct{}

func (u UnsetFlagNameSet) Add(name flag.Name) {
	u[name] = struct{}{}
}

func (u UnsetFlagNameSet) Delete(name flag.Name) {
	delete(u, name)
}

func (u UnsetFlagNameSet) Contains(name flag.Name) bool {
	_, exists := u[name]
	return exists
}

type ParseResult2 struct {
	SectionPath    []string
	CurrentSection *section.SectionT

	CurrentCommandName command.Name
	CurrentCommand     *command.Command

	CurrentFlagName flag.Name
	CurrentFlag     *flag.Flag
	FlagValues      FlagValueMap
	UnsetFlagNames  UnsetFlagNameSet

	HelpPassed bool
	State      ParseState
}

func (a *App) parseArgs(args []string) (ParseResult2, error) {
	pr := ParseResult2{
		SectionPath:    nil,
		CurrentSection: &a.rootSection,

		CurrentCommandName: "",
		CurrentCommand:     nil,

		CurrentFlagName: "",
		CurrentFlag:     nil,
		FlagValues:      make(FlagValueMap),
		UnsetFlagNames:  make(UnsetFlagNameSet),

		HelpPassed: false,
		State:      Parse_ExpectingSectionOrCommand,
	}

	aliasToFlagName := make(map[flag.Name]flag.Name)
	for flagName, fl := range a.globalFlags {
		if fl.Alias != "" {
			aliasToFlagName[flag.Name(fl.Alias)] = flagName
		}
	}

	// fill the FlagValues map with empty values from the app
	for flagName := range a.globalFlags {
		val := a.globalFlags[flagName].EmptyValueConstructor()
		pr.FlagValues[flagName] = val
	}

	for i, arg := range args {

		// --help <helptype> or --help must be the last thing passed and can appear at any state we aren't expecting a flag value
		if i >= len(args)-2 &&
			(flag.Name(arg) == a.helpFlagName || flag.Name(arg) == a.helpFlagAlias) &&
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

				// fill the FlagValues map with empty values from the command
				// All names in (command flag names, command flag aliases, global flag names, global flag aliases)
				// should be unique because app.Validate should have caught any conflicts
				for flagName, f := range pr.CurrentCommand.Flags {
					pr.FlagValues[flagName] = f.EmptyValueConstructor()

					if f.Alias != "" {
						aliasToFlagName[flag.Name(f.Alias)] = flagName
					}

				}
				pr.State = Parse_ExpectingFlagNameOrEnd
			} else {
				return pr, fmt.Errorf("expecting section or command, got %s", arg)
			}

		case Parse_ExpectingFlagNameOrEnd:
			flagName := flag.Name(arg)
			if actualFlagName, exists := aliasToFlagName[flagName]; exists {
				flagName = actualFlagName
			}
			fl := findFlag(flagName, a.globalFlags, pr.CurrentCommand.Flags)
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

func findFlag(flagName flag.Name, globalFlags flag.FlagMap, currentCommandFlags flag.FlagMap) *flag.Flag {
	if fl, exists := globalFlags[flagName]; exists {
		return &fl
	}
	if fl, exists := currentCommandFlags[flagName]; exists {
		return &fl
	}
	return nil
}

func resolveFlag2(
	flagName flag.Name,
	fl flag.Flag,
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

func (a *App) resolveFlags(currentCommand *command.Command, flagValues FlagValueMap, lookupEnv LookupFunc, unsetFlagNames UnsetFlagNameSet) error {
	// resolve config flag first and try to get a reader
	var configReader config.Reader
	if a.configFlagName != "" {
		err := resolveFlag2(
			a.configFlagName, a.globalFlags[a.configFlagName], flagValues, nil, lookupEnv, unsetFlagNames)
		if err != nil {
			return fmt.Errorf("resolveFlag error for flag %s: %w", a.configFlagName, err)
		}
		if flagValues[a.configFlagName].UpdatedBy() != value.UpdatedByUnset {
			configPath := flagValues[a.configFlagName].Get().(path.Path)
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
	pr, err := a.parseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("Parse error: %w", err)
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
	for flagName, flag := range a.globalFlags {
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

	// --config flag...
	// original Parse treats it specially
	// Parse2 expects it to be in app.GlobalFlags
	if app.configFlag != nil {
		app.globalFlags[app.configFlagName] = *app.configFlag
	}

	pr2, err := app.Parse2(parseOptHolder.Args[1:], parseOptHolder.LookupFunc)
	if err != nil {
		return nil, fmt.Errorf("parseWithOptHolder2 err: %w", err)
	}

	// build ftar.AvailableFlags - it's a map of string to flag for the app globals + current command. Don't forget to set each flag.IsCommandFlag and Value for now..
	// TODO:
	ftarAllowedFlags := make(flag.FlagMap)
	for flagName, fl := range app.globalFlags {
		fl.Value = pr2.FlagValues[flagName]
		fl.IsCommandFlag = false
		ftarAllowedFlags.AddFlag(flagName, fl)
	}

	// If we're in Parse_ExpectingSectionOrCommand, we haven't received a command
	if pr2.State != Parse_ExpectingSectionOrCommand {
		for flagName, fl := range pr2.CurrentCommand.Flags {
			fl.Value = pr2.FlagValues[flagName]
			fl.IsCommandFlag = true
			ftarAllowedFlags.AddFlag(flagName, fl)
		}
	}

	// port pfs
	pfs := pr2.FlagValues.ToPassedFlags()

	// port gar.Path
	garPath := pr2.SectionPath
	if pr2.CurrentCommandName != "" {
		garPath = slices.Concat(pr2.SectionPath, []string{string(pr2.CurrentCommandName)})
	}

	// TODO: handle aliases and sentinel values later

	if pr2.CurrentCommand == nil { // we got a section
		// no legit actions, just print the help
		helpInfo := common.HelpInfo{
			AvailableFlags: ftarAllowedFlags,
			RootSection:    app.rootSection,
		}
		// We know the helpFlag has a default so this is safe
		helpType := ftarAllowedFlags[flag.Name(app.helpFlagName)].Value.Get().(string)
		for _, e := range app.helpMappings {
			if e.Name == helpType {
				pr := ParseResult{
					Context: command.Context{
						AppName: app.name,
						Context: parseOptHolder.Context,
						Flags:   pfs,
						Path:    garPath,
						Stderr:  parseOptHolder.Stderr,
						Stdout:  parseOptHolder.Stdout,
						Version: app.version,
					},
					Action: e.SectionHelp(pr2.CurrentSection, helpInfo),
				}
				return &pr, nil
			}
		}
		return nil, fmt.Errorf("some problem with section help: info: %v", helpInfo)
	} else if pr2.CurrentCommand != nil { // we got a command
		if pr2.HelpPassed {
			helpInfo := common.HelpInfo{
				AvailableFlags: ftarAllowedFlags,
				RootSection:    app.rootSection,
			}
			// We know the helpFlag has a default so this is safe
			helpType := ftarAllowedFlags[flag.Name(app.helpFlagName)].Value.Get().(string)
			for _, e := range app.helpMappings {
				if e.Name == helpType {
					pr := ParseResult{
						Context: command.Context{
							AppName: app.name,
							Context: parseOptHolder.Context,
							Flags:   pfs,
							Path:    garPath,
							Stderr:  parseOptHolder.Stderr,
							Stdout:  parseOptHolder.Stdout,
							Version: app.version,
						},
						Action: e.CommandHelp(pr2.CurrentCommand, helpInfo),
					}
					return &pr, nil
				}
			}
			return nil, fmt.Errorf("some problem with command help: info: %v", helpInfo)
		} else {
			pr := ParseResult{
				Context: command.Context{
					AppName: app.name,
					Context: parseOptHolder.Context,
					Flags:   pfs,
					Path:    garPath,
					Stderr:  parseOptHolder.Stderr,
					Stdout:  parseOptHolder.Stdout,
					Version: app.version,
				},
				Action: pr2.CurrentCommand.Action,
			}
			return &pr, nil
		}

	} else {
		return nil, fmt.Errorf("internal Error: invalid parse state: currentSection == %v, currentCommand == %v", pr2.SectionPath, pr2.CurrentCommandName)
	}
}
