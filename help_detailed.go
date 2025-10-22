package warg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"

	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/value"
)

func detailedPrintFlag(w io.Writer, color *gocolor.Color, name string, f *Flag, val value.Value) {
	if f.Alias != "" {
		fmt.Fprintf(
			w,
			"  %s , %s : %s\n",
			fmtFlagName(color, name),
			fmtFlagAlias(color, f.Alias),
			f.HelpShort,
		)
	} else {
		fmt.Fprintf(
			w,
			"  %s : %s\n",
			fmtFlagName(color, name),
			f.HelpShort,
		)
	}
	fmt.Fprintf(
		w,
		"    %s : %s\n",
		color.Add(color.Bold, "type"),
		val.Description(),
	)

	if len(val.Choices()) > 0 {
		fmt.Fprintf(w,
			"    %s : %s\n",
			color.Add(color.Bold, "choices"),
			val.Choices(),
		)
	}

	if val.HasDefault() {
		switch v := val.(type) {
		case value.DictValue:
			fmt.Fprintf(
				w,
				"    %s\n",
				color.Add(color.Bold, "default"),
			)
			def := v.DefaultStringMap()
			for _, key := range sortedKeys(def) {
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
			panic(fmt.Sprintf("Unexpected type: %#v", val))
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
	if f.UnsetSentinel != nil {
		fmt.Fprintf(
			w,
			"    %s : %s\n",
			color.Add(color.Bold, "unsetsentinel"),
			*f.UnsetSentinel,
		)
	}

	if val.UpdatedBy() != value.UpdatedByUnset {
		switch v := val.(type) {
		case value.DictValue:
			fmt.Fprintf(
				w,
				"    %s (set by %s):\n",
				color.Add(color.Bold, "currentvalue"),
				color.Add(color.Bold, string(val.UpdatedBy())),
			)
			m := v.StringMap()
			for _, key := range sortedKeys(m) {
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
				color.Add(color.Bold, string(val.UpdatedBy())),
			)

			for i, e := range v.StringSlice() {
				indexStr := fmt.Sprint(i)
				padding := maxPadding - len(indexStr)
				fmt.Fprintf(
					w,
					"      %s %s\n",
					color.Add(
						color.Bold,
						leftPad(indexStr, "0", padding)+")",
					),
					e,
				)
			}
		case value.ScalarValue:
			fmt.Fprintf(
				w,
				"    %s (set by %s) : %s\n",
				color.Add(color.Bold, "currentvalue"),
				color.Add(color.Bold, string(val.UpdatedBy())),
				v.String(),
			)
		default:
			panic(fmt.Sprintf("unexpected value: %#v", f))
		}
	}

	fmt.Fprintln(w)
}

// detailedCmdHelp returns an Action that prints detailed help for the current command (cmdCtx.ParseState.CurrentCmd). It is expected to not be nil.
func detailedCmdHelp() Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		cur := cmdCtx.ParseState.CurrentCmd

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
			for _, name := range cmdCtx.App.GlobalFlags.SortedNames() {
				f := cmdCtx.App.GlobalFlags[name]
				val := cmdCtx.ParseState.FlagValues[name]
				detailedPrintFlag(&sectionFlagHelp, &col, name, &f, val)
			}

			for _, name := range cmdCtx.ParseState.CurrentCmd.Flags.SortedNames() {
				f := cmdCtx.ParseState.CurrentCmd.Flags[name]
				val := cmdCtx.ParseState.FlagValues[name]
				detailedPrintFlag(&commandFlagHelp, &col, name, &f, val)
			}

			if commandFlagHelp.Len() > 0 {
				fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Command Flags")+":")
				fmt.Fprintln(f)
				_, _ = commandFlagHelp.WriteTo(f)
			}
			if sectionFlagHelp.Len() > 0 {
				fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Global Flags")+":")
				fmt.Fprintln(f)
				_, _ = sectionFlagHelp.WriteTo(f)
			}
		}
		if cur.AllowForwardedArgs {
			fmt.Fprintln(f, col.Add(col.Bold+col.Underline, "Forwarded Arguments")+":")
			fmt.Fprintln(f)
			fmt.Fprintln(f, "  This command accepts forwarded arguments after a `--` separator.")
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

func detailedSectionHelp() Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout

		f := bufio.NewWriter(file)
		defer f.Flush()

		col, err := ConditionallyEnableColor(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		cur := cmdCtx.ParseState.CurrentSection

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
					fmtSectionName(&col, k),
					cur.Sections[k].HelpShort,
				)
			}

			fmt.Fprintln(f)
		}

		// Print commands
		if len(cur.Cmds) > 0 {
			fmt.Fprintln(f, col.Add(col.Underline+col.Bold, "Commands")+":")
			fmt.Fprintln(f)

			for _, k := range cur.Cmds.SortedNames() {
				fmt.Fprintf(
					f,
					"  %s : %s\n",
					fmtCommandName(&col, k),
					cur.Cmds[string(k)].HelpShort,
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
