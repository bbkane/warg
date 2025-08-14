package warg

import (
	"go.bbkane.com/warg/value/scalar"
)

func DefaultHelpCommandMap() CmdMap {
	return CmdMap{
		"default":     HelpToCommand(DetailedCommandHelp, AllCommandsSectionHelp),
		"detailed":    HelpToCommand(DetailedCommandHelp, DetailedSectionHelp),
		"outline":     HelpToCommand(OutlineCommandHelp, OutlineSectionHelp),
		"allcommands": HelpToCommand(DetailedCommandHelp, AllCommandsSectionHelp),
	}
}

func DefaultHelpFlagMap(defaultChoice string, choices []string) FlagMap {
	return FlagMap{
		"--help": NewFlag(
			"Print help",
			scalar.String(
				scalar.Choices(choices...),
				scalar.Default(defaultChoice),
			),
			Alias("-h"),
		),
	}
}
