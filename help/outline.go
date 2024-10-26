package help

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/common"
	"go.bbkane.com/warg/section"
)

func outlineFlagHelper(w io.Writer, color *gocolor.Color, flagName flag.Name, f flag.Flag, indent int) {
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

func outlineHelper(w io.Writer, color *gocolor.Color, sec section.SectionT, indent int) {
	// commands and command flags
	for _, comName := range sec.Commands.SortedNames() {
		com := sec.Commands[command.Name(comName)]
		fmt.Fprintln(w, common.LeftPad("# "+string(com.HelpShort), "  ", indent))
		fmt.Fprintln(
			w,
			common.LeftPad(common.FmtCommandName(color, command.Name(comName)), "  ", indent),
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

func OutlineSectionHelp(_ *section.SectionT, hi common.HelpInfo) command.Action {
	return func(cmdCtx command.Context) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := common.ConditionallyEnableColor(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		fmt.Fprintln(f, "# "+string(hi.RootSection.HelpShort))
		fmt.Fprintf(f, "%s\n", common.FmtSectionName(&col, section.Name(cmdCtx.AppName)))

		outlineHelper(f, &col, hi.RootSection, 1)

		return nil
	}
}

func OutlineCommandHelp(cur *command.Command, helpInfo common.HelpInfo) command.Action {
	return OutlineSectionHelp(nil, helpInfo)
}
