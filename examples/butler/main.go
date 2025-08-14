package main

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value/scalar"
)

func app() *warg.App {
	app := warg.New(
		"butler",
		"v1.0.0",
		warg.NewSection(
			string("A virtual assistant"),
			warg.NewChildCmd(
				"present",
				"Formally present a guest (guests are never introduced, always presented).",
				present,
				warg.NewChildFlag(
					"--name",
					"Guest to address.",
					scalar.String(),
					warg.Alias("-n"),
					warg.EnvVars("BUTLER_PRESENT_NAME", "USER"),
					warg.Required(),
				),
			),
		),
	)
	return &app
}

func present(ctx warg.Context) error {
	// this is a required flag, so we know it exists
	name := ctx.Flags["--name"].(string)
	fmt.Fprintf(ctx.Stdout, "May I present to you %s.\n", name)
	return nil
}

func main() {
	app := app()
	app.MustRun()
}
