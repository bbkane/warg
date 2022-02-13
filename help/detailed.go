package help

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/bbkane/gocolor"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
)

func detailedPrintFlag(w io.Writer, color *gocolor.Color, name string, f *flag.Flag) {
	if f.Alias != "" {
		fmt.Fprintf(
			w,
			"  %s , %s : %s\n",
			color.Add(color.Bold+color.FgYellow, name),
			color.Add(color.Bold+color.FgYellow, f.Alias),
			f.HelpShort,
		)
	} else {
		fmt.Fprintf(
			w,
			"  %s : %s\n",
			color.Add(color.Bold+color.FgYellow, name),
			f.HelpShort,
		)
	}
	fmt.Fprintf(
		w,
		"    %s : %s\n",
		color.Add(color.Bold, "type"),
		f.TypeDescription,
	)

	// TODO: should I print these one by one like I do value?
	if len(f.DefaultValues) > 0 {
		if f.TypeInfo == value.TypeInfoScalar {
			fmt.Fprintf(
				w,
				"    %s : %s\n",
				color.Add(color.Bold, "default"),
				f.DefaultValues[0],
			)
		} else {
			fmt.Fprintf(
				w,
				"    %s : %s\n",
				color.Add(color.Bold, "default"),
				f.DefaultValues,
			)
		}
	}
	if f.ConfigPath != "" {
		fmt.Fprintf(
			w,
			"    %s : %s\n",
			color.Add(color.Bold, "configpath"),
			f.ConfigPath,
		)
	}
	if len(f.EnvVars) > 0 {
		fmt.Fprintf(w,
			"    %s : %s\n",
			color.Add(color.Bold, "envvars"),
			f.EnvVars,
		)
	}

	// TODO: it would be nice if this were red when the value isn't set
	if f.Required {
		fmt.Fprintf(w,
			"    %s : true\n",
			color.Add(color.Bold, "required"),
		)
	}

	if f.SetBy != "" {
		if f.TypeInfo == value.TypeInfoSlice {

			width := len(fmt.Sprint(len(f.Value.StringSlice())))
			fmt.Fprintf(w,
				"    %s (set by %s) :\n",
				color.Add(color.Bold, "currentvalue"),
				color.Add(color.Bold, f.SetBy),
			)

			for i, e := range f.Value.StringSlice() {
				fmt.Fprintf(
					w,
					"      %s %s\n",
					color.Add(
						color.Bold,
						leftPad(fmt.Sprint(i), "0", width)+")",
					),
					e,
				)
			}
		} else {
			fmt.Fprintf(
				w,
				"    %s (set by %s) : %s\n",
				color.Add(color.Bold, "currentvalue"),
				color.Add(color.Bold, f.SetBy),
				f.Value,
			)
		}
	}

	fmt.Fprintln(w)
}

func DetailedCommandHelp(file *os.File, cur command.Command, helpInfo HelpInfo) command.Action {
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

		// compute sections for command flags and inherited flags,
		// then print their headers and them if they're not empty
		var commandFlagHelp bytes.Buffer
		var sectionFlagHelp bytes.Buffer
		{
			// we need to sort these things so we need to use strings here
			keys := make([]string, 0, len(helpInfo.AvailableFlags))
			for k := range helpInfo.AvailableFlags {
				keys = append(keys, string(k))
			}
			sort.Strings(keys)
			for _, name := range keys {
				f := helpInfo.AvailableFlags[flag.Name(name)]
				if f.IsCommandFlag {
					detailedPrintFlag(&commandFlagHelp, &col, name, &f)
				} else {
					detailedPrintFlag(&sectionFlagHelp, &col, name, &f)
				}
			}

			if commandFlagHelp.Len() > 0 {
				fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Command Flags"))
				fmt.Fprintln(f)
				commandFlagHelp.WriteTo(f)
			}
			if sectionFlagHelp.Len() > 0 {
				fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Inherited Section Flags"))
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

func DetailedSectionHelp(file *os.File, cur section.SectionT, _ HelpInfo) command.Action {
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

		// Print sections
		if len(cur.Sections) > 0 {
			fmt.Fprintln(f, col.Add(col.Underline+col.Bold, "Sections")+":")
			fmt.Fprintln(f)

			keys := make([]string, 0, len(cur.Sections))
			for k := range cur.Sections {
				keys = append(keys, string(k))
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(
					f,
					"  %s : %s\n",
					col.Add(col.Bold+col.FgCyan, k),
					cur.Sections[section.Name(k)].HelpShort,
				)
			}

			fmt.Fprintln(f)
		}

		// Print commands
		if len(cur.Commands) > 0 {
			fmt.Fprintln(f, col.Add(col.Underline+col.Bold, "Commands")+":")
			fmt.Fprintln(f)

			keys := make([]string, 0, len(cur.Commands))
			for k := range cur.Commands {
				keys = append(keys, string(k))
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(
					f,
					"  %s : %s\n",
					col.Add(col.Bold+col.FgGreen, k),
					cur.Commands[command.Name(k)].HelpShort,
				)
			}
		}

		if cur.Footer != "" {
			fmt.Fprintln(f)
			fmt.Fprintln(f, col.Add(col.Underline+col.Bold, "Footer"))
			fmt.Fprintln(f)
			fmt.Fprintf(f, "%s\n", cur.Footer)
		}
		return nil
	}
}
