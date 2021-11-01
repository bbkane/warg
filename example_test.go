package warg_test

import (
	"fmt"
	"os"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
)

func login(pf flag.PassedFlags) error {
	url := pf["--url"].(string)

	// timeout doesn't have a default value,
	// so we can't rely on it being passed.
	timeout, exists := pf["--timeout"]
	if exists {
		timeout := timeout.(int)
		fmt.Printf("Logging into %s with timeout %d\n", url, timeout)
		return nil
	}

	fmt.Printf("Logging into %s\n", url)
	return nil
}

func ExampleNew() {
	app := warg.New(
		"blog",
		section.New(
			"work with a fictional blog platform",
			section.WithCommand(
				"login",
				"Login to the platform",
				login,
			),
			section.WithFlag(
				"--timeout",
				"Optional timeout. Defaults to no timeout",
				value.Int,
			),
			section.WithFlag(
				"--url",
				"URL of the blog",
				value.String,
				flag.Default("https://www.myblog.com"),
			),
			section.WithSection(
				"comments",
				"Deal with comments",
				section.WithCommand(
					"list",
					"List all comments",
					// still prototyping how we want this
					// command to look,
					// so use a provided stub action
					command.DoNothing,
				),
			),
		),
	)

	// Of course, in actual code, it would be something like:
	// err := app.Run(os.Args, os.LookupEnv)
	// TODO: add envvar to the example once they're implemented - looks like the Go playground supports them ( https://play.golang.org/p/nHJQcAUewNF )
	err := app.Run([]string{"blog.exe", "login"}, warg.LookupDict(nil))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
