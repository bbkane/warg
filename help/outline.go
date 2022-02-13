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

func outlineHelper(w io.Writer, color *gocolor.Color, sec *section.SectionT, indent int) {

}

func OutlineSectionHelp(file *os.File, _ *section.SectionT, hi HelpInfo) command.Action {
	return func(pf flag.PassedFlags) error {
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(pf, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		outlineHelper(f, &col, &hi.RootSection, 0)

		return nil
	}
}

func OutlineCommandHelp(file *os.File, cur *command.Command, helpInfo HelpInfo) command.Action {
	// return OutlineSectionHelp(file, nil, helpInfo)
	return command.DoNothing
}
