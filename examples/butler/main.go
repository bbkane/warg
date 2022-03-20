package main

import (
	"fmt"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

func buildApp() warg.App {
	app := warg.New(
		section.New(
			section.HelpShort("A virtual assistant"),
			section.Command(
				command.Name("present"),
				command.HelpShort("Formally present a guest (guests are never introduced, always presented)."),
				present,
				command.Flag(
					flag.Name("--name"),
					flag.HelpShort("Guest to address."),
					value.String,
					flag.Alias("-n"),
					flag.EnvVars("BUTLER_PRESENT_NAME", "USER"),
					flag.Required(),
				),
			),
		),
		// Run the validation in a test instead of every
		// time the app is created.
		warg.SkipValidation(),
	)
	return app
}

func present(pf flag.PassedFlags) error {
	// this is a required flag, so we know it exists
	name := pf["--name"].(string)
	fmt.Printf("May I present to you %s.\n", name)
	return nil
}

func main() {
	app := buildApp()
	app.MustRun(os.Args, os.LookupEnv)
}
