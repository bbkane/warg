package help

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/help/common"
)

func outlineFlagHelper(w io.Writer, color *gocolor.Color, flagName string, f cli.Flag, indent int) {
	str := common.FmtFlagName(color, flagName)
	if f.Alias != "" {
		str = str + " , " + common.FmtFlagAlias(color, f.Alias)
	}
	fmt.Fprintln(w, common.LeftPad("# "+string(f.HelpShort), "  ", indent))
	fmt.Fprintln(
		w,
		common.LeftPad(str, "  ", indent),
	)
}

func outlineHelper(w io.Writer, color *gocolor.Color, sec cli.SectionT, indent int) {
	// commands and command flags
	for _, comName := range sec.Commands.SortedNames() {
		com := sec.Commands[string(comName)]
		fmt.Fprintln(w, common.LeftPad("# "+string(com.HelpShort), "  ", indent))
		fmt.Fprintln(
			w,
			common.LeftPad(common.FmtCommandName(color, string(comName)), "  ", indent),
		)
		for _, flagName := range com.Flags.SortedNames() {
			outlineFlagHelper(w, color, flagName, com.Flags[flagName], indent+1)
		}

	}

	// sections
	for _, k := range sec.Sections.SortedNames() {
		childSec := sec.Sections[k]
		fmt.Fprintln(
			w,
			common.LeftPad("# "+string(childSec.HelpShort), "  ", indent),
		)
		fmt.Fprintln(
			w,
			common.LeftPad(common.FmtSectionName(color, k), "  ", indent),
		)
		outlineHelper(w, color, childSec, indent+1)
	}

}

func OutlineSectionHelp(_ *cli.SectionT, hi cli.HelpInfo) cli.Action {
	return func(cmdCtx cli.Context) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := common.ConditionallyEnableColor(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		fmt.Fprintln(f, "# "+string(hi.RootSection.HelpShort))
		fmt.Fprintf(f, "%s\n", common.FmtSectionName(&col, string(cmdCtx.App.Name)))

		outlineHelper(f, &col, hi.RootSection, 1)

		return nil
	}
}

func OutlineCommandHelp(cur *cli.Command, helpInfo cli.HelpInfo) cli.Action {
	return OutlineSectionHelp(nil, helpInfo)
}

func OutlineSectionHelpCommand(cmdCtx cli.Context) error {

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

	return OutlineSectionHelp(sec, hi)(cli.Context{}) //nolint:exhaustruct  // this context is not used and this is temp code to ease the porting

}
