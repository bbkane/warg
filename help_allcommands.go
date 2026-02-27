package warg

import (
	"bufio"
	"fmt"
	"os"

	"go.bbkane.com/warg/styles"
)

func allCmdsSectionHelp() Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout

		f := bufio.NewWriter(file)
		defer f.Flush()

		s, err := conditionallyEnableStyle(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		p := styles.NewPrinter(file)

		cur := cmdCtx.ParseState.CurrentSection

		// Print top help section
		if cur.HelpLong != "" {
			p.Println(cur.HelpLong)
		} else {
			p.Println(cur.HelpShort)
		}

		p.Println()

		p.Println(s.Header("All Commands") + " (use <cmd> -h to see flag details):")
		p.Println()

		path := []string{string(cmdCtx.App.Name)}
		for _, e := range cmdCtx.ParseState.SectionPath {
			path = append(path, string(e))
		}

		it := cur.breadthFirst(path)
		for it.HasNext() {
			flatSec := it.Next()

			for _, name := range flatSec.Sec.Cmds.SortedNames() {

				com := flatSec.Sec.Cmds[name]
				p.Print("  # ")
				p.Println(com.HelpShort)

				p.Print("  ")
				for _, path := range flatSec.Path {
					p.Print(s.CommandName(string(path)) + " ")
				}
				p.Println(s.CommandName(name))

				p.Println()
			}

		}
		if cur.Footer != "" {
			p.Println(s.Header("Footer") + ":")
			p.Println()
			p.Println(cur.Footer)
		}

		return nil
	}
}
