package warg

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"go.bbkane.com/gocolor"
)

func outlineHelper(w io.Writer, color *gocolor.Color, sec Section, indent int) {
	// commands and command flags
	for _, comName := range sec.Commands.SortedNames() {
		fmt.Fprintln(
			w,
			leftPad(fmtCommandName(color, string(comName)), "  ", indent),
		)
	}

	// sections
	for _, k := range sec.Sections.SortedNames() {
		childSec := sec.Sections[k]
		fmt.Fprintln(
			w,
			leftPad(fmtSectionName(color, k), "  ", indent),
		)
		outlineHelper(w, color, childSec, indent+1)
	}

}

func outlineSectionHelp(_ *Section, hi helpInfo) Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		fmt.Fprintln(f, "# "+string(hi.RootSection.HelpShort))
		fmt.Fprintf(f, "%s\n", fmtSectionName(&col, string(cmdCtx.App.Name)))

		outlineHelper(f, &col, hi.RootSection, 1)

		return nil
	}
}

func outlineCommandHelp(cur *Cmd, helpInfo helpInfo) Action {
	return outlineSectionHelp(nil, helpInfo)
}
