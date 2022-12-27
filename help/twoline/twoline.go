package twoline

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/common"
)

func printFlag(f io.Writer, color *gocolor.Color, flagName flag.Name, fl *flag.Flag) {
	common.FprintNoSpace(f, "  ")
	if string(fl.Alias) != "" {
		common.FprintNoSpace(
			f,
			common.FmtFlagAlias(color, fl.Alias),
			", ",
		)
	}
	common.FprintNoSpace(
		f,
		common.FmtFlagName(color, flagName),
	)
	if fl.SetBy != "" {
		common.FprintNoSpace(
			f,
			" = ",
			fl.Value.String(),
		)
	}
	common.FprintlnNoSpace(f)
	common.FprintlnNoSpace(
		f,
		"    : ",
		fl.HelpShort,
	)
}

func TwoLineCommandHelp(file *os.File, cur *command.Command, hi common.HelpInfo) command.Action {
	return func(ctx command.Context) error {
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := common.ConditionallyEnableColor(ctx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		// Print top help section
		if cur.HelpLong != "" {
			fmt.Fprintf(f, "%s\n\n", cur.HelpLong)
		} else {
			fmt.Fprintf(f, "%s\n\n", cur.HelpShort)
		}

		// compute sections for command flags and inherited flags,
		// then print their headers and them if they're not empty
		var commandFlagHelp bytes.Buffer
		var sectionFlagHelp bytes.Buffer

		for _, flagName := range hi.AvailableFlags.SortedNames() {
			fl := hi.AvailableFlags[flagName]
			if fl.IsCommandFlag {
				printFlag(&commandFlagHelp, &col, flagName, &fl)
			} else {
				printFlag(&sectionFlagHelp, &col, flagName, &fl)
			}
		}

		if commandFlagHelp.Len() > 0 {
			fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Command Flags"))
			fmt.Fprintln(f)
			_, _ = commandFlagHelp.WriteTo(f)
			fmt.Fprintln(f)

		}
		if sectionFlagHelp.Len() > 0 {
			fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Inherited Section Flags"))
			fmt.Fprintln(f)
			_, _ = sectionFlagHelp.WriteTo(f)
		}
		if cur.Footer != "" {
			fmt.Fprintln(f)
			fmt.Fprintln(f, col.Add(col.Underline+col.Bold, "Footer")+":")
			fmt.Fprintln(f)
			fmt.Fprintf(f, "%s\n", cur.Footer)
		}
		return nil
	}
}
