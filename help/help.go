package help

import (
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/help/allcommands"
	"go.bbkane.com/warg/help/common"
	"go.bbkane.com/warg/help/detailed"
	"go.bbkane.com/warg/section"
)

type CommandHelp func(cur *command.Command, helpInfo common.HelpInfo) command.Action
type SectionHelp func(cur *section.SectionT, helpInfo common.HelpInfo) command.Action

// HelpFlagMapping adds a new option to your --help flag
type HelpFlagMapping struct {
	Name        string
	CommandHelp CommandHelp
	SectionHelp SectionHelp
}

// BuiltinHelpFlagMappings is a convenience method that contains help flag mappings built into warg.
// Feel free to use it as a base to custimize help functions for your users
func BuiltinHelpFlagMappings() []HelpFlagMapping {
	return []HelpFlagMapping{
		{Name: "default", CommandHelp: detailed.DetailedCommandHelp, SectionHelp: allcommands.AllCommandsSectionHelp},
		{Name: "detailed", CommandHelp: detailed.DetailedCommandHelp, SectionHelp: detailed.DetailedSectionHelp},
		{Name: "outline", CommandHelp: OutlineCommandHelp, SectionHelp: OutlineSectionHelp},
		// allcommands list child commands, so it doesn't really make sense for a command
		{Name: "allcommands", CommandHelp: detailed.DetailedCommandHelp, SectionHelp: allcommands.AllCommandsSectionHelp},
	}
}
