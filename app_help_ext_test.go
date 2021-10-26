package warg_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/help"

	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func RequireEqualBytesOrDiff(t *testing.T, expectedFilePath string, actual []byte, msg string) {
	expectedBytes, readErr := ioutil.ReadFile(expectedFilePath)
	require.Nil(t, readErr)
	if bytes.Equal(expectedBytes, actual) {
		return
	}

	actualTmpFile, err := ioutil.TempFile(os.TempDir(), "go-test-actual-")
	if err != nil {
		t.Fatalf("Error creating tmpfile: %v", err)
	}
	defer actualTmpFile.Close()

	_, err = actualTmpFile.Write(actual)
	if err != nil {
		t.Fatalf("Error writing tmpfile: %v", err)
	}

	t.Fatalf(
		"%s: expected != actual. See diff:\n  vimdiff %s %s\n",
		msg,
		expectedFilePath,
		actualTmpFile.Name(),
	)
}

func TestDefaultSectionHelp(t *testing.T) {
	var actualBuffer bytes.Buffer

	app := warg.New(
		"grabbit",
		s.New(
			"grab those images!",
			s.WithSection(
				"config",
				"change grabbit's config",
				s.WithCommand(
					"edit",
					"edit the config",
					c.DoNothing,
					c.WithFlag(
						"--editor",
						"path to editor",
						v.String,
						f.Default("vi"),
					),
				),
			),
			s.WithCommand(
				"grab",
				"do the grabbity grabbity",
				c.DoNothing,
			),
		),
		warg.OverrideHelp(
			&actualBuffer,
			[]string{"-h", "--help"},
			help.DefaultSectionHelp,
			help.DefaultCommandHelp,
		),
	)
	args := []string{"grabbit", "--help"}
	actualErr := app.Run(args, warg.DictLookup(nil))
	require.Nil(t, actualErr)

	golden := filepath.Join("testdata", t.Name()+".golden.txt")
	if *update {
		mkdirErr := os.MkdirAll("testdata", 0700)
		require.Nil(t, mkdirErr)
		writeErr := ioutil.WriteFile(golden, actualBuffer.Bytes(), 0600)
		require.Nil(t, writeErr)
		t.Logf("Wrote: %v\n", golden)
	}

	RequireEqualBytesOrDiff(
		t,
		golden,
		actualBuffer.Bytes(),
		t.Name(),
	)
}

func TestDefaultCommandHelp(t *testing.T) {
	var actualBuffer bytes.Buffer

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
		s.New(
			"grab those images!",
			s.WithSection(
				"config",
				"Change grabbit's config",
				s.Footer(rootFooter),
				s.WithCommand(
					"edit",
					"Edit the config. A default config will be created if it doesn't exist",
					c.DoNothing,
					c.Footer(configEditFooter),
					c.WithFlag(
						"--editor",
						"path to editor",
						v.String,
						f.Default("vi"),
					),
				),
			),
			s.WithCommand(
				"grab",
				"do the grabbity grabbity",
				c.DoNothing,
			),
		),
		warg.OverrideHelp(
			&actualBuffer,
			[]string{"-h", "--help"},
			help.DefaultSectionHelp,
			help.DefaultCommandHelp,
		),
	)
	args := []string{"grabbit", "config", "edit", "--help"}
	actualErr := app.Run(args, warg.DictLookup(nil))
	require.Nil(t, actualErr)

	golden := filepath.Join("testdata", t.Name()+".golden.txt")
	if *update {
		mkdirErr := os.MkdirAll("testdata", 0700)
		require.Nil(t, mkdirErr)
		writeErr := ioutil.WriteFile(golden, actualBuffer.Bytes(), 0600)
		require.Nil(t, writeErr)
		t.Logf("Wrote: %v\n", golden)
	}
	RequireEqualBytesOrDiff(
		t,
		golden,
		actualBuffer.Bytes(),
		t.Name(),
	)
}
