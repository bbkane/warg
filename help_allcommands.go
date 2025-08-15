package warg

import (
	"bufio"
	"fmt"
	"os"
)

func allCommandsSectionHelp(cur *Section, helpInfo helpInfo) Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout

		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(cmdCtx.Flags, file)
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

		path := []string{string(cmdCtx.App.Name)}
		for _, e := range cmdCtx.ParseState.SectionPath {
			path = append(path, string(e))
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
					fmt.Fprint(f, fmtCommandName(&col, string(p))+" ")
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
