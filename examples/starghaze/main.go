package main

import (
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func app() *wargcore.App {

	downloadCmd := wargcore.NewCmd(
		"Download star info",
		githubStarsDownload,
		wargcore.NewChildFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
		),
		wargcore.NewChildFlag(
			"--max-languages",
			"Max number of languages to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		wargcore.NewChildFlag(
			"--max-repo-topics",
			"Max number of topics to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		wargcore.NewChildFlag(
			"--after-cursor",
			"PageInfo EndCursor to start from",
			scalar.String(),
		),
		wargcore.NewChildFlag(
			"--max-pages",
			"Max number of pages to fetch",
			scalar.Int(
				scalar.Default(1),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--output",
			"Output filepath. Must not exist",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
		),
		wargcore.NewChildFlag(
			"--page-size",
			"Number of starred repos in page",
			scalar.Int(
				scalar.Default(100),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--timeout",
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			scalar.Duration(
				scalar.Default(time.Minute*10),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--token",
			"Github PAT",
			scalar.String(),
			wargcore.EnvVars("STARGHAZE_GITHUB_TOKEN", "GITHUB_TOKEN"),
			wargcore.Required(),
		),
	)

	formatCmd := wargcore.NewCmd(
		"Format downloaded GitHub Stars",
		format,
		wargcore.NewChildFlag(
			"--format",
			"Output format",
			scalar.String(
				scalar.Choices("csv", "jsonl", "sqlite", "zinc"),
				scalar.Default("csv"),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--date-format",
			"Datetime output format. See https://github.com/lestrrat-go/strftime for details. If not passed, the GitHub default is RFC 3339. Consider using '%b %d, %Y' for csv format",
			scalar.String(),
		),
		wargcore.NewChildFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name. Only used for --format sqlite",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
		),
		wargcore.NewChildFlag(
			"--zinc-index-name",
			"Only used for --format zinc.",
			scalar.String(
				scalar.Default("starghaze"),
			),
		),
		wargcore.NewChildFlag(
			"--input",
			"Input file",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--max-line-size",
			"Max line size in the file in MB",
			scalar.Int(
				scalar.Default(10),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--output",
			"output file. Prints to stdout if not passed",
			scalar.Path(),
		),
	)

	sheetFlags := wargcore.FlagMap{
		"--sheet-id": wargcore.NewFlag(
			"ID For the particulare sheet. Viewable from `gid` URL param",
			scalar.Int(),
			wargcore.EnvVars("STARGHAZE_SHEET_ID"),
			wargcore.Required(),
		),
		"--spreadsheet-id": wargcore.NewFlag(
			"ID for the whole spreadsheet. Viewable from URL",
			scalar.String(),
			wargcore.EnvVars("STARGHAZE_SPREADSHEET_ID"),
			wargcore.Required(),
		),
	}

	gsheetsSection := wargcore.NewSection(
		"Google Sheets commands",
		wargcore.NewChildCmd(
			"open",
			"Open spreadsheet in browser",
			gSheetsOpen,
			wargcore.ChildFlagMap(sheetFlags),
		),
		wargcore.NewChildCmd(
			"upload",
			"Upload CSV to Google Sheets. This will overwrite whatever is in the spreadsheet",
			gSheetsUpload,
			wargcore.ChildFlagMap(sheetFlags),
			wargcore.NewChildFlag(
				"--csv-path",
				"CSV file to upload",
				scalar.Path(),
				wargcore.Required(),
			),
			wargcore.NewChildFlag(
				"--timeout",
				"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
				scalar.Duration(
					scalar.Default(time.Minute*10),
				),
				wargcore.Required(),
			),
		),
	)

	searchCmd := wargcore.NewCmd(

		"Full text search SQLite database",
		search,
		wargcore.NewChildFlag(
			"--limit",
			"Max number of results",
			scalar.Int(
				scalar.Default(50),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name.",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
			wargcore.Required(),
		),
		wargcore.NewChildFlag(
			"--term",
			"Search for this term",
			scalar.String(),
			wargcore.Alias("-t"),
			wargcore.Required(),
		),

		// TODO: how many results? limit by date added?
	)

	app := warg.New(
		"starghaze",
		"v1.0.0",
		wargcore.NewSection(
			"Save GitHub Starred Repos",
			wargcore.ChildCmd(
				"download",
				downloadCmd,
			),
			wargcore.ChildCmd(
				"format",
				formatCmd,
			),
			wargcore.ChildCmd(
				"search",
				searchCmd,
			),
			wargcore.ChildSection(
				"gsheets",
				gsheetsSection,
			),
			wargcore.SectionFooter("Homepage: https://github.com/bbkane/starghaze"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
