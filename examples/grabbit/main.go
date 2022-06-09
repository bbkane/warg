package main

import (
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
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
			value.StringSlice,
			flag.Alias("-sn"),
			flag.Default("earthporn", "wallpapers"),
			flag.ConfigPath("subreddits[].name"),
			flag.Required(),
		),
		command.Flag(
			"--subreddit-destination",
			"Where to store the subreddit",
			value.PathSlice,
			flag.Alias("-sd"),
			flag.Default(".", "."),
			flag.ConfigPath("subreddits[].destination"),
			flag.Required(),
		),
		command.Flag(
			"--subreddit-timeframe",
			"Take the top subreddits from this timeframe",
			value.StringEnumSlice("day", "week", "month", "year", "all"),
			flag.Alias("-st"),
			flag.Default("week", "week"),
			flag.ConfigPath("subreddits[].timeframe"),
			flag.Required(),
		),
		command.Flag(
			"--subreddit-limit",
			"Max number of links to try to download",
			value.IntSlice,
			flag.Alias("-sl"),
			flag.Default("2", "3"),
			flag.ConfigPath("subreddits[].limit"),
			flag.Required(),
		),
		command.Flag(
			"--timeout",
			"Timeout for a single download",
			value.Duration,
			flag.Alias("-t"),
			flag.Default("30s"),
			flag.Required(),
		),
	)

	app := warg.New(
		"grabbit",
		section.New(
			"Get top images from subreddits",
			section.ExistingCommand(
				"grab",
				grabCmd,
			),
			section.Footer(appFooter),
			section.Command(
				"version",
				"Print version",
				printVersion,
			),
			section.Flag(
				"--color",
				"Use colorized output",
				value.StringEnum("true", "false", "auto"),
				flag.Default("auto"),
			),
			section.Flag(
				"--log-filename",
				"Log filename",
				value.Path,
				flag.Default("~/.config/grabbit.jsonl"),
				flag.ConfigPath("lumberjacklogger.filename"),
				flag.Required(),
			),
			section.Flag(
				"--log-maxage",
				"Max age before log rotation in days",
				value.Int,
				flag.Default("30"),
				flag.ConfigPath("lumberjacklogger.maxage"),
				flag.Required(),
			),
			section.Flag(
				"--log-maxbackups",
				"Num backups for the log",
				value.Int,
				flag.Default("0"),
				flag.ConfigPath("lumberjacklogger.maxbackups"),
				flag.Required(),
			),
			section.Flag(
				"--log-maxsize",
				"Max size of log in megabytes",
				value.Int,
				flag.Default("5"),
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
						value.String,
						flag.Alias("-e"),
						flag.Default("vi"),
						flag.EnvVars("EDITOR"),
						flag.Required(),
					),
				),
			),
		),
		warg.ConfigFlag(
			"--config",
			yamlreader.New,
			"Config filepath",
			flag.Alias("-c"),
			flag.Default("~/.config/grabbit.yaml"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun(os.Args, os.LookupEnv)
}
