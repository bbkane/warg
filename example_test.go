package warg_test

import (
	"fmt"

	w "github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	"github.com/bbkane/warg/configreader/jsonreader"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

func Example_parse() {
	comAction := func(fvs f.FlagValues) error {
		action := fvs["--flag"].(int)
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
						v.IntEmpty,
						f.Default("10"),
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
	err := app.Run(args)
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
						v.StringEmpty,
						f.Default("vi"),
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
	err := app.Run(args)
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

func Example_grabbit_help() {
	app := w.New(
		"grabbit",
		"v0.0.0",
		s.NewSection(
			"Get top images from subreddits",
			s.WithCommand(
				"grab",
				"Grab images. Use `config edit` first to create a config",
				c.DoNothing,
				c.WithFlag(
					"--subreddit-name",
					"subreddit to grab",
					v.StringSliceEmpty,
					f.Default("wallpapers"),
					f.ConfigPath("subreddits[].name", v.StringSliceFromInterface),
				),
				c.WithFlag(
					"--subreddit-destination",
					"Where to store the subreddit",
					v.StringSliceEmpty,
					f.Default("~/Pictures/grabbit"),
					f.ConfigPath("subreddits[].destination", v.StringSliceFromInterface),
				),
				c.WithFlag(
					"--subreddit-timeframe",
					"Take the top subreddits from this timeframe",
					v.StringSliceEmpty,
					f.Default("week"),
					f.ConfigPath("subreddits[].timeframe", v.StringSliceFromInterface),
				),
				c.WithFlag(
					"--subreddit-limit",
					"max number of links to try to download",
					v.IntSliceEmpty,
					f.Default("5"),
					f.ConfigPath("subreddit[].limit", v.IntSliceFromInterface),
				),
			),
			s.WithFlag(
				"--log-filename",
				"log filename",
				v.StringEmpty,
				f.Default("~/.config/grabbit.jsonl"),
				f.ConfigPath("lumberjacklogger.filename", v.StringFromInterface),
			),
			s.WithFlag(
				"--log-maxage",
				"max age before log rotation in days",
				v.IntEmpty,
				f.Default("30"),
				f.ConfigPath("lumberjacklogger.maxage", v.IntFromInterface),
			),
			s.WithFlag(
				"--log-maxbackups",
				"num backups for the log",
				v.IntEmpty,
				f.Default("0"),
				f.ConfigPath("lumberjacklogger.maxbackups", v.IntFromInterface),
			),
			s.WithFlag(
				"--log-maxsize",
				"max size of log in megabytes",
				v.IntEmpty,
				f.Default("5"),
				f.ConfigPath("lumberjacklogger.maxsize", v.IntFromInterface),
			),
			s.WithSection(
				"config",
				"config commands",
				s.WithCommand(
					"edit",
					"Edit or create configuration file. Uses $EDITOR as a fallback",
					c.DoNothing,
					c.WithFlag(
						"--editor",
						"path to editor",
						v.StringEmpty,
						f.Default("vi"),
					),
				),
			),
		),
		w.ConfigFlag(
			"--config-path",
			jsonreader.NewJSONConfigReader,
			"config filepath",
			f.Default("~/.config/grabbit.yaml"),
		),
	)

	args := []string{"example", "config", "edit", "-h"}
	err := app.Run(args)
	if err != nil {
		panic(err)
	}

	// Output:
	// Edit or create configuration file. Uses $EDITOR as a fallback
	//
	// Flags:
	//
	//   --config-path : config filepath
	//     value : ~/.config/grabbit.yaml
	//     setby : appdefault
	//
	//   --editor : path to editor
	//     value : vi
	//     setby : appdefault
	//
	//   --log-filename : log filename
	//     configpath : lumberjacklogger.filename
	//     value : ~/.config/grabbit.jsonl
	//     setby : appdefault
	//
	//   --log-maxage : max age before log rotation in days
	//     configpath : lumberjacklogger.maxage
	//     value : 30
	//     setby : appdefault
	//
	//   --log-maxbackups : num backups for the log
	//     configpath : lumberjacklogger.maxbackups
	//     value : 0
	//     setby : appdefault
	//
	//   --log-maxsize : max size of log in megabytes
	//     configpath : lumberjacklogger.maxsize
	//     value : 5
	//     setby : appdefault
}
