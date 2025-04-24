package main

import (
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
	"go.bbkane.com/warg/wargcore"
)

func app() *wargcore.App {
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

	logFlags := wargcore.FlagMap{
		"--log-filename": flag.New(
			"Log filename",
			scalar.Path(
				scalar.Default(path.New("~/.config/grabbit.jsonl")),
			),
			flag.ConfigPath("lumberjacklogger.filename"),
			flag.Required(),
		),
		"--log-maxage": flag.New(
			"Max age before log rotation in days", // TODO: change to duration flag
			scalar.Int(
				scalar.Default(30),
			),
			flag.ConfigPath("lumberjacklogger.maxage"),
			flag.Required(),
		),
		"--log-maxbackups": flag.New(
			"Num backups for the log",
			scalar.Int(
				scalar.Default(0),
			),
			flag.ConfigPath("lumberjacklogger.maxbackups"),
			flag.Required(),
		),
		"--log-maxsize": flag.New(
			"Max size of log in megabytes",
			scalar.Int(
				scalar.Default(5),
			),
			flag.ConfigPath("lumberjacklogger.maxsize"),
			flag.Required(),
		),
	}

	app := warg.New(
		"grabbit",
		"v1.0.0",
		section.New(
			"Get top images from subreddits",
			section.NewCommand(
				"grab",
				"Grab images. Optionally use `config edit` first to create a config",
				grab,
				command.FlagMap(logFlags),
				command.NewFlag(
					"--subreddit-name",
					"Subreddit to grab",
					slice.String(
						slice.Default([]string{"earthporn", "wallpapers"}),
					),
					flag.Alias("-sn"),
					flag.ConfigPath("subreddits[].name"),
					flag.Required(),
				),
				command.NewFlag(
					"--subreddit-destination",
					"Where to store the subreddit",
					slice.Path(
						slice.Default([]path.Path{path.New("."), path.New(".")}),
					),
					flag.Alias("-sd"),
					flag.ConfigPath("subreddits[].destination"),
					flag.Required(),
				),
				command.NewFlag(
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
				command.NewFlag(
					"--subreddit-limit",
					"Max number of links to try to download",
					slice.Int(
						slice.Default([]int{2, 3}),
					),
					flag.Alias("-sl"),
					flag.ConfigPath("subreddits[].limit"),
					flag.Required(),
				),
				command.NewFlag(
					"--timeout",
					"Timeout for a single download",
					scalar.Duration(
						scalar.Default(time.Second*30),
					),
					flag.Alias("-t"),
					flag.Required(),
				),
			),

			section.Footer(appFooter),

			section.NewSection(
				"config",
				"Config commands",
				section.NewCommand(
					"edit",
					"Edit or create configuration file.",
					editConfig,
					command.FlagMap(logFlags),
					command.NewFlag(
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
			yamlreader.New,
			wargcore.FlagMap{
				"--config": flag.New(
					"Path to YAML config file",
					scalar.Path(
						scalar.Default(path.New("~/.config/grabbit.yaml")),
					),
					flag.Alias("-c"),
				),
			},
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
