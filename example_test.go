package warg_test

import (
	"fmt"

	w "github.com/bbkane/warg"
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

	app := w.New(
		"test",
		"v0.0.0",
		s.NewSection(
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
						v.IntNew(0),
						f.Default(v.IntNew(10)),
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
	app := w.New(
		"test",
		"v0.0.0",
		s.NewSection(""),
		w.OverrideVersion(
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

func Example_section_help() {
	app := w.New(
		"grabbit",
		"v0.0.0",
		s.NewSection(
			"grab those images!",
			s.WithSection(
				"config",
				"change grabbit's config",
				s.WithCommand(
					"edit",
					"edit the config",
					c.DoNothing,
					c.WithFlag(
						"--editor",
						"path to editor",
						v.StringEmpty(),
						f.Default(v.StringNew("vi")),
					),
				),
			),
			s.WithCommand(
				"grab",
				"do the grabbity grabbity",
				c.DoNothing,
			),
		),
		w.OverrideHelp(
			[]string{"-h", "--help"},
			w.DefaultSectionHelp,
			w.DefaultCommandHelp,
		),
	)
	args := []string{"grabbit", "--help"}
	pr, err := app.Parse(args)
	if err != nil {
		panic(err)
	}
	err = pr.Action(pr.PassedFlags)
	if err != nil {
		panic(err)
	}
	// Output:
	// grab those images!
	//
	// Sections:
	//   config : change grabbit's config
	//
	// Commands:
	//   grab : do the grabbity grabbity
}

func Example_command_help() {
	app := w.New(
		"grabbit",
		"v0.0.0",
		s.NewSection(
			"grab those images!",
			s.WithFlag(
				"--config-path",
				"path to config",
				v.StringEmpty(),
				f.Default(v.StringNew("~/.config/grabbit.yaml")),
			),
			s.WithSection(
				"config",
				"change grabbit's config",
				s.WithCommand(
					"edit",
					"edit the config",
					c.DoNothing,
					c.WithFlag(
						"--editor",
						"path to editor",
						v.StringEmpty(),
						f.Default(v.StringNew("vi")),
					),
				),
			),
			s.WithCommand(
				"grab",
				"do the grabbity grabbity",
				c.DoNothing,
			),
		),
		w.OverrideHelp(
			[]string{"-h", "--help"},
			w.DefaultSectionHelp,
			w.DefaultCommandHelp,
		),
	)
	args := []string{"example", "config", "edit", "-h"}
	pr, err := app.Parse(args)
	if err != nil {
		panic(err)
	}
	err = pr.Action(pr.PassedFlags)
	if err != nil {
		panic(err)
	}
	// Output:
	// edit the config

	// Flags:
	//
	//  --config-path : path to config
	// 	  value : ~/.config/grabbit.yaml
	// 	  setby : appdefault

	//  --editor : path to editor
	// 	  value : vi
	// 	  setby : appdefault
}
