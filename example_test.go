package warg_test

import (
	"fmt"

	a "github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

func Example_parse() {

	comAction := func(vm v.ValueMap) error {
		action := vm["--flag"].Get().(int)
		fmt.Printf("Action Output: %v\n", action)
		return nil
	}

	app := a.New(
		"test",
		"v0.0.0",
		a.WithRootSection(
			"help for test",
			s.WithSection(
				"cat",
				"help for cat",
				s.WithCommand(
					"com",
					"help for com",
					comAction,
					c.WithFlag(
						"--flag",
						"flag help",
						v.NewIntValue(0),
						f.WithDefault(v.NewIntValue(10)),
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
	fmt.Printf("PassedCmd: %v\n", pr.PasssedPath)
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
	app := a.New(
		"test",
		"v0.0.0",
		a.OverrideVersion(
			[]string{"--version"},
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

func Example_help() {
	app := a.New(
		"example",
		"v0.0.0",
		a.OverrideHelp(
			[]string{"-h", "--help"},
		),
		a.WithRootSection(
			"example help!",
			s.WithSection(
				"cat",
				"cat help!",
			),
			s.WithCommand(
				"com",
				"com help!!",
				c.DoNothing,
			),
		),
	)
	args := []string{"example", "--help"}
	pr, err := app.Parse(args)
	if err != nil {
		panic(err)
	}
	err = pr.Action(pr.PassedFlags)
	if err != nil {
		panic(err)
	}
	// Output:
	// Current Category:
	//   example : example help!
	// Subcategories:
	//   cat: cat help!
	// Commands:
	//   com: com help!!
}

func Example_grabbit() {
	_ = a.New(
		"grabbit",
		"v0.0.0",
		a.OverrideHelp([]string{"--help", "-h"}),
		a.OverrideVersion([]string{"--version"}),

		a.WithRootSection(
			"grab pics from reddit!",
			s.WithSection(
				"config",
				"work with the config",
				s.WithCommand(
					"edit",
					"edit the config",
					c.DoNothing,
					c.WithFlag(
						"--editor",
						"path to editor",
						v.NewEmptyStringValue(),
						f.WithDefault(v.NewStringValue("vi")),
					),
				),
			),
			s.WithCommand(
				"grab",
				"download the images!",
				c.DoNothing,
			),
		),
	)
}
