package warg_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	wflag "github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/help"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func RequireEqualBytesOrDiff(t *testing.T, expectedFilePath string, actualFilePath string, msg string) {
	expectedBytes, expectedReadErr := ioutil.ReadFile(expectedFilePath)
	require.Nil(t, expectedReadErr)

	actualBytes, actualReadErr := ioutil.ReadFile(actualFilePath)
	require.Nil(t, actualReadErr)

	if bytes.Equal(expectedBytes, actualBytes) {
		return
	}

	t.Fatalf(
		"%s: expected != actual. See diff:\n  vimdiff %s %s\n",
		msg,
		expectedFilePath,
		actualFilePath,
	)
}

func TestDefaultSectionHelp(t *testing.T) {

	actualHelpTmpFile, err := ioutil.TempFile(os.TempDir(), "go-test-actual-help")
	if err != nil {
		t.Fatalf("Error creating tmpfile: %v", err)
	}

	app := warg.New(
		"grabbit",
		section.New(
			"grab those images!",
			section.Section(
				"config",
				"change grabbit's config",
				section.Command(
					"edit",
					"edit the config",
					command.DoNothing,
					command.Flag(
						"--editor",
						"path to editor",
						value.String,
						wflag.Default("vi"),
					),
				),
			),
			section.Command(
				"grab",
				"do the grabbity grabbity",
				command.DoNothing,
			),
		),
		warg.OverrideHelpFlag(
			[]warg.HelpFlagMapping{
				{Name: "default", CommandHelp: help.DefaultCommandHelp, SectionHelp: help.DefaultSectionHelp},
			},
			actualHelpTmpFile,
			"--help",
			"Print help information",
			wflag.Default("default"),
			wflag.Alias("-h"),
		),
	)
	args := []string{"grabbit", "--help"}
	pr, parseErr := app.Parse(args, warg.LookupMap(nil))
	require.Nil(t, parseErr)
	actualErr := pr.Action(pr.PassedFlags)
	require.Nil(t, actualErr)

	closeErr := actualHelpTmpFile.Close()
	require.Nil(t, closeErr)

	actualHelpBytes, readErr := ioutil.ReadFile(actualHelpTmpFile.Name())
	require.Nil(t, readErr)

	golden := filepath.Join("testdata", t.Name()+".golden.txt")
	if *update {
		mkdirErr := os.MkdirAll("testdata", 0700)
		require.Nil(t, mkdirErr)

		writeErr := ioutil.WriteFile(golden, actualHelpBytes, 0600)
		require.Nil(t, writeErr)

		t.Logf("Wrote: %v\n", golden)
	}

	RequireEqualBytesOrDiff(
		t,
		golden,
		actualHelpTmpFile.Name(),
		t.Name(),
	)
}

func TestDefaultCommandHelp(t *testing.T) {

	actualHelpTmpFile, err := ioutil.TempFile(os.TempDir(), "go-test-actual-help")
	if err != nil {
		t.Fatalf("Error creating tmpfile: %v", err)
	}

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

	app := warg.New(
		"grabbit",
		section.New(
			"grab those images!",
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
						wflag.Default("vi"),
						wflag.ConfigPath("editor"),
						wflag.EnvVars("EDITOR"),
						wflag.Required(),
					),
				),
			),
			section.Command(
				"grab",
				"do the grabbity grabbity",
				command.DoNothing,
			),
		),
		warg.OverrideHelpFlag(
			[]warg.HelpFlagMapping{
				{Name: "default", CommandHelp: help.DefaultCommandHelp, SectionHelp: help.DefaultSectionHelp},
			},
			actualHelpTmpFile,
			"--help",
			"Print help information",
			wflag.Default("default"),
			wflag.Alias("-h"),
		),
	)
	args := []string{"grabbit", "config", "edit", "--help"}
	pr, parseErr := app.Parse(args, warg.LookupMap(map[string]string{"EDITOR": "emacs"}))
	require.Nil(t, parseErr)
	actualErr := pr.Action(pr.PassedFlags)
	require.Nil(t, actualErr)

	closeErr := actualHelpTmpFile.Close()
	require.Nil(t, closeErr)

	actualHelpBytes, readErr := ioutil.ReadFile(actualHelpTmpFile.Name())
	require.Nil(t, readErr)

	golden := filepath.Join("testdata", t.Name()+".golden.txt")
	if *update {
		mkdirErr := os.MkdirAll("testdata", 0700)
		require.Nil(t, mkdirErr)

		writeErr := ioutil.WriteFile(golden, actualHelpBytes, 0600)
		require.Nil(t, writeErr)

		t.Logf("Wrote: %v\n", golden)
	}
	RequireEqualBytesOrDiff(
		t,
		golden,
		actualHelpTmpFile.Name(),
		t.Name(),
	)
}
