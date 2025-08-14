package main

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func app() *wargcore.App {
	app := warg.New(
		"butler",
		"v1.0.0",
		wargcore.NewSection(
			string("A virtual assistant"),
			wargcore.NewChildCmd(
				"present",
				"Formally present a guest (guests are never introduced, always presented).",
				present,
				wargcore.NewChildFlag(
					"--name",
					"Guest to address.",
					scalar.String(),
					wargcore.Alias("-n"),
					wargcore.EnvVars("BUTLER_PRESENT_NAME", "USER"),
					wargcore.Required(),
				),
			),
		),
	)
	return &app
}

func present(ctx wargcore.Context) error {
	// this is a required flag, so we know it exists
	name := ctx.Flags["--name"].(string)
	fmt.Fprintf(ctx.Stdout, "May I present to you %s.\n", name)
	return nil
}

func main() {
	app := app()
	app.MustRun()
}
