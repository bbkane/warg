package warg_test

import (
	"fmt"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/parseopt"
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
		"--timeout": wargcore.NewFlag(
			"Optional timeout. Defaults to no timeout",
			scalar.Int(),
		),
		"--url": wargcore.NewFlag(
			"URL of the blog",
			scalar.String(
				scalar.Default("https://www.myblog.com"),
			),
			wargcore.EnvVars("BLOG_URL"),
		),
	}
	app := warg.New(
		"newAppName",
		"v1.0.0",
		wargcore.NewSection(
			"work with a fictional blog platform",
			wargcore.NewChildCmd(
				"login",
				"Login to the platform",
				login,
				wargcore.ChildFlagMap(commonFlags),
			),
			wargcore.NewChildSection(
				"comments",
				"Deal with comments",
				wargcore.NewChildCmd(
					"list",
					"List all comments",
					// still prototyping how we want this
					// command to look,
					// so use a provided stub action
					wargcore.DoNothing,
					wargcore.ChildFlagMap(commonFlags),
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
