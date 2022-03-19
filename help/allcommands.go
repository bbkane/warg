package help

import (
	"bufio"
	"fmt"
	"os"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
)

func AllCommandsCommandHelp(file *os.File, cur *command.Command, helpInfo HelpInfo) command.Action {
	return AllCommandsSectionHelp(file, nil, helpInfo)
}

func AllCommandsSectionHelp(file *os.File, _ *section.SectionT, helpInfo HelpInfo) command.Action {
	return func(pf flag.PassedFlags) error {

		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(pf, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		cur := helpInfo.RootSection // TODO: acutually use cur
		// Print top help section
		if cur.HelpLong != "" {
			fmt.Fprintf(f, "%s\n", cur.HelpLong)
		} else {
			fmt.Fprintf(f, "%s\n", cur.HelpShort)
		}

		fmt.Fprintln(f)

		fmt.Fprintln(f, fmtHeader(&col, "All Commands")+" (use <cmd> -h to see flag details):")
		fmt.Fprintln(f)

		it := cur.BreadthFirst(section.Name(helpInfo.AppName))
		for it.HasNext() {
			flatSec := it.Next()

			for _, name := range flatSec.Sec.Commands.SortedNames() {

				com := flatSec.Sec.Commands[name]
				fmt.Fprint(f, "  # ")
				fmt.Fprintln(f, com.HelpShort)

				fmt.Fprintf(f, "  ")

				// fmt.Fprintln(f, helpInfo.AppName, helpInfo.Path, flatSec.ParentPath, flatSec.Name, name)

				for _, p := range flatSec.ParentPath {
					fmt.Fprintf(f, fmtCommandName(&col, command.Name(p))+" ")
				}
				fmt.Fprintf(f, fmtCommandName(&col, command.Name(flatSec.Name))+" ")
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
