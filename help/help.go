package help

import (
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/allcommands"
	"go.bbkane.com/warg/help/detailed"
	"go.bbkane.com/warg/value/scalar"
)

func DefaultHelpCommandMap() cli.CommandMap {
	return cli.CommandMap{
		"default":     cli.HelpToCommand(detailed.DetailedCommandHelp, allcommands.AllCommandsSectionHelp),
		"detailed":    cli.HelpToCommand(detailed.DetailedCommandHelp, detailed.DetailedSectionHelp),
		"outline":     cli.HelpToCommand(OutlineCommandHelp, OutlineSectionHelp),
		"allcommands": cli.HelpToCommand(detailed.DetailedCommandHelp, allcommands.AllCommandsSectionHelp),
	}
}

func DefaultHelpFlagMap(defaultChoice string, choices []string) cli.FlagMap {
	return cli.FlagMap{
		"--help": flag.NewFlag(
			"Print help",
			scalar.String(
				scalar.Choices(choices...),
				scalar.Default(defaultChoice),
			),
			flag.Alias("-h"),
		),
	}
}
