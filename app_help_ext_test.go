package warg_test

// Run WARG_TEST_UPDATE_GOLDEN=1 go test ./... to update golden files

import (
	"os"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value/scalar"
)

// A grabbitSection is a simple section to test help
func grabbitSection() warg.Section {

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

	sec := warg.NewSection(
		"grab those images!",
		warg.NewSubCmd(
			"grab",
			"do the grabbity grabbity",
			warg.Unimplemented(),
		),
		warg.NewSubCmd(
			"command2",
			"another command",
			warg.Unimplemented(),
		),
		warg.NewSubCmd(
			"command3",
			"another command",
			warg.Unimplemented(),
		),
		warg.NewSubSection(
			"config",
			"Change grabbit's config",
			warg.SectionFooter(rootFooter),
			warg.NewSubCmd(
				"edit",
				"Edit the config. A default config will be created if it doesn't exist",
				warg.Unimplemented(),
				warg.CmdFooter(configEditFooter),
				warg.NewCmdFlag(
					"--editor",
					"path to editor",
					scalar.String(
						scalar.Default("vi"),
					),
					warg.ConfigPath("editor"),
					warg.EnvVars("EDITOR"),
					warg.Required(),
				),
			),
		),
		warg.NewSubSection(
			"section2",
			"another section",
			warg.NewSubCmd("com", "Dummy command to pass validation", warg.Unimplemented()),
		),
		warg.NewSubSection(
			"section3",
			"another section",
			warg.NewSubCmd("com", "Dummy command to pass validation", warg.Unimplemented()),
		),
	)
	return sec
}

// namedenvLikeSection creates a section structure that reproduces the ordering issue from
// https://github.com/bbkane/warg/issues/74 where sibling sections appear between a
// parent section's direct commands and its nested subsection commands.
func namedenvLikeSection() warg.Section {
	return warg.NewSection(
		"Manage environmental secrets",
		warg.NewSubCmd("print-version", "Print version", warg.Unimplemented()),
		warg.NewSubSection(
			"env",
			"Manage environments",
			warg.NewSubCmd("create", "Create an environment", warg.Unimplemented()),
			warg.NewSubCmd("update", "Update an environment", warg.Unimplemented()),
			warg.NewSubSection(
				"print-script",
				"Print scripts",
				warg.NewSubCmd("export", "Print export script", warg.Unimplemented()),
			),
		),
		warg.NewSubSection(
			"keyring",
			"Manage keyring entries",
			warg.NewSubCmd("create", "Create a keyring entry", warg.Unimplemented()),
		),
	)
}

// TestAllCommandsHelp_Ordering tests that allcommands help output shows commands
// in depth-first order (section's own commands, then subsection commands recursively)
// rather than breadth-first order. See https://github.com/bbkane/warg/issues/74
func TestAllCommandsHelp_Ordering(t *testing.T) {
	updateGolden := os.Getenv("WARG_TEST_UPDATE_GOLDEN") != ""
	app := warg.New(
		"namedenv",
		"v1.0.0",
		namedenvLikeSection(),
		warg.SkipAll(),
	)
	warg.GoldenTest(
		t,
		warg.GoldenTestArgs{
			App:             &app,
			UpdateGolden:    updateGolden,
			ExpectActionErr: false,
			Args:            []string{"--help"},
		},
		warg.ParseWithLookupEnv(warg.LookupMap(nil)),
	)
}

func TestAppHelp(t *testing.T) {
	updateGolden := os.Getenv("WARG_TEST_UPDATE_GOLDEN") != ""
	tests := []struct {
		name   string
		args   []string
		lookup warg.LookupEnv
	}{
		// toplevel just a toplevel help!
		{
			name:   "toplevel",
			args:   []string{"-h", "outline"},
			lookup: warg.LookupMap(nil),
		},

		// allcommands (no command help)
		{
			name:   "allcommandsSection",
			args:   []string{"config", "--help"},
			lookup: warg.LookupMap(nil),
		},

		// detailed
		{
			name:   "detailedCommand",
			args:   []string{"config", "edit", "--help"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name:   "detailedSection",
			args:   []string{"--help", "detailed"},
			lookup: warg.LookupMap(nil),
		},

		// outline
		{
			// TODO: make this print global flags!
			name:   "outlineCommand",
			args:   []string{"config", "edit", "--help", "outline"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			// TODO: make this print global flags!
			name:   "outlineSection",
			args:   []string{"--help", "outline"},
			lookup: warg.LookupMap(nil),
		},

		// compact
		{
			name:   "compactCommand",
			args:   []string{"config", "edit", "--help", "compact"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name:   "compactSection",
			args:   []string{"--help", "compact"},
			lookup: warg.LookupMap(nil),
		},
		{
			name:   "compactCommandTermWidth120",
			args:   []string{"config", "edit", "--term-width", "120", "--help", "compact"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name:   "compactSectionTermWidth120",
			args:   []string{"--help", "compact"},
			lookup: warg.LookupMap(map[string]string{"WARG_TERM_WIDTH": "120"}),
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
					Args:            tt.args,
				},

				warg.ParseWithLookupEnv(tt.lookup),
			)
		})
	}
}
