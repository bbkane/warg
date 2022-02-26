package help

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
)

func outlineFlagHelper(w io.Writer, color *gocolor.Color, flagName flag.Name, f flag.Flag, indent int) {
	str := fmtFlagName(color, flagName)
	if f.Alias != "" {
		str = str + " , " + fmtFlagAlias(color, f.Alias)
	}
	fmt.Fprintln(w, leftPad("# "+string(f.HelpShort), "  ", indent))
	fmt.Fprintln(
		w,
		leftPad(str, "  ", indent),
	)
}

func outlineHelper(w io.Writer, color *gocolor.Color, sec section.SectionT, indent int) {
	// section flags
	for _, k := range sec.Flags.SortedNames() {
		outlineFlagHelper(w, color, k, sec.Flags[flag.Name(k)], indent)
	}

	// commands and command flags
	for _, comName := range sec.Commands.SortedNames() {
		com := sec.Commands[command.Name(comName)]
		fmt.Fprintln(w, leftPad("# "+string(com.HelpShort), "  ", indent))
		fmt.Fprintln(
			w,
			leftPad(fmtCommandName(color, command.Name(comName)), "  ", indent),
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
			leftPad("# "+string(childSec.HelpShort), "  ", indent),
		)
		fmt.Fprintln(
			w,
			leftPad(fmtSectionName(color, k), "  ", indent),
		)
		outlineHelper(w, color, childSec, indent+1)
	}

}

func OutlineSectionHelp(file *os.File, _ *section.SectionT, hi HelpInfo) command.Action {
	return func(pf flag.PassedFlags) error {
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(pf, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		fmt.Fprintln(f, "# "+string(hi.RootSection.HelpShort))
		fmt.Fprintf(f, "%s\n", fmtSectionName(&col, section.Name(hi.AppName)))

		outlineHelper(f, &col, hi.RootSection, 1)

		return nil
	}
}

func OutlineCommandHelp(file *os.File, cur *command.Command, helpInfo HelpInfo) command.Action {
	return OutlineSectionHelp(file, nil, helpInfo)
}
