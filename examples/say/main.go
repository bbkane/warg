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

func hello(pf flag.PassedFlags) error {
	// this is a required flag, so we know it exists
	name := pf["--name"].(string)
	fmt.Printf("Hello %s!\n", name)
	return nil
}

func main() {
	app := warg.New(
		"say",
		section.New(
			"Make the terminal say things!!",
			section.Command(
				"hello",
				"Say hello",
				hello,
				command.Flag(
					"--name",
					"Person we're talking to",
					value.String,
					flag.Alias("-n"),
					flag.EnvVars("SAY_NAME"),
					flag.Required(),
				),
			),
		),
	)
	app.MustRun(os.Args, os.LookupEnv)
}
