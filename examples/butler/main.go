package main

import (
	"fmt"
	"os"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
)

func present(pf flag.PassedFlags) error {
	// this is a required flag, so we know it exists
	name := pf["--name"].(string)
	fmt.Printf("May I present to you %s.\n", name)
	return nil
}

func main() {
	app := warg.New(
		"butler",
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
	)
	app.MustRun(os.Args, os.LookupEnv)
}
