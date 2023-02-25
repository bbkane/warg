package warg

import (
	"fmt"
	"strings"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/common"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

type flagStr struct {
	NameOrAlias string
	Value       string
	Consumed    bool
}

type gatherArgsResult struct {

	// Path holds the path to the current command/section. Does not include the app name
	Path []string
	// FlagStrs is a slice of flags and values passed from the CLI. It can't be a map because flags can have aliases and we need to preserve order
	FlagStrs []flagStr
	// HelpPassed records whether --help was passed. The help flag may be set to a default value, so we need to check whether it's passed explicitly
	// so we can decide whether it needs to be acted upon
	HelpPassed bool
}

func containsString(haystack []string, needle string) bool {
	for _, w := range haystack {
		if w == needle {
			return true
		}
	}
	return false
}

// gatherArgs separates os.Args into a command path, a list of flags and their values from the CLI.
// It also takes note of whether --help was passed. To minimize ambiguitiy between a path element and an optional
// argument to --help, --help must be either not be passed, be the last string passed, or have exactly one value after it.
// See img/warg-gatherArgs-state-machine.png at the root of the repo for a diagram.
func gatherArgs(osArgs []string, helpFlagNames []string) (*gatherArgsResult, error) {
	res := &gatherArgsResult{
		Path:       nil,
		FlagStrs:   nil,
		HelpPassed: false,
	}

	startSt := "startSt"
	helpFlagPassedSt := "helpFlagPassedSt"
	helpValuePassedSt := "helpValuePassedSt"
	flagPassedSt := "flagPassedSt"

	state := startSt
	var currentFlagName string
	// Skip the name of the executable passed.
	for _, arg := range osArgs[1:] {
		// fmt.Printf("state: %v, arg: %v\n", state, arg)

		switch state {
		case startSt:
			if containsString(helpFlagNames, arg) {
				res.HelpPassed = true
				currentFlagName = arg
				state = helpFlagPassedSt
			} else if strings.HasPrefix(arg, "-") {
				currentFlagName = arg
				state = flagPassedSt
			} else { // cmd
				res.Path = append(res.Path, arg)
				state = startSt
			}
		case helpFlagPassedSt:
			res.FlagStrs = append(
				res.FlagStrs,
				flagStr{NameOrAlias: currentFlagName, Value: arg, Consumed: false},
			)
			state = helpValuePassedSt
		case helpValuePassedSt:
			return nil, fmt.Errorf("help flags should take maximally one arg, but more than one passed: %s", arg)
		case flagPassedSt:
			res.FlagStrs = append(
				res.FlagStrs,
				flagStr{NameOrAlias: currentFlagName, Value: arg, Consumed: false},
			)
			state = startSt
		default:
			return nil, fmt.Errorf("internal error: unknown state: %s", state)
		}
	}
	// check the only non-terminal state
	if state == flagPassedSt {
		return nil, fmt.Errorf("flag passed without value( %#v) . All flags must have one value passed. Repeat flags to accumulate values. Example: --level 9000", currentFlagName)
	}
	return res, nil
}

// flagNameToAlias is a map of flag name to flag alias
type flagNameToAlias map[flag.Name]flag.Name

// fitToAppResult holds the result of fitToApp
// Exactly one of Section or Command should hold something. The other should be nil
type fitToAppResult struct {
	Section            *section.SectionT
	Command            *command.Command
	Action             command.Action
	AllowedFlags       flag.FlagMap
	AllowedFlagAliases flagNameToAlias
}

// fitToApp takes the command entered by a user and uses it to "walk" down the apps command tree to build what the command was and what the available flags are.
func fitToApp(rootSection section.SectionT, path []string) (*fitToAppResult, error) {

	// AllowedFlags grows, as we traverse the tree; copy rootSection.Flags.
	// We need a pristine rootsection.Flags for --help printing
	allowedFlags := make(flag.FlagMap)
	for k, v := range rootSection.Flags {
		allowedFlags[k] = v
	}

	// validate passed command and get available flags
	ftar := fitToAppResult{
		Section:            &rootSection,
		Command:            nil, // we start with a section, not a command
		Action:             nil,
		AllowedFlags:       allowedFlags,
		AllowedFlagAliases: make(flagNameToAlias),
	}
	// Add any root flag aliases to AllowedFlagAliases
	// Wonder if I could put all this in one part of the code...
	for flagName, fl := range ftar.AllowedFlags {
		if fl.Alias != "" {
			ftar.AllowedFlagAliases[flagName] = fl.Alias
		}
	}
	childCommands := rootSection.Commands
	childSections := rootSection.Sections
	for _, word := range path {
		if command, exists := childCommands[command.Name(word)]; exists {
			ftar.Command = &command
			ftar.Section = nil
			ftar.Action = command.Action
			// once we're in a commmand, we should be at the end of the path
			// commands have no child commands or child sections
			childCommands = nil
			childSections = nil
			for flagName, fl := range command.Flags {
				// TODO: check if name exists already
				if fl.Alias != "" {
					ftar.AllowedFlagAliases[flagName] = fl.Alias
				}
				fl.IsCommandFlag = true
				ftar.AllowedFlags[flagName] = fl
			}
		} else if section, exists := childSections[section.Name(word)]; exists {
			ftar.Section = &section
			childCommands = section.Commands
			childSections = section.Sections
			for flagName, fl := range section.Flags {
				// TODO: check if key exists already
				if fl.Alias != "" {
					ftar.AllowedFlagAliases[flagName] = fl.Alias
				}
				ftar.AllowedFlags[flagName] = fl
			}
		} else {
			retErr := fmt.Errorf("expected command or section, but got %#v, try --help", word)
			return nil, retErr
		}
	}
	return &ftar, nil
}

// resolveFlag updates a flag's value from the command line, and then from the
// default value. flag should not be nil. deletes from flagStrs
func resolveFlag(
	fl *flag.Flag,
	name flag.Name,
	flagStrs []flagStr,
	configReader config.Reader,
	lookupEnv LookupFunc,
	aliases flagNameToAlias,
) error {
	// TODO: can I delete from flagStrs in the caller? then I wouldn't need to pass
	// flagStrs (just a potential strValues) into here and it's a more pure function

	val, err := fl.EmptyValueConstructor()
	if err != nil {
		return fmt.Errorf("flag error: %v: %w", name, err)
	}
	fl.Value = val

	// try to update from command line and consume from flagStrs
	// need to check flag.SetBy even in the first case because we could be resolving
	// flags multiple times (for instance --config gets resolved before this and also now)
	{
		strValues := []string{}
		for i := range flagStrs {
			// TODO: come back to theses string casts...
			if flagStrs[i].NameOrAlias == string(name) || flag.Name(flagStrs[i].NameOrAlias) == aliases[name] {
				strValues = append(strValues, flagStrs[i].Value)
				flagStrs[i].Consumed = true
			}
		}

		if fl.SetBy == "" && len(strValues) > 0 {

			_, isScalar := val.(value.ScalarValue)
			if isScalar && len(strValues) > 1 {
				return fmt.Errorf("flag error: %v: flag passed multiple times, it's value (type %v), can only be updated once", name, fl.Value.Description())

			}

			for _, v := range strValues {
				err = fl.Value.Update(v)
				if err != nil {
					return fmt.Errorf("error updating flag %v from passed flag value %v: %w", name, v, err)
				}
			}
			fl.SetBy = "passedflag"
		}
	}

	// update from config
	{
		if fl.SetBy == "" && configReader != nil {
			fpr, err := configReader.Search(fl.ConfigPath)
			if err != nil {
				return err
			}
			if fpr.Exists {
				if !fpr.IsAggregated {
					err := fl.Value.ReplaceFromInterface(fpr.IFace)
					if err != nil {
						return fmt.Errorf(
							"could not replace container type value: val: %#v , replacement: %#v, err: %w",
							fl.Value,
							fpr.IFace,
							err,
						)
					}
					fl.SetBy = "config"
				} else {
					v, ok := fl.Value.(value.SliceValue)
					if !ok {
						return fmt.Errorf("could not update scalar value with aggregated value from config: name: %v, configPath: %v", name, fl.ConfigPath)

					}
					under, ok := fpr.IFace.([]interface{})
					if !ok {
						return fmt.Errorf("expected []interface{}, got: %#v", under)
					}
					for _, e := range under {
						err = v.AppendFromInterface(e)
						if err != nil {
							return fmt.Errorf("could not update container type value: err: %w", err)
						}
					}
					fl.SetBy = "config"
					fl.Value = v
				}
			}
		}
	}

	// update from envvars
	{
		if fl.SetBy == "" && len(fl.EnvVars) > 0 {
			for _, e := range fl.EnvVars {
				val, exists := lookupEnv(e)
				if exists {
					err = fl.Value.Update(val)
					if err != nil {
						return fmt.Errorf("error updating flag %v from envvar %v: %w", name, val, err)
					}
					fl.SetBy = "envvar"
					break // stop looking for envvars
				}

			}
		}
	}

	// update from default
	{
		if fl.SetBy == "" && fl.Value.HasDefault() {
			fl.Value.ReplaceFromDefault()
			fl.SetBy = "appdefault"
		}
	}

	return nil
}

// ParseResult holds the result of parsing the command line.
type ParseResult struct {
	// Path to the command invoked. Does not include executable name (os.Args[0])
	Path []string // TODO: consider moving this inside the command.Context
	// Context holds the parsed information
	Context command.Context
	// Action holds the passed command's action to execute.
	Action command.Action
}

// Parse parses the args, but does not execute anything.
func (app *App) Parse(osArgs []string, osLookupEnv LookupFunc) (*ParseResult, error) {
	helpFlagNames := []string{string(app.helpFlagName)}
	if app.helpFlagAlias != "" {
		helpFlagNames = append(helpFlagNames, string(app.helpFlagAlias))
	}

	gar, err := gatherArgs(osArgs, helpFlagNames)
	if err != nil {
		return nil, err
	}

	ftar, err := fitToApp(app.rootSection, gar.Path)
	if err != nil {
		return nil, err
	}

	// fill the flags
	var configReader config.Reader
	// get the value of a potential passed --config flag first so we can use it
	// to resolve further flags
	if app.configFlag != nil {

		// Maybe this should go in fitToApp?
		if app.configFlag.Alias != "" {
			ftar.AllowedFlagAliases[app.configFlagName] = app.configFlag.Alias
		}

		// we're gonna make a config map out of this if everything goes well
		// so pass nil for the configreader now
		err = resolveFlag(
			app.configFlag,
			app.configFlagName,
			gar.FlagStrs,
			nil,
			osLookupEnv,
			ftar.AllowedFlagAliases,
		)
		if err != nil {
			return nil, err
		}
		// NOTE: this *should* always be a string
		configPath := app.configFlag.Value.Get().(string)
		configReader, err = app.newConfigReader(configPath)
		if err != nil {
			return nil, fmt.Errorf("error reading config path ( %s ) : %w", configPath, err)
		}
	}

	// Loop over allowed flags for the passed command and try to resolve them
	for name, fl := range ftar.AllowedFlags {

		err = resolveFlag(
			&fl,
			name,
			gar.FlagStrs,
			configReader,
			osLookupEnv,
			ftar.AllowedFlagAliases,
		)
		if err != nil {
			return nil, err
		}

		if !gar.HelpPassed {
			if fl.Required && fl.SetBy == "" {
				return nil, fmt.Errorf("flag required but not set: %s", name)
			}
		}

		ftar.AllowedFlags[name] = fl
	}

	// add the config flag so both help and actions can see it
	if app.configFlag != nil {
		ftar.AllowedFlags[app.configFlagName] = *app.configFlag
	}

	for _, e := range gar.FlagStrs {
		if !e.Consumed {
			return nil, fmt.Errorf("unrecognized flag: %v -> %v", e.NameOrAlias, e.Value)
		}
	}

	pfs := make(command.PassedFlags)
	for name, fl := range ftar.AllowedFlags {
		if fl.SetBy != "" {
			pfs[string(name)] = fl.Value.Get()
		}
	}

	// OK! Let's make the ParseResult for each case and gtfo
	if ftar.Section != nil && ftar.Command == nil {
		// no legit actions, just print the help
		helpInfo := common.HelpInfo{
			AppName:        app.name,
			Path:           gar.Path,
			AvailableFlags: ftar.AllowedFlags,
			RootSection:    app.rootSection,
		}
		// We know the helpFlag has a default so this is safe
		helpType := ftar.AllowedFlags[flag.Name(app.helpFlagName)].Value.Get().(string)
		for _, e := range app.helpMappings {
			if e.Name == helpType {
				pr := ParseResult{
					Path: gar.Path,
					Context: command.Context{
						Flags:  pfs,
						Stderr: app.Stderr,
						Stdout: app.Stdout,
					},
					Action: e.SectionHelp(ftar.Section, helpInfo),
				}
				return &pr, nil
			}
		}
		return nil, fmt.Errorf("some problem with section help: info: %v", helpInfo)
	} else if ftar.Section == nil && ftar.Command != nil {
		if gar.HelpPassed {
			helpInfo := common.HelpInfo{
				AppName:        app.name,
				Path:           gar.Path,
				AvailableFlags: ftar.AllowedFlags,
				RootSection:    app.rootSection,
			}
			// We know the helpFlag has a default so this is safe
			helpType := ftar.AllowedFlags[flag.Name(app.helpFlagName)].Value.Get().(string)
			for _, e := range app.helpMappings {
				if e.Name == helpType {
					pr := ParseResult{
						Path: gar.Path,
						Context: command.Context{
							Flags:  pfs,
							Stderr: app.Stderr,
							Stdout: app.Stdout,
						},
						Action: e.CommandHelp(ftar.Command, helpInfo),
					}
					return &pr, nil
				}
			}
			return nil, fmt.Errorf("some problem with section help: info: %v", helpInfo)
		} else {

			pr := ParseResult{
				Path: gar.Path,
				Context: command.Context{
					Flags:  pfs,
					Stderr: app.Stderr,
					Stdout: app.Stdout,
				},
				Action: ftar.Action,
			}
			return &pr, nil
		}
	} else {
		return nil, fmt.Errorf("internal Error: invalid parse state: currentSection == %v, currentCommand == %v", ftar.Section, ftar.Command)
	}
}
