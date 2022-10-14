package main

import (
	"os"
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
)

func app() *warg.App {

	downloadCmd := command.New(
		"Download star info",
		githubStarsDownload,
		command.Flag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
		),
		command.Flag(
			"--max-languages",
			"Max number of languages to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		command.Flag(
			"--max-repo-topics",
			"Max number of topics to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		command.Flag(
			"--after-cursor",
			"PageInfo EndCursor to start from",
			scalar.String(),
		),
		command.Flag(
			"--max-pages",
			"Max number of pages to fetch",
			scalar.Int(
				scalar.Default(1),
			),
			flag.Required(),
		),
		command.Flag(
			"--output",
			"Output filepath. Must not exist",
			scalar.Path(
				scalar.Default("starghaze_download.jsonl"),
			),
		),
		command.Flag(
			"--page-size",
			"Number of starred repos in page",
			scalar.Int(
				scalar.Default(100),
			),
			flag.Required(),
		),
		command.Flag(
			"--timeout",
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			scalar.Duration(
				scalar.Default(time.Minute*10),
			),
			flag.Required(),
		),
		command.Flag(
			"--token",
			"Github PAT",
			scalar.String(),
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
			scalar.String(
				scalar.Choices("csv", "jsonl", "sqlite", "zinc"),
				scalar.Default("csv"),
			),
			flag.Required(),
		),
		command.Flag(
			"--date-format",
			"Datetime output format. See https://github.com/lestrrat-go/strftime for details. If not passed, the GitHub default is RFC 3339. Consider using '%b %d, %Y' for csv format",
			scalar.String(),
		),
		command.Flag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
			flag.Required(),
		),
		command.Flag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name. Only used for --format sqlite",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
		),
		command.Flag(
			"--zinc-index-name",
			"Only used for --format zinc.",
			scalar.String(
				scalar.Default("starghaze"),
			),
		),
		command.Flag(
			"--input",
			"Input file",
			scalar.Path(
				scalar.Default("starghaze_download.jsonl"),
			),
			flag.Required(),
		),
		command.Flag(
			"--max-line-size",
			"Max line size in the file in MB",
			scalar.Int(
				scalar.Default(10),
			),
			flag.Required(),
		),
		command.Flag(
			"--output",
			"output file. Prints to stdout if not passed",
			scalar.Path(),
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
				scalar.Path(),
				flag.Required(),
			),
			command.Flag(
				"--timeout",
				"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
				scalar.Duration(
					scalar.Default(time.Minute*10),
				),
				flag.Required(),
			),
		),
		section.Flag(
			"--sheet-id",
			"ID For the particulare sheet. Viewable from `gid` URL param",
			scalar.Int(),
			flag.EnvVars("STARGHAZE_SHEET_ID"),
			flag.Required(),
		),
		section.Flag(
			"--spreadsheet-id",
			"ID for the whole spreadsheet. Viewable from URL",
			scalar.String(),
			flag.EnvVars("STARGHAZE_SPREADSHEET_ID"),
			flag.Required(),
		),
	)

	searchCmd := command.New(

		"Full text search SQLite database",
		search,
		command.Flag(
			"--limit",
			"Max number of results",
			scalar.Int(
				scalar.Default(50),
			),
			flag.Required(),
		),
		command.Flag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name.",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
			flag.Required(),
		),
		command.Flag(
			"--term",
			"Search for this term",
			scalar.String(),
			flag.Alias("-t"),
			flag.Required(),
		),

		// TODO: how many results? limit by date added?
	)

	app := warg.New(
		"starghaze",
		section.New(
			"Save GitHub Starred Repos",
			section.ExistingCommand(
				"download",
				downloadCmd,
			),
			section.ExistingCommand(
				"format",
				formatCmd,
			),
			section.ExistingCommand(
				"search",
				searchCmd,
			),
			section.ExistingSection(
				"gsheets",
				gsheetsSection,
			),
			section.Footer("Homepage: https://github.com/bbkane/starghaze"),
		),
		warg.AddVersionCommand(version),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun(os.Args, os.LookupEnv)
}
