package warg

import (
	"bufio"
	"fmt"
	"os"

	"go.bbkane.com/warg/styles"
)

func outlineHelper(p *styles.Printer, s *styles.Styles, sec Section, indent int) {
	// commands and command flags
	for _, comName := range sec.Cmds.SortedNames() {
		p.Println(
			leftPad(s.CommandName(string(comName)), "  ", indent),
		)
	}

	// sections
	for _, k := range sec.Sections.SortedNames() {
		childSec := sec.Sections[k]
		p.Println(
			leftPad(s.SectionName(k), "  ", indent),
		)
		outlineHelper(p, s, childSec, indent+1)
	}

}

func outlineHelp() Cmd {
	action := func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		s, err := conditionallyEnableStyle(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		p := styles.NewPrinter(f)

		p.Println("# " + string(cmdCtx.App.RootSection.HelpShort))
		p.Println(s.SectionName(string(cmdCtx.App.Name)))

		outlineHelper(p, &s, cmdCtx.App.RootSection, 1)

		return nil
	}
	return NewCmd("", action)
}
