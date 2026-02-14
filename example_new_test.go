package warg_test

import (
	"fmt"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value/scalar"
)

func login(ctx warg.CmdContext) error {
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
	commonFlags := warg.FlagMap{
		"--timeout": warg.NewFlag(
			"Optional timeout. Defaults to no timeout",
			scalar.Int(),
		),
		"--url": warg.NewFlag(
			"URL of the blog",
			scalar.String(
				scalar.Default("https://www.myblog.com"),
			),
			warg.EnvVars("BLOG_URL"),
		),
	}
	app := warg.New(
		"newAppName",
		"v1.0.0",
		warg.NewSection(
			"work with a fictional blog platform",
			warg.NewSubCmd(
				"login",
				"Login to the platform",
				login,
				warg.CmdFlagMap(commonFlags),
			),
			warg.NewSubSection(
				"comments",
				"Deal with comments",
				warg.NewSubCmd(
					"list",
					"List all comments",
					// still prototyping how we want this
					// command to look,
					// so use a provided stub action
					warg.Unimplemented(),
					warg.CmdFlagMap(commonFlags),
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
	app.Run([]string{"login"})
	// Output:
	// Logging into https://envvar.com
}
