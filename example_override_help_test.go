package warg_test

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/parseopt"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/wargcore"
)

func exampleHelpFlaglogin(_ wargcore.Context) error {
	fmt.Println("Logging in")
	return nil
}

func customHelpCmd() wargcore.Cmd {
	return command.NewCmd(
		"", // this command will be launched by the help flag, so users will never see the help
		func(ctx wargcore.Context) error {
			file := ctx.Stdout
			fmt.Fprintln(file, "Custom help command output")
			return nil
		},
	)
}

func ExampleHelpFlag() {

	// create a custom help command map by grabbing the default one
	// and adding our custom help command
	helpCommands := help.DefaultHelpCommandMap()
	helpCommands["custom"] = customHelpCmd()

	app := warg.New(
		"newAppName",
		"v1.0.0",
		section.NewSection(
			"work with a fictional blog platform",
			section.NewChildCmd(
				"login",
				"Login to the platform",
				exampleHelpFlaglogin,
			),
		),
		warg.HelpFlag(
			helpCommands,
			nil,
		),
	)

	app.MustRun(parseopt.Args([]string{"blog.exe", "-h", "custom"}))
	// Output:
	// Custom help command output
}
