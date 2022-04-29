package warg_test

import (
	"fmt"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/section"
)

func exampleOverrideHelpFlaglogin(_ command.Context) error {
	fmt.Println("Logging in")
	return nil
}

func exampleOverrideHelpFlagCustomCommandHelp(file *os.File, _ *command.Command, _ help.HelpInfo) command.Action {
	return func(_ command.Context) error {
		fmt.Fprintln(file, "Custom command help")
		return nil
	}
}

func exampleOverrideHelpFlagCustomSectionHelp(file *os.File, _ *section.SectionT, _ help.HelpInfo) command.Action {
	return func(_ command.Context) error {
		fmt.Fprintln(file, "Custom section help")
		return nil
	}
}

func ExampleOverrideHelpFlag() {
	app := warg.New(
		"newAppName",
		section.New(
			"work with a fictional blog platform",
			section.Command(
				"login",
				"Login to the platform",
				exampleOverrideHelpFlaglogin,
			),
		),
		warg.OverrideHelpFlag(
			[]help.HelpFlagMapping{
				{
					Name:        "default",
					CommandHelp: help.DetailedCommandHelp,
					SectionHelp: help.DetailedSectionHelp,
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
