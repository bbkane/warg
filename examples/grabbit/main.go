package main

import (
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

func app() *warg.App {
	appFooter := `Examples (assuming BASH-like shell):

  # Grab from passed flags
  grabbit grab \
      --subreddit-destination . \
      --subreddit-limit 5 \
      --subreddit-name wallpapers \
      --subreddit-timeframe day

  # Create/Edit config file
  grabbit config edit --editor /path/to/editor

  # Grab from config file
  grabbit grab

Homepage: https://github.com/bbkane/grabbit
`
	grabCmd := command.New(
		"Grab images. Optionally use `config edit` first to create a config",
		grab,
		command.Flag(
			"--subreddit-name",
			"Subreddit to grab",
			slice.String(
				slice.Default([]string{"earthporn", "wallpapers"}),
			),
			flag.Alias("-sn"),
			flag.ConfigPath("subreddits[].name"),
			flag.Required(),
		),
		command.Flag(
			"--subreddit-destination",
			"Where to store the subreddit",
			slice.Path(
				slice.Default([]string{".", "."}),
			),
			flag.Alias("-sd"),
			flag.ConfigPath("subreddits[].destination"),
			flag.Required(),
		),
		command.Flag(
			"--subreddit-timeframe",
			"Take the top subreddits from this timeframe",
			slice.String(
				slice.Choices("day", "week", "month", "year", "all"),
				slice.Default([]string{"week", "week"}),
			),
			flag.Alias("-st"),
			flag.ConfigPath("subreddits[].timeframe"),
			flag.Required(),
		),
		command.Flag(
			"--subreddit-limit",
			"Max number of links to try to download",
			slice.Int(
				slice.Default([]int{2, 3}),
			),
			flag.Alias("-sl"),
			flag.ConfigPath("subreddits[].limit"),
			flag.Required(),
		),
		command.Flag(
			"--timeout",
			"Timeout for a single download",
			scalar.Duration(
				scalar.Default(time.Second*30),
			),
			flag.Alias("-t"),
			flag.Required(),
		),
	)

	app := warg.New(
		"grabbit",
		section.New(
			"Get top images from subreddits",
			section.ExistingCommand("version", warg.VersionCommand()),

			section.ExistingCommand(
				"grab",
				grabCmd,
			),
			section.ExistingFlag("--color", warg.ColorFlag()),
			section.Footer(appFooter),
			section.Flag(
				"--log-filename",
				"Log filename",
				scalar.Path(
					scalar.Default("~/.config/grabbit.jsonl"),
				),
				flag.ConfigPath("lumberjacklogger.filename"),
				flag.Required(),
			),
			section.Flag(
				"--log-maxage",
				"Max age before log rotation in days", // TODO: change to duration flag
				scalar.Int(
					scalar.Default(30),
				),
				flag.ConfigPath("lumberjacklogger.maxage"),
				flag.Required(),
			),
			section.Flag(
				"--log-maxbackups",
				"Num backups for the log",
				scalar.Int(
					scalar.Default(0),
				),
				flag.ConfigPath("lumberjacklogger.maxbackups"),
				flag.Required(),
			),
			section.Flag(
				"--log-maxsize", // TODO: make bytesize package similar to time?
				"Max size of log in megabytes",
				scalar.Int(
					scalar.Default(5),
				),
				flag.ConfigPath("lumberjacklogger.maxsize"),
				flag.Required(),
			),
			section.Section(
				"config",
				"Config commands",
				section.Command(
					"edit",
					"Edit or create configuration file.",
					editConfig,
					command.Flag(
						"--editor",
						"Path to editor",
						scalar.String(
							scalar.Default("vi"),
						),
						flag.Alias("-e"),
						flag.EnvVars("EDITOR"),
						flag.Required(),
					),
				),
			),
		),
		warg.ConfigFlag(
			"--config",
			[]scalar.ScalarOpt[string]{
				scalar.Default("~/.config/grabbit.yaml"),
			},
			yamlreader.New,
			"Config filepath",
			flag.Alias("-c"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
