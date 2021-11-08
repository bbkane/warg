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

func printFlag(w io.Writer, name string, flag *f.Flag) {
	fmt.Fprintf(w, "  %s : %s\n", name, flag.Help)
	fmt.Fprintf(w, "    type : %s\n", flag.TypeDescription)
	if flag.SetBy != "" {
		fmt.Fprintf(w, "    value : %s\n", flag.Value)
		fmt.Fprintf(w, "    setby : %s\n", flag.SetBy)
	}
	if len(flag.DefaultValues) > 0 {
		fmt.Fprintf(w, "    default : %s\n", flag.DefaultValues)
	}
	if flag.ConfigPath != "" {
		fmt.Fprintf(w, "    configpath : %s\n", flag.ConfigPath)
	}
	if len(flag.EnvVars) > 0 {
		fmt.Fprintf(w, "    envvars : %s\n", flag.EnvVars)
	}
	if flag.Required {
		fmt.Fprintf(w, "    required : true\n")
	}
	fmt.Fprintln(w)
}

func DefaultCommandHelp(file *os.File, cur c.Command, helpInfo HelpInfo) c.Action {
	return func(_ f.PassedFlags) error {
		f := bufio.NewWriter(file)
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
		if cur.HelpLong == "" {
			fmt.Fprintf(f, "%s\n", cur.Help)
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
					color.Add(color.ForegroundCyan, k),
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
					color.Add(color.ForegroundGreen, k),
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
