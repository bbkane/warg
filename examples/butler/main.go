package main

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
)

func app() *cli.App {
	app := warg.NewApp(
		"butler",
		"v1.0.0",
		section.New(
			string("A virtual assistant"),
			section.NewCommand(
				"present",
				"Formally present a guest (guests are never introduced, always presented).",
				present,
				command.NewFlag(
					"--name",
					"Guest to address.",
					scalar.String(),
					flag.Alias("-n"),
					flag.EnvVars("BUTLER_PRESENT_NAME", "USER"),
					flag.Required(),
				),
			),
			section.CommandMap(warg.VersionCommandMap()),
		),
		warg.GlobalFlagMap(warg.ColorFlagMap()),
	)
	return &app
}

func present(ctx cli.Context) error {
	// this is a required flag, so we know it exists
	name := ctx.Flags["--name"].(string)
	fmt.Fprintf(ctx.Stdout, "May I present to you %s.\n", name)
	return nil
}

func main() {
	app := app()
	app.MustRun()
}
