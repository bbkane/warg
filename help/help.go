package help

import (
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/allcommands"
	"go.bbkane.com/warg/help/detailed"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func DefaultHelpCommandMap() wargcore.CommandMap {
	return wargcore.CommandMap{
		"default":     wargcore.HelpToCommand(detailed.DetailedCommandHelp, allcommands.AllCommandsSectionHelp),
		"detailed":    wargcore.HelpToCommand(detailed.DetailedCommandHelp, detailed.DetailedSectionHelp),
		"outline":     wargcore.HelpToCommand(OutlineCommandHelp, OutlineSectionHelp),
		"allcommands": wargcore.HelpToCommand(detailed.DetailedCommandHelp, allcommands.AllCommandsSectionHelp),
	}
}

func DefaultHelpFlagMap(defaultChoice string, choices []string) wargcore.FlagMap {
	return wargcore.FlagMap{
		"--help": flag.New(
			"Print help",
			scalar.String(
				scalar.Choices(choices...),
				scalar.Default(defaultChoice),
			),
			flag.Alias("-h"),
		),
	}
}
