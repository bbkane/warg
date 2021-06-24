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

	app := a.New(
		a.RootSection(
			s.WithSection(
				"cat",
				s.WithCommand(
					"com",
					c.WithAction(comAction),
					c.WithFlag(
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
	app := a.New(
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
	app := a.New(
		a.EnableHelpFlag(
			[]string{"-h", "--help"},
			"example",
		),
		a.RootSection(
			s.HelpShort("example help!"),
			s.WithSection(
				"cat",
				s.HelpShort("cat help!"),
			),
			s.WithCommand(
				"com",
				c.HelpShort("com help!!"),
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
	_ = a.New2(
		[]a.AppOpt{
			a.EnableHelpFlag([]string{"--help", "-h"}, "grabbit"),
			a.EnableVersionFlag([]string{"--version"}, "v.0.0.1"),
		},
		s.HelpShort("grab pics from reddit!"),
		s.WithSection(
			"config",
			s.HelpShort("work with the config"),
			s.WithCommand(
				"edit",
				c.HelpShort("edit the config"),
				c.WithFlag(
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
