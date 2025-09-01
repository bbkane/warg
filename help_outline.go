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
	for _, comName := range sec.Cmds.SortedNames() {
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

func outlineHelp() Cmd {
	action := func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		fmt.Fprintln(f, "# "+string(cmdCtx.App.RootSection.HelpShort))
		fmt.Fprintf(f, "%s\n", fmtSectionName(&col, string(cmdCtx.App.Name)))

		outlineHelper(f, &col, cmdCtx.App.RootSection, 1)

		return nil
	}
	return NewCmd("", action)
}
