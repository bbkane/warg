package main

import (
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/path"
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
		"--log-filename": wargcore.NewFlag(
			"Log filename",
			scalar.Path(
				scalar.Default(path.New("~/.config/grabbit.jsonl")),
			),
			wargcore.ConfigPath("lumberjacklogger.filename"),
			wargcore.Required(),
		),
		"--log-maxage": wargcore.NewFlag(
			"Max age before log rotation in days", // TODO: change to duration flag
			scalar.Int(
				scalar.Default(30),
			),
			wargcore.ConfigPath("lumberjacklogger.maxage"),
			wargcore.Required(),
		),
		"--log-maxbackups": wargcore.NewFlag(
			"Num backups for the log",
			scalar.Int(
				scalar.Default(0),
			),
			wargcore.ConfigPath("lumberjacklogger.maxbackups"),
			wargcore.Required(),
		),
		"--log-maxsize": wargcore.NewFlag(
			"Max size of log in megabytes",
			scalar.Int(
				scalar.Default(5),
			),
			wargcore.ConfigPath("lumberjacklogger.maxsize"),
			wargcore.Required(),
		),
	}

	app := warg.New(
		"grabbit",
		"v1.0.0",
		wargcore.NewSection(
			"Get top images from subreddits",
			wargcore.NewChildCmd(
				"grab",
				"Grab images. Optionally use `config edit` first to create a config",
				grab,
				wargcore.ChildFlagMap(logFlags),
				wargcore.NewChildFlag(
					"--subreddit-name",
					"Subreddit to grab",
					slice.String(
						slice.Default([]string{"earthporn", "wallpapers"}),
					),
					wargcore.Alias("-sn"),
					wargcore.ConfigPath("subreddits[].name"),
					wargcore.Required(),
				),
				wargcore.NewChildFlag(
					"--subreddit-destination",
					"Where to store the subreddit",
					slice.Path(
						slice.Default([]path.Path{path.New("."), path.New(".")}),
					),
					wargcore.Alias("-sd"),
					wargcore.ConfigPath("subreddits[].destination"),
					wargcore.Required(),
				),
				wargcore.NewChildFlag(
					"--subreddit-timeframe",
					"Take the top subreddits from this timeframe",
					slice.String(
						slice.Choices("day", "week", "month", "year", "all"),
						slice.Default([]string{"week", "week"}),
					),
					wargcore.Alias("-st"),
					wargcore.ConfigPath("subreddits[].timeframe"),
					wargcore.Required(),
				),
				wargcore.NewChildFlag(
					"--subreddit-limit",
					"Max number of links to try to download",
					slice.Int(
						slice.Default([]int{2, 3}),
					),
					wargcore.Alias("-sl"),
					wargcore.ConfigPath("subreddits[].limit"),
					wargcore.Required(),
				),
				wargcore.NewChildFlag(
					"--timeout",
					"Timeout for a single download",
					scalar.Duration(
						scalar.Default(time.Second*30),
					),
					wargcore.Alias("-t"),
					wargcore.Required(),
				),
			),

			wargcore.SectionFooter(appFooter),

			wargcore.NewChildSection(
				"config",
				"Config commands",
				wargcore.NewChildCmd(
					"edit",
					"Edit or create configuration file.",
					editConfig,
					wargcore.ChildFlagMap(logFlags),
					wargcore.NewChildFlag(
						"--editor",
						"Path to editor",
						scalar.String(
							scalar.Default("vi"),
						),
						wargcore.Alias("-e"),
						wargcore.EnvVars("EDITOR"),
						wargcore.Required(),
					),
				),
			),
		),
		warg.ConfigFlag(
			yamlreader.New,
			wargcore.FlagMap{
				"--config": wargcore.NewFlag(
					"Path to YAML config file",
					scalar.Path(
						scalar.Default(path.New("~/.config/grabbit.yaml")),
					),
					wargcore.Alias("-c"),
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
