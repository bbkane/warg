package warg_test

import (
	"fmt"
	"log"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/parseopt"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
	"go.bbkane.com/warg/wargcore"
)

func exampleConfigFlagTextAdd(ctx wargcore.Context) error {
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
		wargcore.NewSection(
			"do math",
			wargcore.NewChildCmd(
				string("add"),
				"add integers",
				exampleConfigFlagTextAdd,
				wargcore.NewChildFlag(
					string("--addend"),
					"Integer to add. Flag is repeatible",
					slice.Int(),
					wargcore.ConfigPath("add.addends"),
					wargcore.Required(),
				),
			),
		),
		warg.ConfigFlag(
			yamlreader.New,
			wargcore.FlagMap{
				"--config": wargcore.NewFlag(
					"Path to YAML config file",
					scalar.Path(
						scalar.Default(path.New("~/.config/calc.yaml")),
					),
					wargcore.Alias("-c"),
				),
			},
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
		parseopt.Args([]string{"calc", "add", "-c", "testdata/ExampleConfigFlag/calc.yaml"}),
	)
	// Output:
	// Sum: 6
}
