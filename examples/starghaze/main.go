package main

import (
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
)

func app() *warg.App {

	downloadCmd := command.New(
		"Download star info",
		githubStarsDownload,
		command.NewFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
		),
		command.NewFlag(
			"--max-languages",
			"Max number of languages to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		command.NewFlag(
			"--max-repo-topics",
			"Max number of topics to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		command.NewFlag(
			"--after-cursor",
			"PageInfo EndCursor to start from",
			scalar.String(),
		),
		command.NewFlag(
			"--max-pages",
			"Max number of pages to fetch",
			scalar.Int(
				scalar.Default(1),
			),
			flag.Required(),
		),
		command.NewFlag(
			"--output",
			"Output filepath. Must not exist",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
		),
		command.NewFlag(
			"--page-size",
			"Number of starred repos in page",
			scalar.Int(
				scalar.Default(100),
			),
			flag.Required(),
		),
		command.NewFlag(
			"--timeout",
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			scalar.Duration(
				scalar.Default(time.Minute*10),
			),
			flag.Required(),
		),
		command.NewFlag(
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
		command.NewFlag(
			"--format",
			"Output format",
			scalar.String(
				scalar.Choices("csv", "jsonl", "sqlite", "zinc"),
				scalar.Default("csv"),
			),
			flag.Required(),
		),
		command.NewFlag(
			"--date-format",
			"Datetime output format. See https://github.com/lestrrat-go/strftime for details. If not passed, the GitHub default is RFC 3339. Consider using '%b %d, %Y' for csv format",
			scalar.String(),
		),
		command.NewFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
			flag.Required(),
		),
		command.NewFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name. Only used for --format sqlite",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
		),
		command.NewFlag(
			"--zinc-index-name",
			"Only used for --format zinc.",
			scalar.String(
				scalar.Default("starghaze"),
			),
		),
		command.NewFlag(
			"--input",
			"Input file",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
			flag.Required(),
		),
		command.NewFlag(
			"--max-line-size",
			"Max line size in the file in MB",
			scalar.Int(
				scalar.Default(10),
			),
			flag.Required(),
		),
		command.NewFlag(
			"--output",
			"output file. Prints to stdout if not passed",
			scalar.Path(),
		),
	)

	sheetFlags := flag.FlagMap{
		"--sheet-id": flag.NewFlag(
			"ID For the particulare sheet. Viewable from `gid` URL param",
			scalar.Int(),
			flag.EnvVars("STARGHAZE_SHEET_ID"),
			flag.Required(),
		),
		"--spreadsheet-id": flag.NewFlag(
			"ID for the whole spreadsheet. Viewable from URL",
			scalar.String(),
			flag.EnvVars("STARGHAZE_SPREADSHEET_ID"),
			flag.Required(),
		),
	}

	gsheetsSection := section.New(
		"Google Sheets commands",
		section.NewCommand(
			"open",
			"Open spreadsheet in browser",
			gSheetsOpen,
			command.FlagMap(sheetFlags),
		),
		section.NewCommand(
			"upload",
			"Upload CSV to Google Sheets. This will overwrite whatever is in the spreadsheet",
			gSheetsUpload,
			command.FlagMap(sheetFlags),
			command.NewFlag(
				"--csv-path",
				"CSV file to upload",
				scalar.Path(),
				flag.Required(),
			),
			command.NewFlag(
				"--timeout",
				"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
				scalar.Duration(
					scalar.Default(time.Minute*10),
				),
				flag.Required(),
			),
		),
	)

	searchCmd := command.New(

		"Full text search SQLite database",
		search,
		command.NewFlag(
			"--limit",
			"Max number of results",
			scalar.Int(
				scalar.Default(50),
			),
			flag.Required(),
		),
		command.NewFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name.",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
			flag.Required(),
		),
		command.NewFlag(
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
		"v1.0.0",
		section.New(
			"Save GitHub Starred Repos",
			section.CommandMap(warg.VersionCommandMap()),
			section.Command(
				"download",
				downloadCmd,
			),
			section.Command(
				"format",
				formatCmd,
			),
			section.Command(
				"search",
				searchCmd,
			),
			section.Section(
				"gsheets",
				gsheetsSection,
			),
			section.Footer("Homepage: https://github.com/bbkane/starghaze"),
		),
		warg.GlobalFlagMap(warg.ColorFlagMap()),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
