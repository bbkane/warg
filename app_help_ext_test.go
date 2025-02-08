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
		section.NewCommand(
			"grab",
			"do the grabbity grabbity",
			command.DoNothing,
		),
		section.NewCommand(
			"command2",
			"another command",
			command.DoNothing,
		),
		section.NewCommand(
			"command3",
			"another command",
			command.DoNothing,
		),
		section.NewSection(
			"config",
			"Change grabbit's config",
			section.Footer(rootFooter),
			section.NewCommand(
				"edit",
				"Edit the config. A default config will be created if it doesn't exist",
				command.DoNothing,
				command.Footer(configEditFooter),
				command.NewFlag(
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
		section.NewSection(
			"section2",
			"another section",
			section.NewCommand("com", "Dummy command to pass validation", command.DoNothing),
		),
		section.NewSection(
			"section3",
			"another section",
			section.NewCommand("com", "Dummy command to pass validation", command.DoNothing),
		),
	)
	return sec
}

func TestAppHelp(t *testing.T) {
	updateGolden := os.Getenv("WARG_TEST_UPDATE_GOLDEN") != ""
	tests := []struct {
		name   string
		args   []string
		lookup warg.LookupFunc
	}{
		// toplevel just a toplevel help!
		{
			name:   "toplevel",
			args:   []string{"grabbit", "-h", "outline"},
			lookup: warg.LookupMap(nil),
		},

		// allcommands (no command help)
		{
			name:   "allcommandsSection",
			args:   []string{"grabbit", "config", "--help"},
			lookup: warg.LookupMap(nil),
		},

		// detailed
		{
			name:   "detailedCommand",
			args:   []string{"grabbit", "config", "edit", "--help"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name:   "detailedSection",
			args:   []string{"grabbit", "--help", "detailed"},
			lookup: warg.LookupMap(nil),
		},

		// outline
		{
			// TODO: make this print global flags!
			name:   "outlineCommand",
			args:   []string{"grabbit", "config", "edit", "--help", "outline"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			// TODO: make this print global flags!
			name:   "outlineSection",
			args:   []string{"grabbit", "--help", "outline"},
			lookup: warg.LookupMap(nil),
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
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(tt.lookup),
			)
		})
	}
}
