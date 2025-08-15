package warg

import (
	"os"
	"sort"

	"github.com/mattn/go-isatty"
	"go.bbkane.com/gocolor"
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

type CmdHelp func(cur *Cmd, helpInfo HelpInfo) Action
type SectionHelp func(cur *Section, helpInfo HelpInfo) Action

// HelpInfo lists common information available to a help function
type HelpInfo struct {

	// AvailableFlags for the current section or commmand, including inherted flags from parent sections.
	// All flags are Resolved if possible (i.e., flag.SetBy != "")
	AvailableFlags FlagMap
	// RootSection of the app. Especially useful for printing all sections and commands
	RootSection Section
}

func HelpToCommand(commandHelp CmdHelp, secHelp SectionHelp) Cmd {
	return Cmd{ //nolint:exhaustruct  // This help is never used since this is a generated command
		Action: func(cmdCtx CmdContext) error {
			// build ftar.AvailableFlags - it's a map of string to flag for the app globals + current command. Don't forget to set each flag.IsCommandFlag and Value for now..
			// TODO:
			ftarAllowedFlags := make(FlagMap)
			for flagName, fl := range cmdCtx.App.GlobalFlags {
				fl.Value = cmdCtx.ParseState.FlagValues[flagName]
				fl.IsCommandFlag = false
				ftarAllowedFlags.AddFlag(flagName, fl)
			}

			// If we're in Parse_ExpectingSectionOrCommand, we haven't received a command
			if cmdCtx.ParseState.ExpectingArg != ExpectingArg_SectionOrCommand {
				for flagName, fl := range cmdCtx.ParseState.CurrentCommand.Flags {
					fl.Value = cmdCtx.ParseState.FlagValues[flagName]
					fl.IsCommandFlag = true
					ftarAllowedFlags.AddFlag(flagName, fl)
				}
			}

			hi := HelpInfo{
				AvailableFlags: ftarAllowedFlags,
				RootSection:    cmdCtx.App.RootSection,
			}
			com := cmdCtx.ParseState.CurrentCommand
			if com != nil {
				return commandHelp(com, hi)(cmdCtx)
			} else {
				return secHelp(cmdCtx.ParseState.CurrentSection, hi)(cmdCtx)
			}
		},
	}

}

// LeftPad pads a string `s` with pad `pad` `plength` times
//
// In Python: (pad * plength) + s
func LeftPad(s string, pad string, plength int) string {
	// https://stackoverflow.com/a/45456649/2958070
	for i := 0; i < plength; i++ {
		s = pad + s
	}
	return s
}

// ConditionallyEnableColor looks for a passed --color flag with an underlying string value. If
// it exists and is set to "true", or if it exists, is set to "auto",
// and the passed file is a TTY, an enabled Color is returned.
func ConditionallyEnableColor(pf PassedFlags, file *os.File) (gocolor.Color, error) {
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

func FmtHeader(col *gocolor.Color, header string) string {
	return col.Add(col.Bold+col.Underline, header)
}

func FmtSectionName(col *gocolor.Color, sectionName string) string {
	return col.Add(col.Bold+col.FgCyan, string(sectionName))
}

func FmtCommandName(col *gocolor.Color, commandName string) string {
	return col.Add(col.Bold+col.FgGreen, string(commandName))
}

func FmtFlagName(col *gocolor.Color, flagName string) string {
	return col.Add(col.Bold+col.FgYellow, string(flagName))
}

func FmtFlagAlias(col *gocolor.Color, flagAlias string) string {
	return col.Add(col.Bold+col.FgYellow, string(flagAlias))
}

// SortedKeys returns the keys of the map m in sorted order.
// copied and modified from https://cs.opensource.google/go/x/exp/+/master:maps/maps.go;l=10;drc=79cabaa25d7518588d46eb676385c8dff49670c3
func SortedKeys[M ~map[string]V, V any](m M) []string {
	r := make([]string, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}
