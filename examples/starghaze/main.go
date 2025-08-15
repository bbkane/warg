package main

import (
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value/scalar"
)

func app() *warg.App {

	downloadCmd := warg.NewCmd(
		"Download star info",
		githubStarsDownload,
		warg.NewCmdFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
		),
		warg.NewCmdFlag(
			"--max-languages",
			"Max number of languages to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		warg.NewCmdFlag(
			"--max-repo-topics",
			"Max number of topics to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		warg.NewCmdFlag(
			"--after-cursor",
			"PageInfo EndCursor to start from",
			scalar.String(),
		),
		warg.NewCmdFlag(
			"--max-pages",
			"Max number of pages to fetch",
			scalar.Int(
				scalar.Default(1),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--output",
			"Output filepath. Must not exist",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
		),
		warg.NewCmdFlag(
			"--page-size",
			"Number of starred repos in page",
			scalar.Int(
				scalar.Default(100),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--timeout",
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			scalar.Duration(
				scalar.Default(time.Minute*10),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--token",
			"Github PAT",
			scalar.String(),
			warg.EnvVars("STARGHAZE_GITHUB_TOKEN", "GITHUB_TOKEN"),
			warg.Required(),
		),
	)

	formatCmd := warg.NewCmd(
		"Format downloaded GitHub Stars",
		format,
		warg.NewCmdFlag(
			"--format",
			"Output format",
			scalar.String(
				scalar.Choices("csv", "jsonl", "sqlite", "zinc"),
				scalar.Default("csv"),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--date-format",
			"Datetime output format. See https://github.com/lestrrat-go/strftime for details. If not passed, the GitHub default is RFC 3339. Consider using '%b %d, %Y' for csv format",
			scalar.String(),
		),
		warg.NewCmdFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name. Only used for --format sqlite",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
		),
		warg.NewCmdFlag(
			"--zinc-index-name",
			"Only used for --format zinc.",
			scalar.String(
				scalar.Default("starghaze"),
			),
		),
		warg.NewCmdFlag(
			"--input",
			"Input file",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--max-line-size",
			"Max line size in the file in MB",
			scalar.Int(
				scalar.Default(10),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--output",
			"output file. Prints to stdout if not passed",
			scalar.Path(),
		),
	)

	sheetFlags := warg.FlagMap{
		"--sheet-id": warg.NewFlag(
			"ID For the particulare sheet. Viewable from `gid` URL param",
			scalar.Int(),
			warg.EnvVars("STARGHAZE_SHEET_ID"),
			warg.Required(),
		),
		"--spreadsheet-id": warg.NewFlag(
			"ID for the whole spreadsheet. Viewable from URL",
			scalar.String(),
			warg.EnvVars("STARGHAZE_SPREADSHEET_ID"),
			warg.Required(),
		),
	}

	gsheetsSection := warg.NewSection(
		"Google Sheets commands",
		warg.NewSubCmd(
			"open",
			"Open spreadsheet in browser",
			gSheetsOpen,
			warg.CmdFlagMap(sheetFlags),
		),
		warg.NewSubCmd(
			"upload",
			"Upload CSV to Google Sheets. This will overwrite whatever is in the spreadsheet",
			gSheetsUpload,
			warg.CmdFlagMap(sheetFlags),
			warg.NewCmdFlag(
				"--csv-path",
				"CSV file to upload",
				scalar.Path(),
				warg.Required(),
			),
			warg.NewCmdFlag(
				"--timeout",
				"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
				scalar.Duration(
					scalar.Default(time.Minute*10),
				),
				warg.Required(),
			),
		),
	)

	searchCmd := warg.NewCmd(

		"Full text search SQLite database",
		search,
		warg.NewCmdFlag(
			"--limit",
			"Max number of results",
			scalar.Int(
				scalar.Default(50),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name.",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
			warg.Required(),
		),
		warg.NewCmdFlag(
			"--term",
			"Search for this term",
			scalar.String(),
			warg.Alias("-t"),
			warg.Required(),
		),

		// TODO: how many results? limit by date added?
	)

	app := warg.New(
		"starghaze",
		"v1.0.0",
		warg.NewSection(
			"Save GitHub Starred Repos",
			warg.SubCmd(
				"download",
				downloadCmd,
			),
			warg.SubCmd(
				"format",
				formatCmd,
			),
			warg.SubCmd(
				"search",
				searchCmd,
			),
			warg.SubSection(
				"gsheets",
				gsheetsSection,
			),
			warg.SectionFooter("Homepage: https://github.com/bbkane/starghaze"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
