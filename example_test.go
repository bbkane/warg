package warg_test

import (
	"fmt"
	"os"

	w "github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	"github.com/bbkane/warg/configreader/yamlreader"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

func Example_grabbit_help() {
	app := w.New(
		"grabbit",
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
					f.ConfigPath("subreddits[].name"),
				),
				c.WithFlag(
					"--subreddit-destination",
					"Where to store the subreddit",
					v.StringSliceEmpty,
					f.Default("~/Pictures/grabbit"),
					f.ConfigPath("subreddits[].destination"),
				),
				c.WithFlag(
					"--subreddit-timeframe",
					"Take the top subreddits from this timeframe",
					v.StringSliceEmpty,
					f.Default("week"),
					f.ConfigPath("subreddits[].timeframe"),
				),
				c.WithFlag(
					"--subreddit-limit",
					"max number of links to try to download",
					v.IntSliceEmpty,
					f.Default("5"),
					f.ConfigPath("subreddit[].limit"),
				),
			),
			s.WithFlag(
				"--log-filename",
				"log filename",
				v.StringEmpty,
				f.Default("~/.config/grabbit.jsonl"),
				f.ConfigPath("lumberjacklogger.filename"),
			),
			s.WithFlag(
				"--log-maxage",
				"max age before log rotation in days",
				v.IntEmpty,
				f.Default("30"),
				f.ConfigPath("lumberjacklogger.maxage"),
			),
			s.WithFlag(
				"--log-maxbackups",
				"num backups for the log",
				v.IntEmpty,
				f.Default("0"),
				f.ConfigPath("lumberjacklogger.maxbackups"),
			),
			s.WithFlag(
				"--log-maxsize",
				"max size of log in megabytes",
				v.IntEmpty,
				f.Default("5"),
				f.ConfigPath("lumberjacklogger.maxsize"),
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
			yamlreader.NewYAMLConfigReader,
			"config filepath",
			f.Default("/path/to/config.yaml"),
		),
		w.OverrideHelp(
			os.Stdout,
			[]string{"-h", "--help"},
			w.DefaultSectionHelp,
			w.DefaultCommandHelp,
		),
	)

	args := []string{"example", "config", "edit", "-h"}
	err := app.Run(args)
	if err != nil {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	// Output:
	// Edit or create configuration file. Uses $EDITOR as a fallback
	//
	// Flags:
	//
	//   --config-path : config filepath
	//     value : /path/to/config.yaml
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
