package warg_test

import (
	"fmt"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/parseopt"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func login(ctx wargcore.Context) error {
	url := ctx.Flags["--url"].(string)

	// timeout doesn't have a default value,
	// so we can't rely on it being passed.
	timeout, exists := ctx.Flags["--timeout"]
	if exists {
		timeout := timeout.(int)
		fmt.Printf("Logging into %s with timeout %d\n", url, timeout)
		return nil
	}

	fmt.Printf("Logging into %s\n", url)
	return nil
}

func ExampleNew() {
	commonFlags := wargcore.FlagMap{
		"--timeout": flag.New(
			"Optional timeout. Defaults to no timeout",
			scalar.Int(),
		),
		"--url": flag.New(
			"URL of the blog",
			scalar.String(
				scalar.Default("https://www.myblog.com"),
			),
			flag.EnvVars("BLOG_URL"),
		),
	}
	app := warg.New(
		"newAppName",
		"v1.0.0",
		section.NewSection(
			"work with a fictional blog platform",
			section.NewChildCmd(
				"login",
				"Login to the platform",
				login,
				command.ChildFlagMap(commonFlags),
			),
			section.NewChildSection(
				"comments",
				"Deal with comments",
				section.NewChildCmd(
					"list",
					"List all comments",
					// still prototyping how we want this
					// command to look,
					// so use a provided stub action
					command.DoNothing,
					command.ChildFlagMap(commonFlags),
				),
			),
		),
	)

	// normally we would rely on the user to set the environment variable,
	// bu this is an example
	err := os.Setenv("BLOG_URL", "https://envvar.com")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app.MustRun(parseopt.Args([]string{"blog.exe", "login"}))
	// Output:
	// Logging into https://envvar.com
}
