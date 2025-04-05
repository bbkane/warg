package warg_test

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/section"
)

func exampleHelpFlaglogin(_ cli.Context) error {
	fmt.Println("Logging in")
	return nil
}

func customHelpCmd() cli.Command {
	return command.New(
		"", // this command will be launched by the help flag, so users will never see the help
		func(ctx cli.Context) error {
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

	app := warg.NewApp(
		"newAppName",
		"v1.0.0",
		section.New(
			"work with a fictional blog platform",
			section.NewCommand(
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

	app.MustRun(cli.OverrideArgs([]string{"blog.exe", "-h", "custom"}))
	// Output:
	// Custom help command output
}
