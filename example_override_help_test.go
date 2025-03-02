package warg_test

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/help/common"
	"go.bbkane.com/warg/help/detailed"
	"go.bbkane.com/warg/section"
)

func exampleOverrideHelpFlaglogin(_ warg.CommandContext) error {
	fmt.Println("Logging in")
	return nil
}

func exampleOverrideHelpFlagCustomCommandHelp(_ *warg.Command, _ common.HelpInfo) warg.Action {
	return func(ctx warg.CommandContext) error {
		file := ctx.Stdout
		fmt.Fprintln(file, "Custom command help")
		return nil
	}
}

func exampleOverrideHelpFlagCustomSectionHelp(_ *section.SectionT, _ common.HelpInfo) warg.Action {
	return func(ctx warg.CommandContext) error {
		file := ctx.Stdout
		fmt.Fprintln(file, "Custom section help")
		return nil
	}
}

func ExampleOverrideHelpFlag() {
	app := warg.New(
		"newAppName",
		"v1.0.0",
		section.New(
			"work with a fictional blog platform",
			section.NewCommand(
				"login",
				"Login to the platform",
				exampleOverrideHelpFlaglogin,
			),
		),
		warg.OverrideHelpFlag(
			[]help.HelpFlagMapping{
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

	app.MustRun(warg.OverrideArgs([]string{"blog.exe", "-h", "custom"}))
	// Output:
	// Custom section help
}
