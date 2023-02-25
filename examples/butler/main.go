package main

import (
	"fmt"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
)

func app() *warg.App {
	app := warg.New(
		"butler",
		section.New(
			section.HelpShort("A virtual assistant"),
			section.Command(
				"present",
				"Formally present a guest (guests are never introduced, always presented).",
				present,
				command.Flag(
					"--name",
					"Guest to address.",
					scalar.String(),
					flag.Alias("-n"),
					flag.EnvVars("BUTLER_PRESENT_NAME", "USER"),
					flag.Required(),
				),
			),
		),
		warg.AddVersionCommand(""),
		warg.AddColorFlag(),
		warg.SkipValidation(),
	)
	return &app
}

func present(ctx command.Context) error {
	// this is a required flag, so we know it exists
	name := ctx.Flags["--name"].(string)
	fmt.Fprintf(ctx.Stdout, "May I present to you %s.\n", name)
	return nil
}

func main() {
	app := app()
	app.MustRun(os.Args, os.LookupEnv)
}
