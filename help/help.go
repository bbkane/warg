package help

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/bbkane/go-color"
	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
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
	AvailableFlags f.FlagMap
	// RootSection of the app. Especially useful for printing all sections and commands
	RootSection s.Section
}

type CommandHelp = func(file *os.File, cur c.Command, helpInfo HelpInfo) c.Action
type SectionHelp = func(file *os.File, cur s.Section, helpInfo HelpInfo) c.Action

// https://stackoverflow.com/a/45456649/2958070
func leftPad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func printFlag(w io.Writer, name string, flag *f.Flag) {
	if flag.Alias != "" {
		fmt.Fprintf(
			w,
			"  %s , %s : %s\n",
			color.Add(color.Bold+color.ForegroundYellow, name),
			color.Add(color.Bold+color.ForegroundYellow, flag.Alias),
			flag.Help,
		)
	} else {
		fmt.Fprintf(
			w,
			"  %s : %s\n",
			color.Add(color.Bold+color.ForegroundYellow, name),
			flag.Help,
		)
	}
	fmt.Fprintf(
		w,
		"    %s : %s\n",
		color.Add(color.Bold, "type"),
		flag.TypeDescription,
	)

	// TODO: should I print these one by one like I do value?
	if len(flag.DefaultValues) > 0 {
		if flag.TypeInfo == value.TypeInfoScalar {
			fmt.Fprintf(
				w,
				"    %s : %s\n",
				color.Add(color.Bold, "default"),
				flag.DefaultValues[0],
			)
		} else {
			fmt.Fprintf(
				w,
				"    %s : %s\n",
				color.Add(color.Bold, "default"),
				flag.DefaultValues,
			)
		}
	}
	if flag.ConfigPath != "" {
		fmt.Fprintf(
			w,
			"    %s : %s\n",
			color.Add(color.Bold, "configpath"),
			flag.ConfigPath,
		)
	}
	if len(flag.EnvVars) > 0 {
		fmt.Fprintf(w,
			"    %s : %s\n",
			color.Add(color.Bold, "envvars"),
			flag.EnvVars,
		)
	}

	// TODO: it would be nice if this were red when the value isn't set
	if flag.Required {
		fmt.Fprintf(w,
			"    %s : true\n",
			color.Add(color.Bold, "required"),
		)
	}

	if flag.SetBy != "" {
		if flag.TypeInfo == value.TypeInfoSlice {

			width := len(fmt.Sprint(len(flag.Value.StringSlice())))
			fmt.Fprintf(w,
				"    %s (set by %s) :\n",
				color.Add(color.Bold, "value"),
				color.Add(color.Bold, flag.SetBy),
			)

			for i, e := range flag.Value.StringSlice() {
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
				color.Add(color.Bold, flag.SetBy),
				flag.Value,
			)
		}
	}

	fmt.Fprintln(w)
}

func DefaultCommandHelp(file *os.File, cur c.Command, helpInfo HelpInfo) c.Action {
	return func(pf f.PassedFlags) error {
		f := bufio.NewWriter(file)
		defer f.Flush()

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
			keys := make([]string, 0, len(helpInfo.AvailableFlags))
			for k := range helpInfo.AvailableFlags {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, name := range keys {
				flag := helpInfo.AvailableFlags[name]
				if flag.IsCommandFlag {
					printFlag(&commandFlagHelp, name, &flag)
				} else {
					printFlag(&sectionFlagHelp, name, &flag)
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

func DefaultSectionHelp(file *os.File, cur s.Section, _ HelpInfo) c.Action {
	return func(pf f.PassedFlags) error {

		f := bufio.NewWriter(file)
		defer f.Flush()

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
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(
					f,
					"  %s : %s\n",
					color.Add(color.Bold+color.ForegroundCyan, k),
					cur.Sections[k].Help,
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
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(
					f,
					"  %s : %s\n",
					color.Add(color.Bold+color.ForegroundGreen, k),
					cur.Commands[k].Help,
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
