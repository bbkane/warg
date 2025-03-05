package allcommands

import (
	"bufio"
	"fmt"
	"os"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/help/common"
	"go.bbkane.com/warg/section"
)

func AllCommandsSectionHelp(cur *section.SectionT, helpInfo common.HelpInfo) command.Action {
	return func(cmdCtx command.Context) error {
		file := cmdCtx.Stdout

		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := common.ConditionallyEnableColor(cmdCtx.Flags, file)
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

		fmt.Fprintln(f, common.FmtHeader(&col, "All Commands")+" (use <cmd> -h to see flag details):")
		fmt.Fprintln(f)

		path := []section.Name{section.Name(cmdCtx.AppName)}
		for _, e := range cmdCtx.Path {
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
					fmt.Fprint(f, common.FmtCommandName(&col, string(p))+" ")
				}
				fmt.Fprintln(f, common.FmtCommandName(&col, name))

				fmt.Fprintln(f)
			}

		}
		if cur.Footer != "" {
			fmt.Fprintln(f, common.FmtHeader(&col, "Footer")+":")
			fmt.Fprintln(f)
			fmt.Fprintf(f, "%s\n", cur.Footer)
		}

		return nil
	}
}
