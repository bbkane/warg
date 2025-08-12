package help

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/help/common"
	"go.bbkane.com/warg/wargcore"
)

func outlineHelper(w io.Writer, color *gocolor.Color, sec wargcore.Section, indent int) {
	// commands and command flags
	for _, comName := range sec.Commands.SortedNames() {
		fmt.Fprintln(
			w,
			common.LeftPad(common.FmtCommandName(color, string(comName)), "  ", indent),
		)
	}

	// sections
	for _, k := range sec.Sections.SortedNames() {
		childSec := sec.Sections[k]
		fmt.Fprintln(
			w,
			common.LeftPad(common.FmtSectionName(color, k), "  ", indent),
		)
		outlineHelper(w, color, childSec, indent+1)
	}

}

func OutlineSectionHelp(_ *wargcore.Section, hi wargcore.HelpInfo) wargcore.Action {
	return func(cmdCtx wargcore.Context) error {
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

func OutlineCommandHelp(cur *wargcore.Cmd, helpInfo wargcore.HelpInfo) wargcore.Action {
	return OutlineSectionHelp(nil, helpInfo)
}
