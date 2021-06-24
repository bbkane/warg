package warg_test

import (
	"fmt"

	a "github.com/bbkane/warg/app"
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

	app := a.NewApp(
		a.AppRootCategory(
			s.WithCategory(
				"cat",
				s.WithCommand(
					"com",
					c.WithAction(comAction),
					c.WithCommandFlag(
						"--flag",
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
	app := a.NewApp(
		a.EnableVersionFlag(
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

func Example_help() {
	app := a.NewApp(
		a.EnableHelpFlag(
			[]string{"-h", "--help"},
			"example",
		),
		a.AppRootCategory(
			s.WithCategoryHelpShort("example help!"),
			s.WithCategory(
				"cat",
				s.WithCategoryHelpShort("cat help!"),
			),
			s.WithCommand(
				"com",
				c.WithCommandHelpShort("com help!!"),
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
}

func Example_grabbit() {
	_ = a.NewApp2(
		[]a.AppOpt{
			a.EnableHelpFlag([]string{"--help", "-h"}, "grabbit"),
			a.EnableVersionFlag([]string{"--version"}, "v.0.0.1"),
		},
		s.WithCategoryHelpShort("grab pics from reddit!"),
		s.WithCategory(
			"config",
			s.WithCategoryHelpShort("work with the config"),
			s.WithCommand(
				"edit",
				c.WithCommandHelpShort("edit the config"),
				c.WithCommandFlag(
					"--editor",
					v.NewEmptyStringValue(),
					f.WithDefault(v.NewStringValue("vi")),
					f.WithFlagHelpShort("path to editor"),
				),
			),
		),
		s.WithCommand(
			"grab",
			c.WithAction(
				func(vm v.ValueMap) error {
					return nil
				},
			),
		),
	)
}
