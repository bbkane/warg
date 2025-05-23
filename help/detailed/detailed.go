package detailed

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/help/common"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/wargcore"
)

func detailedPrintFlag(w io.Writer, color *gocolor.Color, name string, f *wargcore.Flag) {
	if f.Alias != "" {
		fmt.Fprintf(
			w,
			"  %s , %s : %s\n",
			common.FmtFlagName(color, name),
			common.FmtFlagAlias(color, f.Alias),
			f.HelpShort,
		)
	} else {
		fmt.Fprintf(
			w,
			"  %s : %s\n",
			common.FmtFlagName(color, name),
			f.HelpShort,
		)
	}
	fmt.Fprintf(
		w,
		"    %s : %s\n",
		color.Add(color.Bold, "type"),
		f.Value.Description(),
	)

	if len(f.Value.Choices()) > 0 {
		fmt.Fprintf(w,
			"    %s : %s\n",
			color.Add(color.Bold, "choices"),
			f.Value.Choices(),
		)
	}

	if f.Value.HasDefault() {
		switch v := f.Value.(type) {
		case value.DictValue:
			fmt.Fprintf(
				w,
				"    %s\n",
				color.Add(color.Bold, "default"),
			)
			def := v.DefaultStringMap()
			for _, key := range common.SortedKeys(def) {
				fmt.Fprintf(
					w,
					"      %s : %s\n",
					color.Add(color.Bold, key),
					def[key],
				)
			}
		case value.ScalarValue:
			fmt.Fprintf(
				w,
				"    %s : %s\n",
				color.Add(color.Bold, "default"),
				v.DefaultString(),
			)
		case value.SliceValue:
			fmt.Fprintf(
				w,
				"    %s : %s\n",
				color.Add(color.Bold, "default"),
				v.DefaultStringSlice(),
			)
		default:
			panic(fmt.Sprintf("Unexpected type: %#v", f.Value))
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
	if f.UnsetSentinel != "" {
		fmt.Fprintf(
			w,
			"    %s : %s\n",
			color.Add(color.Bold, "unsetsentinel"),
			f.UnsetSentinel,
		)
	}

	if f.Value.UpdatedBy() != value.UpdatedByUnset {
		switch v := f.Value.(type) {
		case value.DictValue:
			fmt.Fprintf(
				w,
				"    %s (set by %s):\n",
				color.Add(color.Bold, "currentvalue"),
				color.Add(color.Bold, string(f.Value.UpdatedBy())),
			)
			m := v.StringMap()
			for _, key := range common.SortedKeys(m) {
				fmt.Fprintf(
					w,
					"      %s : %s\n",
					color.Add(color.Bold, key),
					m[key],
				)
			}
		case value.SliceValue:
			sliceLen := len(fmt.Sprint(len(v.StringSlice())))

			// Find the max padding for a specified length
			// 0 - 9 : 0  # no padding needed
			// 10 - 99 : 1  # need 0 for single digit numbers
			//  100 - 999 : 2
			maxPadding := int(math.Ceil(math.Log10(float64(sliceLen)))) + 1

			fmt.Fprintf(w,
				"    %s (set by %s) :\n",
				color.Add(color.Bold, "currentvalue"),
				color.Add(color.Bold, string(f.Value.UpdatedBy())),
			)

			for i, e := range v.StringSlice() {
				indexStr := fmt.Sprint(i)
				padding := maxPadding - len(indexStr)
				fmt.Fprintf(
					w,
					"      %s %s\n",
					color.Add(
						color.Bold,
						common.LeftPad(indexStr, "0", padding)+")",
					),
					e,
				)
			}
		case value.ScalarValue:
			fmt.Fprintf(
				w,
				"    %s (set by %s) : %s\n",
				color.Add(color.Bold, "currentvalue"),
				color.Add(color.Bold, string(f.Value.UpdatedBy())),
				v.String(),
			)
		default:
			panic(fmt.Sprintf("unexpected value: %#v", f))
		}
	}

	fmt.Fprintln(w)
}

func DetailedCommandHelp(cur *wargcore.Command, helpInfo wargcore.HelpInfo) wargcore.Action {
	return func(pf wargcore.Context) error {
		file := pf.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := common.ConditionallyEnableColor(pf.Flags, file)
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

			for _, name := range helpInfo.AvailableFlags.SortedNames() {
				f := helpInfo.AvailableFlags[name]
				if f.IsCommandFlag {
					detailedPrintFlag(&commandFlagHelp, &col, name, &f)
				} else {
					detailedPrintFlag(&sectionFlagHelp, &col, name, &f)
				}
			}

			if commandFlagHelp.Len() > 0 {
				fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Command Flags")+":")
				fmt.Fprintln(f)
				_, _ = commandFlagHelp.WriteTo(f)
			}
			if sectionFlagHelp.Len() > 0 {
				fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Inherited Section Flags")+":")
				fmt.Fprintln(f)
				_, _ = sectionFlagHelp.WriteTo(f)
			}
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

func DetailedSectionHelp(cur *wargcore.Section, _ wargcore.HelpInfo) wargcore.Action {
	return func(pf wargcore.Context) error {
		file := pf.Stdout

		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := common.ConditionallyEnableColor(pf.Flags, file)
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

			for _, k := range cur.Sections.SortedNames() {
				fmt.Fprintf(
					f,
					"  %s : %s\n",
					common.FmtSectionName(&col, k),
					cur.Sections[k].HelpShort,
				)
			}

			fmt.Fprintln(f)
		}

		// Print commands
		if len(cur.Commands) > 0 {
			fmt.Fprintln(f, col.Add(col.Underline+col.Bold, "Commands")+":")
			fmt.Fprintln(f)

			for _, k := range cur.Commands.SortedNames() {
				fmt.Fprintf(
					f,
					"  %s : %s\n",
					common.FmtCommandName(&col, k),
					cur.Commands[string(k)].HelpShort,
				)
			}
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
