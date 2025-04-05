package warg_test

import (
	"fmt"
	"log"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

// ExampleApp_Parse_flag_value_options shows a couple combinations of flag/value options.
// It's also possible to use '--help detailed' to see the current value of a flag and what set it.
func ExampleApp_Parse_flag_value_options() {

	action := func(ctx cli.Context) error {
		// flag marked as Required(), so no need to check for existance
		scalarVal := ctx.Flags["--scalar-flag"].(string)
		// flag might not exist in config, so check for existance
		// TODO: does this panic on nil?
		sliceVal, sliceValExists := ctx.Flags["--slice-flag"].([]int)

		fmt.Printf("--scalar-flag: %#v\n", scalarVal)
		if sliceValExists {
			fmt.Printf("--slice-flag: %#v\n", sliceVal)
		} else {
			fmt.Printf("--slice-flag value not filled!\n")
		}
		return nil
	}

	app := warg.New(
		"flag-overrides",
		"v1.0.0",
		section.New(
			"demo flag overrides",
			section.NewCommand(
				string("show"),
				"Show final flag values",
				action,
				command.NewFlag(
					"--scalar-flag",
					"Demo scalar flag",
					scalar.String(
						scalar.Choices("a", "b"),
						scalar.Default("a"),
					),
					flag.ConfigPath("args.scalar-flag"),
					flag.Required(),
				),
				command.NewFlag(
					"--slice-flag",
					"Demo slice flag",
					slice.Int(
						slice.Choices(1, 2, 3),
					),
					flag.Alias("-slice"),
					flag.ConfigPath("args.slice-flag"),
					flag.EnvVars("SLICE", "SLICE_ARG"),
				),
			),
		),
		warg.ConfigFlag(
			yamlreader.New,
			cli.FlagMap{
				"--config": flag.New(
					"Path to YAML config file",
					scalar.Path(
						scalar.Default(path.New("~/.config/flag-overrides.yaml")),
					),
					flag.Alias("-c"),
				),
			},
		),
	)

	err := os.WriteFile(
		"testdata/ExampleFlagValueOptions/config.yaml",
		[]byte(`args:
  slice-flag:
    - 1
    - 2
    - 3
`),
		0644,
	)
	if err != nil {
		log.Fatalf("write error: %e", err)
	}
	app.MustRun(
		cli.OverrideArgs([]string{"calc", "show", "-c", "testdata/ExampleFlagValueOptions/config.yaml", "--scalar-flag", "b"}),
	)
	// Output:
	// --scalar-flag: "b"
	// --slice-flag: []int{1, 2, 3}
}
