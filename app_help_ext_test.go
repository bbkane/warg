package warg_test

// Run WARG_TEST_UPDATE_GOLDEN=1 go test ./... to update golden files

import (
	"os"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/parseopt"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

// A grabbitSection is a simple section to test help
func grabbitSection() wargcore.Section {

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

	sec := wargcore.NewSection(
		"grab those images!",
		wargcore.NewChildCmd(
			"grab",
			"do the grabbity grabbity",
			wargcore.DoNothing,
		),
		wargcore.NewChildCmd(
			"command2",
			"another command",
			wargcore.DoNothing,
		),
		wargcore.NewChildCmd(
			"command3",
			"another command",
			wargcore.DoNothing,
		),
		wargcore.NewChildSection(
			"config",
			"Change grabbit's config",
			wargcore.SectionFooter(rootFooter),
			wargcore.NewChildCmd(
				"edit",
				"Edit the config. A default config will be created if it doesn't exist",
				wargcore.DoNothing,
				wargcore.CmdFooter(configEditFooter),
				wargcore.NewChildFlag(
					"--editor",
					"path to editor",
					scalar.String(
						scalar.Default("vi"),
					),
					wargcore.ConfigPath("editor"),
					wargcore.EnvVars("EDITOR"),
					wargcore.Required(),
				),
			),
		),
		wargcore.NewChildSection(
			"section2",
			"another section",
			wargcore.NewChildCmd("com", "Dummy command to pass validation", wargcore.DoNothing),
		),
		wargcore.NewChildSection(
			"section3",
			"another section",
			wargcore.NewChildCmd("com", "Dummy command to pass validation", wargcore.DoNothing),
		),
	)
	return sec
}

func TestAppHelp(t *testing.T) {
	updateGolden := os.Getenv("WARG_TEST_UPDATE_GOLDEN") != ""
	tests := []struct {
		name   string
		args   []string
		lookup wargcore.LookupEnv
	}{
		// toplevel just a toplevel help!
		{
			name:   "toplevel",
			args:   []string{"grabbit", "-h", "outline"},
			lookup: wargcore.LookupMap(nil),
		},

		// allcommands (no command help)
		{
			name:   "allcommandsSection",
			args:   []string{"grabbit", "config", "--help"},
			lookup: wargcore.LookupMap(nil),
		},

		// detailed
		{
			name:   "detailedCommand",
			args:   []string{"grabbit", "config", "edit", "--help"},
			lookup: wargcore.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name:   "detailedSection",
			args:   []string{"grabbit", "--help", "detailed"},
			lookup: wargcore.LookupMap(nil),
		},

		// outline
		{
			// TODO: make this print global flags!
			name:   "outlineCommand",
			args:   []string{"grabbit", "config", "edit", "--help", "outline"},
			lookup: wargcore.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			// TODO: make this print global flags!
			name:   "outlineSection",
			args:   []string{"grabbit", "--help", "outline"},
			lookup: wargcore.LookupMap(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := warg.New(
				"grabbit",
				"v1.0.0",
				grabbitSection(),
				warg.SkipValidation(),
			)
			warg.GoldenTest(
				t,
				warg.GoldenTestArgs{
					App:             &app,
					UpdateGolden:    updateGolden,
					ExpectActionErr: false,
				},
				parseopt.Args(tt.args),
				parseopt.LookupEnv(tt.lookup),
			)
		})
	}
}
