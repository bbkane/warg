package common

import (
	"os"
	"sort"

	"github.com/mattn/go-isatty"
	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
)

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
func ConditionallyEnableColor(pf cli.PassedFlags, file *os.File) (gocolor.Color, error) {
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

func SectionHelpToCommand(secHelp cli.SectionHelp) cli.Command {
	return command.NewCommand(
		"", // This is never visible to the user as this command is generated from the help flag
		func(cmdCtx cli.Context) error {
			// build ftar.AvailableFlags - it's a map of string to flag for the app globals + current command. Don't forget to set each flag.IsCommandFlag and Value for now..
			// TODO:
			ftarAllowedFlags := make(cli.FlagMap)
			for flagName, fl := range cmdCtx.App.GlobalFlags {
				fl.Value = cmdCtx.ParseResult.FlagValues[flagName]
				fl.IsCommandFlag = false
				ftarAllowedFlags.AddFlag(flagName, fl)
			}

			// If we're in Parse_ExpectingSectionOrCommand, we haven't received a command
			if cmdCtx.ParseResult.State != cli.Parse_ExpectingSectionOrCommand {
				for flagName, fl := range cmdCtx.ParseResult.CurrentCommand.Flags {
					fl.Value = cmdCtx.ParseResult.FlagValues[flagName]
					fl.IsCommandFlag = true
					ftarAllowedFlags.AddFlag(flagName, fl)
				}
			}

			sec := cmdCtx.ParseResult.CurrentSection
			hi := cli.HelpInfo{
				AvailableFlags: ftarAllowedFlags,
				RootSection:    cmdCtx.App.RootSection,
			}

			return secHelp(sec, hi)(cli.Context{}) //nolint:exhaustruct  // this context is not used and this is temp code to ease the porting
		},
	)
}

func CommandHelpToCommand(commandHelp cli.CommandHelp) cli.Command {
	return command.NewCommand(
		"", // This is never visible to the user as this command is generated from the help flag
		func(cmdCtx cli.Context) error {
			// build ftar.AvailableFlags - it's a map of string to flag for the app globals + current command. Don't forget to set each flag.IsCommandFlag and Value for now..
			// TODO:
			ftarAllowedFlags := make(cli.FlagMap)
			for flagName, fl := range cmdCtx.App.GlobalFlags {
				fl.Value = cmdCtx.ParseResult.FlagValues[flagName]
				fl.IsCommandFlag = false
				ftarAllowedFlags.AddFlag(flagName, fl)
			}

			// If we're in Parse_ExpectingSectionOrCommand, we haven't received a command
			if cmdCtx.ParseResult.State != cli.Parse_ExpectingSectionOrCommand {
				for flagName, fl := range cmdCtx.ParseResult.CurrentCommand.Flags {
					fl.Value = cmdCtx.ParseResult.FlagValues[flagName]
					fl.IsCommandFlag = true
					ftarAllowedFlags.AddFlag(flagName, fl)
				}
			}

			com := cmdCtx.ParseResult.CurrentCommand
			hi := cli.HelpInfo{
				AvailableFlags: ftarAllowedFlags,
				RootSection:    cmdCtx.App.RootSection,
			}

			return commandHelp(com, hi)(cli.Context{}) //nolint:exhaustruct  // this context is not used and this is temp code to ease the porting
		},
	)
}
