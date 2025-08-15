package warg_test

import (
	"fmt"

	"go.bbkane.com/warg"
)

func exampleHelpFlaglogin(_ warg.CmdContext) error {
	fmt.Println("Logging in")
	return nil
}

func customHelpCmd() warg.Cmd {
	return warg.NewCmd(
		"", // this command will be launched by the help flag, so users will never see the help
		func(ctx warg.CmdContext) error {
			file := ctx.Stdout
			fmt.Fprintln(file, "Custom help command output")
			return nil
		},
	)
}

func ExampleHelpFlag() {

	// create a custom help command map by grabbing the default one
	// and adding our custom help command
	helpCommands := warg.DefaultHelpCommandMap()
	helpCommands["custom"] = customHelpCmd()

	app := warg.New(
		"newAppName",
		"v1.0.0",
		warg.NewSection(
			"work with a fictional blog platform",
			warg.NewSubCmd(
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

	app.MustRun(warg.ParseWithArgs([]string{"blog.exe", "-h", "custom"}))
	// Output:
	// Custom help command output
}
