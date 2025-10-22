package main

import (
	"fmt"
	"strings"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value/scalar"
)

func app() *warg.App {
	app := warg.New(
		"butler",
		"v1.0.0",
		warg.NewSection(
			string("A virtual assistant"),
			warg.NewSubCmd(
				"present",
				"Formally present a guest (guests are never introduced, always presented).",
				present,
				warg.NewCmdFlag(
					"--name",
					"Guest to address.",
					scalar.String(),
					warg.Alias("-n"),
					warg.EnvVars("BUTLER_PRESENT_NAME", "USER"),
					warg.Required(),
				),
				warg.AllowForwardedArgs(),
				warg.CmdFooter("Treat butler well, and he will serve you faithfully."),
			),
		),
	)
	return &app
}

func present(ctx warg.CmdContext) error {
	// this is a required flag, so we know it exists
	name := ctx.Flags["--name"].(string)
	fmt.Fprintf(ctx.Stdout, "May I present to you %s.\n", name)
	if len(ctx.ForwardedArgs) > 0 {
		fmt.Fprintf(ctx.Stdout, "And: %s\n", strings.Join(ctx.ForwardedArgs, " "))
	}
	return nil
}

func main() {
	app := app()
	app.MustRun()
}
