package help

import (
	"bufio"
	"fmt"
	"os"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
)

func AllCommandsSectionHelp(file *os.File, cur *section.SectionT, helpInfo HelpInfo) command.Action {
	return func(pf flag.PassedFlags) error {

		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(pf, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		// Print top help section
		if cur.HelpLong != "" {
			fmt.Fprintf(f, "%s\n", cur.HelpLong)
		} else {
			fmt.Fprintf(f, "%s\n", cur.HelpShort)
		}

		fmt.Fprintln(f)

		fmt.Fprintln(f, fmtHeader(&col, "All Commands")+" (use <cmd> -h to see flag details):")
		fmt.Fprintln(f)

		path := []section.Name{section.Name(helpInfo.AppName)}
		for _, e := range helpInfo.Path {
			path = append(path, section.Name(e))
		}

		it := cur.BreadthFirst(path)
		for it.HasNext() {
			flatSec := it.Next()

			for _, name := range flatSec.Sec.Commands.SortedNames() {

				com := flatSec.Sec.Commands[name]
				fmt.Fprint(f, "  # ")
				fmt.Fprintln(f, com.HelpShort)

				fmt.Fprintf(f, "  ")

				for _, p := range flatSec.Path {
					fmt.Fprintf(f, fmtCommandName(&col, command.Name(p))+" ")
				}
				fmt.Fprintln(f, fmtCommandName(&col, name))

				fmt.Fprintln(f)
			}

		}
		if cur.Footer != "" {
			fmt.Fprintln(f, fmtHeader(&col, "Footer")+":")
			fmt.Fprintln(f)
			fmt.Fprintf(f, "%s\n", cur.Footer)
		}

		return nil
	}
}
