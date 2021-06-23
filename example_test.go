package warg_test

import (
	"fmt"

	"github.com/bbkane/warg"
	w "github.com/bbkane/warg"
)

func Example_parse() {

	comAction := func(vm w.ValueMap) error {
		action := vm["--flag"].Get().(int)
		fmt.Printf("Action Output: %v\n", action)
		return nil
	}

	app := w.NewApp(
		w.AppRootCategory(
			w.WithCategory(
				"cat",
				w.WithCommand(
					"com",
					w.WithAction(comAction),
					w.WithCommandFlag(
						"--flag",
						w.NewIntValue(0),
						w.WithDefault(w.NewIntValue(10)),
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
	app := w.NewApp(
		w.EnableVersionFlag(
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
	app := w.NewApp(
		w.EnableHelpFlag(
			[]string{"-h", "--help"},
			"example",
		),
		w.AppRootCategory(
			w.WithCategoryHelpShort("example help!"),
			w.WithCategory(
				"cat",
				w.WithCategoryHelpShort("cat help!"),
			),
			w.WithCommand(
				"com",
				w.WithCommandHelpShort("com help!!"),
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
	_ = w.NewApp2(
		[]w.AppOpt{
			w.EnableHelpFlag([]string{"--help", "-h"}, "grabbit"),
			w.EnableVersionFlag([]string{"--version"}, "v.0.0.1"),
		},
		w.WithCategoryHelpShort("grab pics from reddit!"),
		w.WithCategory(
			"config",
			w.WithCategoryHelpShort("work with the config"),
			w.WithCommand(
				"edit",
				w.WithCommandHelpShort("edit the config"),
				w.WithCommandFlag(
					"--editor",
					w.NewEmptyStringValue(),
					w.WithDefault(w.NewStringValue("vi")),
					w.WithFlagHelpShort("path to editor"),
				),
			),
		),
		w.WithCommand(
			"grab",
			w.WithAction(
				func(vm warg.ValueMap) error {
					return nil
				},
			),
		),
	)
}
