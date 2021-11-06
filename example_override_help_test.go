package warg_test

import (
	"fmt"
	"io"
	"os"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/help"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
)

func exampleOverrideHelpFlaglogin(pf flag.PassedFlags) error {
	url := pf["--url"].(string)

	// timeout doesn't have a default value,
	// so we can't rely on it being passed.
	timeout, exists := pf["--timeout"]
	if exists {
		timeout := timeout.(int)
		fmt.Printf("Logging into %s with timeout %d\n", url, timeout)
		return nil
	}

	fmt.Printf("Logging into %s\n", url)
	return nil
}

func exampleOverrideHelpFlagCustomCommandHelp(w io.Writer, _ command.Command, _ help.HelpInfo) command.Action {
	return func(_ flag.PassedFlags) error {
		fmt.Fprintln(w, "Custom command help")
		return nil
	}
}

func exampleOverrideHelpFlagCustomSectionHelp(w io.Writer, _ section.Section, _ help.HelpInfo) command.Action {
	return func(_ flag.PassedFlags) error {
		fmt.Fprintln(w, "Custom section help")
		return nil
	}
}

func ExampleOverrideHelpFlag() {
	app := warg.New(
		"blog",
		section.New(
			"work with a fictional blog platform",
			section.WithCommand(
				"login",
				"Login to the platform",
				exampleOverrideHelpFlaglogin,
			),
			section.WithFlag(
				"--timeout",
				"Optional timeout. Defaults to no timeout",
				value.Int,
			),
			section.WithFlag(
				"--url",
				"URL of the blog",
				value.String,
				flag.Default("https://www.myblog.com"),
				flag.EnvVars("BLOG_URL"),
			),
			section.WithSection(
				"comments",
				"Deal with comments",
				section.WithCommand(
					"list",
					"List all comments",
					// still prototyping how we want this
					// command to look,
					// so use a provided stub action
					command.DoNothing,
				),
			),
		),
		warg.OverrideHelpFlag(
			[]warg.HelpFlagMapping{
				{
					Name:        "default",
					CommandHelp: help.DefaultCommandHelp,
					SectionHelp: help.DefaultSectionHelp,
				},
				{
					Name:        "custom",
					CommandHelp: exampleOverrideHelpFlagCustomCommandHelp,
					SectionHelp: exampleOverrideHelpFlagCustomSectionHelp,
				},
			},
			os.Stdout,
			"--help",
			"Print help",
			flag.Alias("-h"),
			// the flag default should match a name in the HelpFlagMapping
			flag.Default("default"),
		),
	)

	err := app.Run([]string{"blog.exe", "-h", "custom"}, os.LookupEnv)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	// Output:
	// Custom section help
}
