package help

import (
	"os"

	"go.bbkane.com/gocolor"

	"github.com/mattn/go-isatty"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
)

// HelpInfo lists common information available to a help function
type HelpInfo struct {
	// AppName as defined by warg.New()
	AppName string
	// Path passed either to a command or a section
	Path []string
	// AvailableFlags for the current section or commmand, including inherted flags from parent sections.
	// All flags are Resolved if possible (i.e., flag.SetBy != "")
	AvailableFlags flag.FlagMap
	// RootSection of the app. Especially useful for printing all sections and commands
	RootSection section.SectionT
}

type CommandHelp func(file *os.File, cur *command.Command, helpInfo HelpInfo) command.Action
type SectionHelp func(file *os.File, cur *section.SectionT, helpInfo HelpInfo) command.Action

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
		{Name: "default", CommandHelp: DetailedCommandHelp, SectionHelp: AllCommandsSectionHelp},
		{Name: "detailed", CommandHelp: DetailedCommandHelp, SectionHelp: DetailedSectionHelp},
		{Name: "outline", CommandHelp: OutlineCommandHelp, SectionHelp: OutlineSectionHelp},
		// allcommands list child commands, so it doesn't really make sense for a command
		{Name: "allcommands", CommandHelp: DetailedCommandHelp, SectionHelp: AllCommandsSectionHelp},
	}
}

// leftPad pads a string `s` with pad `pad` `plength` times
//
// In Python: (pad * plength) + s
func leftPad(s string, pad string, plength int) string {
	// https://stackoverflow.com/a/45456649/2958070
	for i := 0; i < plength; i++ {
		s = pad + s
	}
	return s
}

// SetColor looks for a passed --color flag with an underlying string value. If
// it exists and is set to "true", or if it exists, is set to "auto",
// and the passed file is a TTY, an enabled Color is returned.
func ConditionallyEnableColor(pf flag.PassedFlags, file *os.File) (gocolor.Color, error) {
	// default to trying to use color
	useColor := "auto"
	// respect a --color string
	if useColorI, exists := pf["--color"]; exists {
		if useColorUnder, isStr := useColorI.(string); isStr {
			useColor = useColorUnder
		}
	}

	startEnabled := useColor == "true" || (useColor == "auto" && isatty.IsTerminal(file.Fd()))
	return gocolor.Prepare(startEnabled)

}

func fmtHeader(col *gocolor.Color, header string) string {
	return col.Add(col.Bold+col.Underline, header)
}

func fmtSectionName(col *gocolor.Color, sectionName section.Name) string {
	return col.Add(col.Bold+col.FgCyan, string(sectionName))
}

func fmtCommandName(col *gocolor.Color, commandName command.Name) string {
	return col.Add(col.Bold+col.FgGreen, string(commandName))
}

func fmtFlagName(col *gocolor.Color, flagName flag.Name) string {
	return col.Add(col.Bold+col.FgYellow, string(flagName))
}

func fmtFlagAlias(col *gocolor.Color, flagAlias flag.Name) string {
	return col.Add(col.Bold+col.FgYellow, string(flagAlias))
}
