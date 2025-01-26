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

func exampleConfigFlagTextAdd(ctx command.Context) error {
	addends := ctx.Flags["--addend"].([]int)
	sum := 0
	for _, a := range addends {
		sum += a
	}
	fmt.Printf("Sum: %d\n", sum)
	return nil
}

func ExampleConfigFlag() {
	app := warg.New(
		"newAppName",
		"v1.0.0",
		section.New(
			"do math",
			section.Command(
				command.Name("add"),
				"add integers",
				exampleConfigFlagTextAdd,
				command.Flag(
					flag.Name("--addend"),
					"Integer to add. Flag is repeatible",
					slice.Int(),
					flag.ConfigPath("add.addends"),
					flag.Required(),
				),
			),
		),
		warg.ConfigFlag(
			"--config",
			[]scalar.ScalarOpt[path.Path]{
				scalar.Default(path.New("~/.config/calc.yaml")),
			},
			yamlreader.New,
			"path to YAML config file",
			flag.Alias("-c"),
		),
	)

	err := os.WriteFile(
		"testdata/ExampleConfigFlag/calc.yaml",
		[]byte(`add:
  addends:
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
		warg.OverrideArgs([]string{"calc", "add", "-c", "testdata/ExampleConfigFlag/calc.yaml"}),
	)
	// Output:
	// Sum: 6
}
