package help

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/bbkane/gocolor"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
)

func outlineFlagHelper(w io.Writer, color *gocolor.Color, flagName string, f flag.Flag, indent int) {
	str := fmtFlagName(color, string(flagName))
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
	{
		flagKeys := make([]string, 0, len(sec.Flags))
		for k := range sec.Flags {
			flagKeys = append(flagKeys, string(k))
		}
		sort.Strings(flagKeys)
		for _, k := range flagKeys {
			outlineFlagHelper(w, color, k, sec.Flags[flag.Name(k)], indent)
		}
	}

	// commands and command flags
	{
		comKeys := make([]string, 0, len(sec.Commands))
		for comName := range sec.Commands {
			comKeys = append(comKeys, string(comName))
		}
		sort.Strings(comKeys)
		for _, comName := range comKeys {
			com := sec.Commands[command.Name(comName)]
			fmt.Fprintln(w, leftPad("# "+string(com.HelpShort), "  ", indent))
			fmt.Fprintln(
				w,
				leftPad(fmtCommandName(color, comName), "  ", indent),
			)
			// command flags
			flagKeys := make([]string, 0, len(com.Flags))
			for flagName := range com.Flags {
				flagKeys = append(flagKeys, string(flagName))
			}
			sort.Strings(flagKeys)
			for _, flagName := range flagKeys {
				outlineFlagHelper(w, color, flagName, com.Flags[flag.Name(flagName)], indent+1)
			}

		}
	}

	// sections
	{
		keys := make([]string, 0, len(sec.Sections))
		for k := range sec.Sections {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		for _, k := range keys {
			childSec := sec.Sections[section.Name(k)]
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
		fmt.Fprintf(f, "%s\n", fmtSectionName(&col, hi.AppName))

		outlineHelper(f, &col, hi.RootSection, 1)

		return nil
	}
}

func OutlineCommandHelp(file *os.File, cur *command.Command, helpInfo HelpInfo) command.Action {
	return OutlineSectionHelp(file, nil, helpInfo)
}
