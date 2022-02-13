package help

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/bbkane/gocolor"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
)

func outlineFlagHelper(w io.Writer, color *gocolor.Color, name flag.Name, f *flag.Flag, indent int) {
	fmt.Fprintf(
		w,
		"%s\n",
		leftPad(string(name), "  ", indent),
	)
}

func outlineHelper(w io.Writer, color *gocolor.Color, sec *section.SectionT, indent int) {
	for name, fl := range sec.Flags {
		// TODO: sort
		outlineFlagHelper(w, color, name, &fl, indent)
	}

	for name, com := range sec.Commands {
		// TOOD: sort
		fmt.Fprintf(
			w,
			"%s\n",
			leftPad(string(name), "  ", indent),
		)
		for name, fl := range com.Flags {
			outlineFlagHelper(w, color, name, &fl, indent+1)
		}
	}

	for name, sec := range sec.Sections {
		fmt.Fprintf(
			w,
			"%s\n",
			leftPad(string(name), "  ", indent),
		)
		outlineHelper(w, color, &sec, indent+1)
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

		fmt.Fprintf(f, "%s\n", hi.AppName)

		outlineHelper(f, &col, &hi.RootSection, 1)

		return nil
	}
}

func OutlineCommandHelp(file *os.File, cur *command.Command, helpInfo HelpInfo) command.Action {
	return OutlineSectionHelp(file, nil, helpInfo)
}
