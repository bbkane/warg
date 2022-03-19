package warg_test

import (
	"bytes"
	stdlibflag "flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

var update = stdlibflag.Bool("update", false, "update golden files")

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
					value.String,
					flag.Default("vi"),
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

func tmpFile(t *testing.T) *os.File {
	actualHelpTmpFile, err := ioutil.TempFile(os.TempDir(), "warg-test-")
	if err != nil {
		t.Fatalf("Error creating tmpfile: %v", err)
	}
	return actualHelpTmpFile
}

func TestAppHelp(t *testing.T) {
	tests := []struct {
		name     string
		app      warg.App
		args     []string
		lookup   warg.LookupFunc
		helpFile *os.File
	}{

		// allcommands
		{
			name: "allcommandsCommand",
			app: warg.New(
				"allcommandsCommand",
				grabbitSection(),
				warg.OverrideHelpFlag(
					[]help.HelpFlagMapping{
						{Name: "allcommands", CommandHelp: help.AllCommandsCommandHelp, SectionHelp: help.AllCommandsSectionHelp},
					},
					tmpFile(t),
					"--help",
					"Print help information",
					flag.Default("allcommands"),
					flag.Alias("-h"),
				),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "config", "edit", "--help"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name: "allcommandsSection",
			app: warg.New(
				"allcommandsSection",
				grabbitSection(),
				warg.OverrideHelpFlag(
					[]help.HelpFlagMapping{
						{Name: "allcommands", CommandHelp: help.AllCommandsCommandHelp, SectionHelp: help.AllCommandsSectionHelp},
					},
					tmpFile(t),
					"--help",
					"Print help information",
					flag.Default("allcommands"),
					flag.Alias("-h"),
				),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "--help"},
			lookup: warg.LookupMap(nil),
		},

		// detailed
		{
			name: "detailedCommand",
			app: warg.New(
				"detailedCommand",
				grabbitSection(),
				warg.OverrideHelpFlag(
					[]help.HelpFlagMapping{
						{Name: "detailed", CommandHelp: help.DetailedCommandHelp, SectionHelp: help.DetailedSectionHelp},
					},
					tmpFile(t),
					"--help",
					"Print help information",
					flag.Default("detailed"),
					flag.Alias("-h"),
				),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "config", "edit", "--help"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name: "detailedSection",
			app: warg.New(
				"detailedSection",
				grabbitSection(),
				warg.OverrideHelpFlag(
					[]help.HelpFlagMapping{
						{Name: "detailed", CommandHelp: help.DetailedCommandHelp, SectionHelp: help.DetailedSectionHelp},
					},
					tmpFile(t),
					"--help",
					"Print help information",
					flag.Default("detailed"),
					flag.Alias("-h"),
				),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "--help"},
			lookup: warg.LookupMap(nil),
		},

		// outline
		{
			name: "outlineCommand",
			app: warg.New(
				"outlineCommand",
				grabbitSection(),
				warg.OverrideHelpFlag(
					[]help.HelpFlagMapping{
						{Name: "outline", CommandHelp: help.OutlineCommandHelp, SectionHelp: help.OutlineSectionHelp},
					},
					tmpFile(t),
					"--help",
					"Print help information",
					flag.Default("outline"),
					flag.Alias("-h"),
				),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "config", "edit", "--help"},
			lookup: warg.LookupMap(map[string]string{"EDITOR": "emacs"}),
		},
		{
			name: "outlineSection",
			app: warg.New(
				"outlineSection",
				grabbitSection(),
				warg.OverrideHelpFlag(
					[]help.HelpFlagMapping{
						{Name: "outline", CommandHelp: help.OutlineCommandHelp, SectionHelp: help.OutlineSectionHelp},
					},
					tmpFile(t),
					"--help",
					"Print help information",
					flag.Default("outline"),
					flag.Alias("-h"),
				),
				warg.SkipValidation(),
			),
			args:   []string{"grabbit", "--help"},
			lookup: warg.LookupMap(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Validate()
			require.Nil(t, err)

			pr, parseErr := tt.app.Parse(tt.args, tt.lookup)
			require.Nil(t, parseErr)

			actionErr := pr.Action(pr.PassedFlags)
			require.Nil(t, actionErr)

			closeErr := tt.app.HelpFile.Close()
			require.Nil(t, closeErr)

			actualHelpBytes, readErr := ioutil.ReadFile(tt.app.HelpFile.Name())
			require.Nil(t, readErr)

			goldenDir := filepath.Join("testdata", t.Name())
			goldenFilePath := filepath.Join(goldenDir, "golden.txt")
			goldenFilePath, err = filepath.Abs(goldenFilePath)
			require.Nil(t, err)

			if *update {
				mkdirErr := os.MkdirAll(goldenDir, 0700)
				require.Nil(t, mkdirErr)

				writeErr := ioutil.WriteFile(goldenFilePath, actualHelpBytes, 0600)
				require.Nil(t, writeErr)

				t.Logf("Wrote: %v\n", goldenFilePath)
			}

			expectedBytes, expectedReadErr := ioutil.ReadFile(goldenFilePath)
			require.Nil(t, expectedReadErr)

			if !bytes.Equal(expectedBytes, actualHelpBytes) {
				t.Fatalf(
					"expected != actual. See diff:\n  vimdiff %s %s\n",
					goldenFilePath,
					tt.app.HelpFile.Name(),
				)
			}

		})
	}
}
