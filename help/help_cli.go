package help

import (
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/help/common"
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
