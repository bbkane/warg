package warg_test

import (
	"fmt"
	"log"
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

func exampleConfigFlagTextAdd(ctx warg.CmdContext) error {
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
		warg.NewSection(
			"do math",
			warg.NewSubCmd(
				string("add"),
				"add integers",
				exampleConfigFlagTextAdd,
				warg.NewCmdFlag(
					string("--addend"),
					"Integer to add. Flag is repeatible",
					slice.Int(),
					warg.ConfigPath("add.addends"),
					warg.Required(),
				),
			),
		),
		warg.ConfigFlag(
			yamlreader.New,
			warg.FlagMap{
				"--config": warg.NewFlag(
					"Path to YAML config file",
					scalar.Path(
						scalar.Default(path.New("~/.config/calc.yaml")),
					),
					warg.Alias("-c"),
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
		warg.ParseWithArgs([]string{"calc", "add", "-c", "testdata/ExampleConfigFlag/calc.yaml"}),
	)
	// Output:
	// Sum: 6
}
