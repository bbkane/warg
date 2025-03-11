package warg_test

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/detailed"
	"go.bbkane.com/warg/section"
)

func exampleOverrideHelpFlaglogin(_ cli.Context) error {
	fmt.Println("Logging in")
	return nil
}

func exampleOverrideHelpFlagCustomCommandHelp(_ *cli.Command, _ cli.HelpInfo) cli.Action {
	return func(ctx cli.Context) error {
		file := ctx.Stdout
		fmt.Fprintln(file, "Custom command help")
		return nil
	}
}

func exampleOverrideHelpFlagCustomSectionHelp(_ *cli.SectionT, _ cli.HelpInfo) cli.Action {
	return func(ctx cli.Context) error {
		file := ctx.Stdout
		fmt.Fprintln(file, "Custom section help")
		return nil
	}
}

func ExampleOverrideHelpFlag() {
	app := warg.NewApp(
		"newAppName",
		"v1.0.0",
		section.NewSectionT(
			"work with a fictional blog platform",
			section.NewCommand(
				"login",
				"Login to the platform",
				exampleOverrideHelpFlaglogin,
			),
		),
		warg.OverrideHelpFlag(
			[]cli.HelpFlagMapping{
				{
					Name:        "default",
					CommandHelp: detailed.DetailedCommandHelp,
					SectionHelp: detailed.DetailedSectionHelp,
				},
				{
					Name:        "custom",
					CommandHelp: exampleOverrideHelpFlagCustomCommandHelp,
					SectionHelp: exampleOverrideHelpFlagCustomSectionHelp,
				},
			},
			"default",
			"--help",
			"Print help",
			flag.Alias("-h"),
			// the flag default should match a name in the HelpFlagMapping
		),
	)

	app.MustRun(cli.OverrideArgs([]string{"blog.exe", "-h", "custom"}))
	// Output:
	// Custom section help
}
