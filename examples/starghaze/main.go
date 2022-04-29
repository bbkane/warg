package main

import (
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

func app() *warg.App {

	downloadCmd := command.New(
		"Download star info",
		githubStarsDownload,
		command.Flag(
			"--include-readmes",
			"Search for README.md.",
			value.Bool,
			flag.Default("false"),
		),
		command.Flag(
			"--max-languages",
			"Max number of languages to query on a repo",
			value.Int,
			flag.Default("20"),
		),
		command.Flag(
			"--max-repo-topics",
			"Max number of topics to query on a repo",
			value.Int,
			flag.Default("20"),
		),
		command.Flag(
			"--after-cursor",
			"PageInfo EndCursor to start from",
			value.String,
		),
		command.Flag(
			"--max-pages",
			"Max number of pages to fetch",
			value.Int,
			flag.Default("1"),
			flag.Required(),
		),
		command.Flag(
			"--output",
			"Output filepath. Must not exist",
			value.Path,
			flag.Default("starghaze_download.jsonl"),
		),
		command.Flag(
			"--page-size",
			"Number of starred repos in page",
			value.Int,
			flag.Default("100"),
			flag.Required(),
		),
		command.Flag(
			"--timeout",
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			value.Duration,
			flag.Default("10m"),
			flag.Required(),
		),
		command.Flag(
			"--token",
			"Github PAT",
			value.String,
			flag.EnvVars("STARGHAZE_GITHUB_TOKEN", "GITHUB_TOKEN"),
			flag.Required(),
		),
	)

	formatCmd := command.New(
		"Format downloaded GitHub Stars",
		format,
		command.Flag(
			"--format",
			"Output format",
			value.StringEnum("csv", "jsonl", "sqlite", "zinc"),
			flag.Default("csv"),
			flag.Required(),
		),
		command.Flag(
			"--date-format",
			"Datetime output format. See https://github.com/lestrrat-go/strftime for details. If not passed, the GitHub default is RFC 3339. Consider using '%b %d, %Y' for csv format",
			value.String,
		),
		command.Flag(
			"--include-readmes",
			"Search for README.md.",
			value.Bool,
			flag.Default("false"),
			flag.Required(),
		),
		command.Flag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name. Only used for --format sqlite",
			value.String,
			flag.Default("starghaze.db"),
		),
		command.Flag(
			"--zinc-index-name",
			"Only used for --format zinc.",
			value.String,
			flag.Default("starghaze"),
		),
		command.Flag(
			"--input",
			"Input file",
			value.String,
			flag.Required(),
			flag.Default("starghaze_download.jsonl"),
		),
		command.Flag(
			"--max-line-size",
			"Max line size in the file in MB",
			value.Int,
			flag.Default("10"),
			flag.Required(),
		),
		command.Flag(
			"--output",
			"output file. Prints to stdout if not passed",
			value.Path,
		),
	)

	gsheetsSection := section.New(
		"Google Sheets commands",
		section.Command(
			"open",
			"Open spreadsheet in browser",
			gSheetsOpen,
		),
		section.Command(
			"upload",
			"Upload CSV to Google Sheets. This will overwrite whatever is in the spreadsheet",
			gSheetsUpload,
			command.Flag(
				"--csv-path",
				"CSV file to upload",
				value.Path,
				flag.Required(),
			),
			command.Flag(
				"--timeout",
				"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
				value.Duration,
				flag.Default("10m"),
				flag.Required(),
			),
		),
		section.Flag(
			"--sheet-id",
			"ID For the particulare sheet. Viewable from `gid` URL param",
			value.Int,
			flag.EnvVars("STARGHAZE_SHEET_ID"),
			flag.Required(),
		),
		section.Flag(
			"--spreadsheet-id",
			"ID for the whole spreadsheet. Viewable from URL",
			value.String,
			flag.EnvVars("STARGHAZE_SPREADSHEET_ID"),
			flag.Required(),
		),
	)

	app := warg.New(
		section.New(
			"Save GitHub Starred Repos",
			section.Command(
				"version",
				"Print version",
				printVersion,
			),
			section.ExistingCommand(
				"download",
				downloadCmd,
			),
			section.ExistingCommand(
				"format",
				formatCmd,
			),
			section.ExistingSection(
				"gsheets",
				gsheetsSection,
			),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun(os.Args, os.LookupEnv)
}
