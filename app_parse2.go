package warg

import (
	"fmt"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

// -- FlagValue

type FlagValue struct {
	SetBy string
	Value value.Value
}

type FlagValueMap map[flag.Name]*FlagValue

func (m FlagValueMap) ToPassedFlags() command.PassedFlags {
	pf := make(command.PassedFlags)
	for name, f := range m {
		if f.SetBy != "" {
			pf[string(name)] = f.Value
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
		val, err := a.globalFlags[flagName].EmptyValueConstructor()
		// TODO: make this not an error!
		if err != nil {
			panic(err)
		}
		pr.FlagValues[flagName] = &FlagValue{
			SetBy: "",
			Value: val,
		}
	}

	for i, arg := range args {

		// --help <helptype> or --help must be the last thing passed and can appear at any state we aren't expecting a flag value
		if i >= len(args)-2 &&
			flag.Name(arg) == a.helpFlagName &&
			pr.State != Parse_ExpectingFlagValue {

			pr.HelpPassed = true
			// set the value of --help if an arg was passed, otherwise let it resolve with the rest of them...
			if i == len(args)-2 {
				err := pr.FlagValues[a.helpFlagName].Value.Update(args[i+1])
				if err != nil {
					return pr, fmt.Errorf("error updating help flag: %w", err)
				}
				pr.FlagValues[a.helpFlagName].SetBy = "passedarg"
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

				for flagName := range pr.CurrentCommand.Flags {
					_, exists := pr.FlagValues[flagName]
					if exists {
						// NOTE: move this check to app construction
						panic("app flags and command flags cannot share a name: " + flagName)
					}
					pr.FlagValues[flagName] = &FlagValue{}
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
			err := pr.FlagValues[pr.CurrentFlagName].Value.Update(arg)
			if err != nil {
				return pr, err
			}
			pr.FlagValues[pr.CurrentFlagName].SetBy = "passedarg"
			pr.State = Parse_ExpectingFlagNameOrEnd

		default:
			panic("unexpected state: " + pr.State)
		}
	}
	return pr, nil
}
