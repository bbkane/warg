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

func outlineHelper(w io.Writer, color *gocolor.Color, sec cli.Section, indent int) {
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

func OutlineSectionHelp(_ *cli.Section, hi cli.HelpInfo) cli.Action {
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
