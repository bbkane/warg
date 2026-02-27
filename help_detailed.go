package warg

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"

	"go.bbkane.com/warg/styles"
	"go.bbkane.com/warg/value"
)

func detailedPrintFlag(p *styles.Printer, s *styles.Styles, name string, f *Flag, val value.Value) {
	if f.Alias != "" {
		p.Printf(
			"  %s , %s : %s\n",
			s.FlagName(name),
			s.FlagAlias(f.Alias),
			f.HelpShort,
		)
	} else {
		p.Printf(
			"  %s : %s\n",
			s.FlagName(name),
			f.HelpShort,
		)
	}
	p.Printf(
		"    %s : %s\n",
		s.Label("type"),
		val.Description(),
	)

	if len(val.Choices()) > 0 {
		p.Printf(
			"    %s : %s\n",
			s.Label("choices"),
			val.Choices(),
		)
	}

	if val.HasDefault() {
		switch v := val.(type) {
		case value.DictValue:
			p.Printf(
				"    %s\n",
				s.Label("default"),
			)
			def := v.DefaultStringMap()
			for _, key := range sortedKeys(def) {
				p.Printf(
					"      %s : %s\n",
					s.Label(key),
					def[key],
				)
			}
		case value.ScalarValue:
			p.Printf(
				"    %s : %s\n",
				s.Label("default"),
				v.DefaultString(),
			)
		case value.SliceValue:
			p.Printf(
				"    %s : %s\n",
				s.Label("default"),
				v.DefaultStringSlice(),
			)
		default:
			panic(fmt.Sprintf("Unexpected type: %#v", val))
		}
	}
	if f.ConfigPath != "" {
		p.Printf(
			"    %s : %s\n",
			s.Label("configpath"),
			f.ConfigPath,
		)
	}
	if len(f.EnvVars) > 0 {
		p.Printf(
			"    %s : %s\n",
			s.Label("envvars"),
			f.EnvVars,
		)
	}

	// TODO: it would be nice if this were red when the value isn't set
	if f.Required {
		p.Printf(
			"    %s : true\n",
			s.Label("required"),
		)
	}
	if f.UnsetSentinel != nil {
		p.Printf(
			"    %s : %s\n",
			s.Label("unsetsentinel"),
			*f.UnsetSentinel,
		)
	}

	if val.UpdatedBy() != value.UpdatedByUnset {
		switch v := val.(type) {
		case value.DictValue:
			p.Printf(
				"    %s (set by %s):\n",
				s.Label("currentvalue"),
				s.Label(string(val.UpdatedBy())),
			)
			m := v.StringMap()
			for _, key := range sortedKeys(m) {
				p.Printf(
					"      %s : %s\n",
					s.Label(key),
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

			p.Printf(
				"    %s (set by %s) :\n",
				s.Label("currentvalue"),
				s.Label(string(val.UpdatedBy())),
			)

			for i, e := range v.StringSlice() {
				indexStr := fmt.Sprint(i)
				padding := maxPadding - len(indexStr)
				p.Printf(
					"      %s %s\n",
					s.Label(leftPad(indexStr, "0", padding)+")"),
					e,
				)
			}
		case value.ScalarValue:
			p.Printf(
				"    %s (set by %s) : %s\n",
				s.Label("currentvalue"),
				s.Label(string(val.UpdatedBy())),
				v.String(),
			)
		default:
			panic(fmt.Sprintf("unexpected value: %#v", f))
		}
	}

	p.Println()
}

// detailedCmdHelp returns an Action that prints detailed help for the current command (cmdCtx.ParseState.CurrentCmd). It is expected to not be nil.
func detailedCmdHelp() Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		s, err := conditionallyEnableStyle(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		p := styles.NewPrinter(f)

		cur := cmdCtx.ParseState.CurrentCmd

		// Print top help section
		if cur.HelpLong != "" {
			p.Println(cur.HelpLong)
		} else {
			p.Println(cur.HelpShort)
		}

		p.Println()

		// compute sections for command flags and inherited flags,
		// then print their headers and them if they're not empty
		var commandFlagHelp bytes.Buffer
		var sectionFlagHelp bytes.Buffer
		{
			for _, name := range cmdCtx.App.GlobalFlags.SortedNames() {
				f := cmdCtx.App.GlobalFlags[name]
				val := cmdCtx.ParseState.FlagValues[name]
				detailedPrintFlag(styles.NewPrinter(&sectionFlagHelp), &s, name, &f, val)
			}

			for _, name := range cmdCtx.ParseState.CurrentCmd.Flags.SortedNames() {
				f := cmdCtx.ParseState.CurrentCmd.Flags[name]
				val := cmdCtx.ParseState.FlagValues[name]
				detailedPrintFlag(styles.NewPrinter(&commandFlagHelp), &s, name, &f, val)
			}

			if commandFlagHelp.Len() > 0 {
				p.Println(s.Header("Command Flags") + ":")
				p.Println()
				_, _ = commandFlagHelp.WriteTo(f)
			}
			if sectionFlagHelp.Len() > 0 {
				p.Println(s.Header("Global Flags") + ":")
				p.Println()
				_, _ = sectionFlagHelp.WriteTo(f)
			}
		}
		if cur.AllowForwardedArgs {
			p.Println(s.Header("Forwarded Arguments") + ":")
			p.Println()
			p.Println("  This command accepts forwarded arguments after a `--` separator.")
		}
		if cur.Footer != "" {
			p.Println()
			p.Println(s.Header("Footer") + ":")
			p.Println()
			p.Println(cur.Footer)
		}
		return nil

	}
}

func detailedSectionHelp() Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout

		f := bufio.NewWriter(file)
		defer f.Flush()

		s, err := conditionallyEnableStyle(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		p := styles.NewPrinter(f)

		cur := cmdCtx.ParseState.CurrentSection

		// Print top help section
		if cur.HelpLong != "" {
			p.Println(cur.HelpLong)
		} else {
			p.Println(cur.HelpShort)
		}

		p.Println()

		// Print sections
		if len(cur.Sections) > 0 {
			p.Println(s.Header("Sections") + ":")
			p.Println()

			for _, k := range cur.Sections.SortedNames() {
				p.Printf(
					"  %s : %s\n",
					s.SectionName(k),
					cur.Sections[k].HelpShort,
				)
			}

			p.Println()
		}

		// Print commands
		if len(cur.Cmds) > 0 {
			p.Println(s.Header("Commands") + ":")
			p.Println()

			for _, k := range cur.Cmds.SortedNames() {
				p.Printf(
					"  %s : %s\n",
					s.CommandName(k),
					cur.Cmds[string(k)].HelpShort,
				)
			}
		}

		if cur.Footer != "" {
			p.Println()
			p.Println(s.Header("Footer") + ":")
			p.Println()
			p.Println(cur.Footer)
		}
		return nil
	}
}
