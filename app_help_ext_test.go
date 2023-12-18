package warg_test

// Run WARG_TEST_UPDATE_GOLDEN=1 go test ./... to update golden files

import (
	"os"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
)

// A grabbitSection is a simple section to test help
func grabbitSection() section.SectionT {

	rootFooter := `Examples:

	# Grab without config
	grabbit grab

	# Edit config, then grab
	grabbit config edit
	grabbit grab
	`

	configEditFooter := `Examples:

	# Use defaults
	grabbit config edit

	# Override defaults
	grabbit config edit --config-path /path/to/config --editor code
	`

	sec := section.New(
		"grab those images!",
		section.Command(
			"grab",
			"do the grabbity grabbity",
			command.DoNothing,
		),
		section.Command(
			"command2",
			"another command",
			command.DoNothing,
		),
		section.Command(
			"command3",
			"another command",
			command.DoNothing,
		),
		section.Section(
			"config",
			"Change grabbit's config",
			section.Footer(rootFooter),
			section.Command(
				"edit",
				"Edit the config. A default config will be created if it doesn't exist",
				command.DoNothing,
				command.Footer(configEditFooter),
				command.Flag(
					"--editor",
					"path to editor",
					scalar.String(
						scalar.Default("vi"),
					),
					flag.ConfigPath("editor"),
					flag.EnvVars("EDITOR"),
					flag.Required(),
				),
			),
		),
		section.Section(
			"section2",
			"another section",
			section.Command("com", "Dummy command to pass validation", command.DoNothing),
		),
		section.Section(
			"section3",
			"another section",
			section.Command("com", "Dummy command to pass validation", command.DoNothing),
		),
	)
	return sec
}

func TestAppHelp(t *testing.T) {
	updateGolden := os.Getenv("WARG_TEST_UPDATE_GOLDEN") != ""
	tests := []struct {
		name   string
		app    warg.App
		args   []string
		lookup warg.LookupFunc
	}{

		// allcommands (no command help)
		{
			name: "allcommandsSection",
			app: warg.New(
				"grabbit",
				grabbitSection(),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "config", "--help"},
			lookup: warg.LookupMap(nil),
		},

		// detailed
		{
			name: "detailedCommand",
			app: warg.New(
				"newAppName",
				grabbitSection(),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "config", "edit", "--help"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name: "detailedSection",
			app: warg.New(
				"newAppName",
				grabbitSection(),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "--help", "detailed"},
			lookup: warg.LookupMap(nil),
		},

		// outline
		{
			name: "outlineCommand",
			app: warg.New(
				"grabbit",
				grabbitSection(),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "config", "edit", "--help", "outline"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name: "outlineSection",
			app: warg.New(
				"grabbit",
				grabbitSection(),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "--help", "outline"},
			lookup: warg.LookupMap(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warg.GoldenTest(
				t,
				tt.app,
				updateGolden,
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(tt.lookup),
			)
		})
	}
}
