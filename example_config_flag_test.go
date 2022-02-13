package warg_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/config/yamlreader"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
)

func exampleConfigFlagTextAdd(pf flag.PassedFlags) error {
	addends := pf["--addend"].([]int)
	sum := 0
	for _, a := range addends {
		sum += a
	}
	fmt.Printf("Sum: %d\n", sum)
	return nil
}

func ExampleConfigFlag() {
	app := warg.New(
		"calc",
		section.New(
			"do math",
			section.Command(
				command.Name("add"),
				"add integers",
				exampleConfigFlagTextAdd,
				command.Flag(
					flag.Name("--addend"),
					"Integer to add. Floats will be truncated. Flag is repeatible",
					value.IntSlice,
					flag.ConfigPath("add.addends"),
					flag.Required(),
				),
			),
		),
		warg.ConfigFlag(
			"--config",
			yamlreader.New,
			"path to YAML config file",
			flag.Alias("-c"),
			flag.Default("~/.config/calc.yaml"),
		),
	)

	err := ioutil.WriteFile(
		"/tmp/calc.yaml",
		[]byte(
			`
add:
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
	app.MustRun([]string{"calc", "-c", "calc.yaml", "add"}, os.LookupEnv)
	// Output:
	// Sum: 6
}