package main

import (
	"time"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func app() *wargcore.App {

	downloadCmd := command.NewCmd(
		"Download star info",
		githubStarsDownload,
		command.NewChildFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
		),
		command.NewChildFlag(
			"--max-languages",
			"Max number of languages to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		command.NewChildFlag(
			"--max-repo-topics",
			"Max number of topics to query on a repo",
			scalar.Int(
				scalar.Default(20),
			),
		),
		command.NewChildFlag(
			"--after-cursor",
			"PageInfo EndCursor to start from",
			scalar.String(),
		),
		command.NewChildFlag(
			"--max-pages",
			"Max number of pages to fetch",
			scalar.Int(
				scalar.Default(1),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--output",
			"Output filepath. Must not exist",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
		),
		command.NewChildFlag(
			"--page-size",
			"Number of starred repos in page",
			scalar.Int(
				scalar.Default(100),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--timeout",
			"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
			scalar.Duration(
				scalar.Default(time.Minute*10),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--token",
			"Github PAT",
			scalar.String(),
			flag.EnvVars("STARGHAZE_GITHUB_TOKEN", "GITHUB_TOKEN"),
			flag.Required(),
		),
	)

	formatCmd := command.NewCmd(
		"Format downloaded GitHub Stars",
		format,
		command.NewChildFlag(
			"--format",
			"Output format",
			scalar.String(
				scalar.Choices("csv", "jsonl", "sqlite", "zinc"),
				scalar.Default("csv"),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--date-format",
			"Datetime output format. See https://github.com/lestrrat-go/strftime for details. If not passed, the GitHub default is RFC 3339. Consider using '%b %d, %Y' for csv format",
			scalar.String(),
		),
		command.NewChildFlag(
			"--include-readmes",
			"Search for README.md.",
			scalar.Bool(
				scalar.Default(false),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name. Only used for --format sqlite",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
		),
		command.NewChildFlag(
			"--zinc-index-name",
			"Only used for --format zinc.",
			scalar.String(
				scalar.Default("starghaze"),
			),
		),
		command.NewChildFlag(
			"--input",
			"Input file",
			scalar.Path(
				scalar.Default(path.New("starghaze_download.jsonl")),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--max-line-size",
			"Max line size in the file in MB",
			scalar.Int(
				scalar.Default(10),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--output",
			"output file. Prints to stdout if not passed",
			scalar.Path(),
		),
	)

	sheetFlags := wargcore.FlagMap{
		"--sheet-id": flag.New(
			"ID For the particulare sheet. Viewable from `gid` URL param",
			scalar.Int(),
			flag.EnvVars("STARGHAZE_SHEET_ID"),
			flag.Required(),
		),
		"--spreadsheet-id": flag.New(
			"ID for the whole spreadsheet. Viewable from URL",
			scalar.String(),
			flag.EnvVars("STARGHAZE_SPREADSHEET_ID"),
			flag.Required(),
		),
	}

	gsheetsSection := section.NewSection(
		"Google Sheets commands",
		section.NewChildCmd(
			"open",
			"Open spreadsheet in browser",
			gSheetsOpen,
			command.ChildFlagMap(sheetFlags),
		),
		section.NewChildCmd(
			"upload",
			"Upload CSV to Google Sheets. This will overwrite whatever is in the spreadsheet",
			gSheetsUpload,
			command.ChildFlagMap(sheetFlags),
			command.NewChildFlag(
				"--csv-path",
				"CSV file to upload",
				scalar.Path(),
				flag.Required(),
			),
			command.NewChildFlag(
				"--timeout",
				"Timeout for a run. Use https://pkg.go.dev/time#Duration to build it",
				scalar.Duration(
					scalar.Default(time.Minute*10),
				),
				flag.Required(),
			),
		),
	)

	searchCmd := command.NewCmd(

		"Full text search SQLite database",
		search,
		command.NewChildFlag(
			"--limit",
			"Max number of results",
			scalar.Int(
				scalar.Default(50),
			),
			flag.Required(),
		),
		command.NewChildFlag(
			"--sqlite-dsn",
			"Sqlite DSN. Usually the file name.",
			scalar.String(
				scalar.Default("starghaze.db"),
			),
			flag.Required(),
		),
		command.NewChildFlag(
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
		section.NewSection(
			"Save GitHub Starred Repos",
			section.ChildCmd(
				"download",
				downloadCmd,
			),
			section.ChildCmd(
				"format",
				formatCmd,
			),
			section.ChildCmd(
				"search",
				searchCmd,
			),
			section.ChildSection(
				"gsheets",
				gsheetsSection,
			),
			section.SectionFooter("Homepage: https://github.com/bbkane/starghaze"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
