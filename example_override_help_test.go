package warg_test

import (
	"fmt"
	"os"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/help"
	"github.com/bbkane/warg/section"
)

func exampleOverrideHelpFlaglogin(pf flag.PassedFlags) error {
	fmt.Println("Logging in")
	return nil
}

func exampleOverrideHelpFlagCustomCommandHelp(file *os.File, _ command.Command, _ help.HelpInfo) command.Action {
	return func(_ flag.PassedFlags) error {
		fmt.Fprintln(file, "Custom command help")
		return nil
	}
}

func exampleOverrideHelpFlagCustomSectionHelp(file *os.File, _ section.SectionT, _ help.HelpInfo) command.Action {
	return func(_ flag.PassedFlags) error {
		fmt.Fprintln(file, "Custom section help")
		return nil
	}
}

func ExampleOverrideHelpFlag() {
	app := warg.New(
		"blog",
		section.New(
			"work with a fictional blog platform",
			section.Command(
				"login",
				"Login to the platform",
				exampleOverrideHelpFlaglogin,
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

	app.MustRun([]string{"blog.exe", "-h", "custom"}, os.LookupEnv)
	// Output:
	// Custom section help
}
