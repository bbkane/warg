package help

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/bbkane/go-color"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
	"github.com/mattn/go-isatty"
)

// HelpInfo lists common information available to a help function
type HelpInfo struct {
	// AppName as defined by warg.New()
	AppName string
	// Path passed either to a command or a section
	Path []string
	// AvailableFlags for the section or commmand.
	// All flags are Resolved if possible (i.e., flag.SetBy != "")
	AvailableFlags flag.FlagMap
	// RootSection of the app. Especially useful for printing all sections and commands
	RootSection section.SectionT
}

type CommandHelp = func(file *os.File, cur command.Command, helpInfo HelpInfo) command.Action
type SectionHelp = func(file *os.File, cur section.SectionT, helpInfo HelpInfo) command.Action

// https://stackoverflow.com/a/45456649/2958070
func leftPad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func printFlag(w io.Writer, name string, f *flag.Flag) {
	if f.Alias != "" {
		fmt.Fprintf(
			w,
			"  %s , %s : %s\n",
			color.Add(color.Bold+color.ForegroundYellow, name),
			color.Add(color.Bold+color.ForegroundYellow, f.Alias),
			f.Help,
		)
	} else {
		fmt.Fprintf(
			w,
			"  %s : %s\n",
			color.Add(color.Bold+color.ForegroundYellow, name),
			f.Help,
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
				color.Add(color.Bold, "value"),
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
				color.Add(color.Bold, "value"),
				color.Add(color.Bold, f.SetBy),
				f.Value,
			)
		}
	}

	fmt.Fprintln(w)
}

// SetColor looks for a passed --color flag with an underlying string value. If
// it exists and is set to "true", color is enabled. If it exists, is set to
// "auto", and the passed file is a terminal, color is enabled
func ConditionallyEnableColor(pf flag.PassedFlags, file *os.File) {
	// default to trying to use color
	useColor := "auto"
	// respect a --color string
	if useColorI, exists := pf["--color"]; exists {
		if useColorUnder, isStr := useColorI.(string); isStr {
			useColor = useColorUnder
		}
	}

	if useColor == "true" || (useColor == "auto" && isatty.IsTerminal(file.Fd())) {
		color.Enable()
	}
}

func DefaultCommandHelp(file *os.File, cur command.Command, helpInfo HelpInfo) command.Action {
	return func(pf flag.PassedFlags) error {
		f := bufio.NewWriter(file)
		defer f.Flush()

		ConditionallyEnableColor(pf, file)

		// Print top help section
		if cur.HelpLong != "" {
			fmt.Fprintf(f, "%s\n", cur.HelpLong)
		} else {
			fmt.Fprintf(f, "%s\n", cur.Help)
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
					printFlag(&commandFlagHelp, name, &f)
				} else {
					printFlag(&sectionFlagHelp, name, &f)
				}
			}

			if commandFlagHelp.Len() > 0 {
				fmt.Fprintln(f, color.Add(color.Bold+color.Underline, "Command Flags"))
				fmt.Fprintln(f)
				commandFlagHelp.WriteTo(f)
			}
			if sectionFlagHelp.Len() > 0 {
				fmt.Fprintln(f, color.Add(color.Bold+color.Underline, "Inherited Section Flags"))
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

func DefaultSectionHelp(file *os.File, cur section.SectionT, _ HelpInfo) command.Action {
	return func(pf flag.PassedFlags) error {

		f := bufio.NewWriter(file)
		defer f.Flush()

		ConditionallyEnableColor(pf, file)

		// Print top help section
		if cur.HelpLong != "" {
			fmt.Fprintf(f, "%s\n", cur.HelpLong)
		} else {
			fmt.Fprintf(f, "%s\n", cur.Help)
		}

		fmt.Fprintln(f)

		// Print sections
		if len(cur.Sections) > 0 {
			fmt.Fprintln(f, color.Add(color.Underline+color.Bold, "Sections"))
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
					color.Add(color.Bold+color.ForegroundCyan, k),
					cur.Sections[section.Name(k)].Help,
				)
			}

			fmt.Fprintln(f)
		}

		// Print commands
		if len(cur.Commands) > 0 {
			fmt.Fprintln(f, color.Add(color.Underline+color.Bold, "Commands"))
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
					color.Add(color.Bold+color.ForegroundGreen, k),
					cur.Commands[command.Name(k)].Help,
				)
			}
		}

		if cur.Footer != "" {
			fmt.Fprintln(f)
			fmt.Fprintln(f, color.Add(color.Underline+color.Bold, "Footer"))
			fmt.Fprintln(f)
			fmt.Fprintf(f, "%s\n", cur.Footer)
		}
		return nil
	}
}
