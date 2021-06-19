package clide_test

import (
	"fmt"

	c "github.com/bbkane/clide"
)

func Example_parse() {

	comAction := func(vm c.ValueMap) error {
		action := vm["--flag"].Get().(int)
		fmt.Printf("Action Output: %v\n", action)
		return nil
	}

	app := c.NewApp(
		c.AppRootCategory(
			c.WithCategory(
				"cat",
				c.WithCommand(
					"com",
					c.WithAction(comAction),
					c.WithCommandFlag(
						"--flag",
						c.NewIntValue(0),
						c.WithDefault(c.NewIntValue(10)),
					),
				),
			),
		),
	)

	args := []string{"example", "cat", "com", "--flag", "3"}
	pr, err := app.Parse(args)
	if err != nil {
		fmt.Printf("Parse Error: %#v\n", err)
		return
	}
	fmt.Printf("PassedCmd: %v\n", pr.PassedCmd)
	// fmt.Printf("PassedFlags: %#v\n", pr.PassedFlags)
	if pr.Action == nil {
		fmt.Println("Action is nil..")
	}
	err = pr.Action(pr.PassedFlags)
	if err != nil {
		fmt.Printf("Action Error: %v\n", err)
		return
	}
	// Output:
	// PassedCmd: [cat com]
	// Action Output: 3
}

func Example_version() {
	app := c.NewApp(
		c.EnableVersionFlag(
			[]string{"--version"},
			"v0.0.0",
		),
	)
	args := []string{"example", "--version"}
	pr, err := app.Parse(args)
	if err != nil {
		panic(err)
	}
	err = pr.Action(pr.PassedFlags)
	if err != nil {
		panic(err)
	}
	// Output:
	// v0.0.0
}

// func Example_help() {
// 	app := c.NewApp(
// 		c.EnableHelpFlag(
// 			[]string{"-h", "--help"},
// 			"example",
// 		),
// 		c.AppRootCategory(),
// 	)
// 	args := []string{"example", "--version"}
// 	pr, err := app.Parse(args)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = pr.Action(pr.PassedFlags)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Output:
// 	// v0.0.0
// }
