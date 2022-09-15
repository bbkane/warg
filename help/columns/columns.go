package columns

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/common"
)

const flagHelpSep = " : "
const flagIndent = "  "
const flagAliasNameSep = ", "
const flagNameValueSep = " = "

// padding returns the spaces to pad a name to fit a given length
func padding(n flag.Name, length int) string {
	lenFlagName := len(n)
	if lenFlagName > length {
		panic("lenFlagName > length")
	}
	return strings.Repeat(" ", length-lenFlagName)
}

// maxFlagColWidth calculates the max width of everything in the flag column.
// NOTE: flagIndent is just considered it's own column
func maxFlagColWidth(fm flag.FlagMap) int {
	m := 0
	for name, fl := range fm {
		dynamicLen := len(name)
		if fl.Alias != "" {
			dynamicLen = dynamicLen + len(fl.Alias) + len(flagAliasNameSep)
		}
		if fl.SetBy != "" {
			// TODO: account for compound type values!
			dynamicLen = dynamicLen + len(flagNameValueSep) + len(fl.Value.String())
		}
		if dynamicLen > m {
			m = dynamicLen
		}
	}
	return m
}

func fPrintNoSpace(w io.Writer, s ...any) {
	for _, si := range s {
		fmt.Fprint(w, si)
	}
}

func fPrintlnNoSpace(w io.Writer, s ...any) {
	fPrintNoSpace(w, s...)
	fmt.Fprintln(w)
}

func printFlag(f io.Writer, color *gocolor.Color, flagName flag.Name, fl *flag.Flag, maxFlagColWidth_ int) {
	fPrintNoSpace(f, flagIndent)
	paddingWidth := maxFlagColWidth_
	// Adjust padding if we get an alias
	if string(fl.Alias) != "" {
		fPrintNoSpace(
			f,
			string(fl.Alias),
			flagAliasNameSep,
		)
		paddingWidth = paddingWidth - len(fl.Alias) - len(flagAliasNameSep)
	}
	fPrintNoSpace(
		f,
		string(flagName),
	)
	if fl.SetBy != "" {
		fPrintNoSpace(
			f,
			flagNameValueSep,
			fl.Value.String(),
		)
		paddingWidth = paddingWidth - len(flagNameValueSep) - len(fl.Value.String())
	}
	fPrintlnNoSpace(
		f,
		padding(flagName, paddingWidth),
		flagHelpSep,
		string(fl.HelpShort),
		" (",
		fl.TypeDescription,
		")",
	)
}

func NoWrapCommandHelp(file *os.File, cur *command.Command, hi common.HelpInfo) command.Action {
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

		maxFlagColWidth_ := maxFlagColWidth(hi.AvailableFlags)
		// fmt.Println(flagIndent + strings.Repeat("_", maxFlagColWidth_))
		for _, flagName := range hi.AvailableFlags.SortedNames() {
			fl := hi.AvailableFlags[flagName]
			if fl.IsCommandFlag {
				printFlag(&commandFlagHelp, &col, flagName, &fl, maxFlagColWidth_)
			} else {
				printFlag(&sectionFlagHelp, &col, flagName, &fl, maxFlagColWidth_)
			}

		}
		if commandFlagHelp.Len() > 0 {
			fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Command Flags"))
			fmt.Fprintln(f)
			commandFlagHelp.WriteTo(f)
			fmt.Fprintln(f)

		}
		if sectionFlagHelp.Len() > 0 {
			fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Inherited Section Flags"))
			fmt.Fprintln(f)
			sectionFlagHelp.WriteTo(f)
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
