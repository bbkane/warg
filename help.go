package warg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"

	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
)

type CommandHelp = func(w io.Writer, appName string, path []string, cur c.Command, flagMap f.FlagMap) c.Action

type SectionHelp = func(w io.Writer, appName string, path []string, cur s.Section, flagMap f.FlagMap) c.Action

func printFlag(w io.Writer, name string, flag *f.Flag) {
	fmt.Fprintf(w, "  %s : %s\n", name, flag.Help)
	if flag.ConfigPath != "" {
		fmt.Fprintf(w, "    configpath : %s\n", flag.ConfigPath)
	}
	fmt.Fprintf(w, "    type : %s\n", flag.TypeDescription)
	if flag.SetBy != "" {
		fmt.Fprintf(w, "    value : %s\n", flag.Value)
		fmt.Fprintf(w, "    setby : %s\n", flag.SetBy)
	}
	fmt.Fprintln(w)
}

func DefaultCommandHelp(
	w io.Writer,
	appName string,
	path []string,
	cur c.Command,
	flagMap f.FlagMap,
) c.Action {
	return func(_ f.FlagValues) error {
		f := bufio.NewWriter(w)
		defer f.Flush()
		// Print top help section
		if cur.HelpLong == "" {
			fmt.Fprintf(f, "%s\n", cur.Help)
		} else {
			fmt.Fprintf(f, "%s\n", cur.Help)
		}

		fmt.Fprintln(f)

		// compute sections for command flags and inherited flags,
		// then print their headers and them if they're not empty
		var commandFlagHelp bytes.Buffer
		var sectionFlagHelp bytes.Buffer
		{
			keys := make([]string, 0, len(flagMap))
			for k := range flagMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, name := range keys {
				flag := flagMap[name]
				if flag.IsCommandFlag {
					printFlag(&commandFlagHelp, name, &flag)
				} else {
					printFlag(&sectionFlagHelp, name, &flag)
				}
			}

			if commandFlagHelp.Len() > 0 {
				fmt.Fprintf(f, "Command Flags:\n")
				fmt.Fprintln(f)
				commandFlagHelp.WriteTo(f)
			}
			if sectionFlagHelp.Len() > 0 {
				fmt.Fprintf(f, "Inherited Section Flags:\n")
				fmt.Fprintln(f)
				sectionFlagHelp.WriteTo(f)
			}
		}
		if cur.Footer != "" {
			fmt.Fprintf(f, "%s\n", cur.Footer)
		}
		return nil
	}
}

func DefaultSectionHelp(
	w io.Writer,
	appName string,
	path []string,
	cur s.Section,
	flagMap f.FlagMap,
) c.Action {
	return func(_ f.FlagValues) error {
		f := bufio.NewWriter(w)
		defer f.Flush()

		// Print top help section
		if cur.HelpLong == "" {
			fmt.Fprintf(f, "%s\n", cur.Help)
		} else {
			fmt.Fprintf(f, "%s\n", cur.Help)
		}

		fmt.Fprintln(f)

		// Print sections
		if len(cur.Sections) > 0 {
			fmt.Fprintf(f, "Sections:\n")
			fmt.Fprintln(f)
		}
		{
			keys := make([]string, 0, len(cur.Sections))
			for k := range cur.Sections {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(f, "  %s : %s\n", k, cur.Sections[k].Help)
			}
		}

		fmt.Fprintln(f)

		// Print commands
		if len(cur.Commands) > 0 {
			fmt.Fprintf(f, "Commands:\n")
			fmt.Fprintln(f)
		}
		{
			keys := make([]string, 0, len(cur.Commands))
			for k := range cur.Commands {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(f, "  %s : %s\n", k, cur.Commands[k].Help)
			}
		}
		if cur.Footer != "" {
			fmt.Fprintln(f)
			fmt.Fprintf(f, "%s\n", cur.Footer)
		}
		return nil
	}
}