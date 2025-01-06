package warg_test

import (
	"fmt"
	"log"
	"os"

	"go.bbkane.com/warg"
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

	action := func(ctx command.Context) error {
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
		section.New(
			"demo flag overrides",
			section.Command(
				command.Name("show"),
				"Show final flag values",
				action,
				command.Flag(
					"--scalar-flag",
					"Demo scalar flag",
					scalar.String(
						scalar.Choices("a", "b"),
						scalar.Default("a"),
					),
					flag.ConfigPath("args.scalar-flag"),
					flag.Required(),
				),
				command.Flag(
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
			"--config",
			[]scalar.ScalarOpt[path.Path]{
				scalar.Default(path.New("~/.config/flag-overrides.yaml")),
			},
			yamlreader.New,
			"path to YAML config file",
			flag.Alias("-c"),
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
		warg.OverrideArgs([]string{"calc", "show", "-c", "testdata/ExampleFlagValueOptions/config.yaml", "--scalar-flag", "b"}),
	)
	// Output:
	// --scalar-flag: "b"
	// --slice-flag: []int{1, 2, 3}
}
