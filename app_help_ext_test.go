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
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestDefaultSectionHelp(t *testing.T) {
	var actualBuffer bytes.Buffer

	app := warg.New(
		"grabbit",
		s.NewSection(
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
						v.StringEmpty,
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
			warg.DefaultSectionHelp,
			warg.DefaultCommandHelp,
		),
	)
	args := []string{"grabbit", "--help"}
	actualErr := app.Run(args)
	require.Nil(t, actualErr)

	golden := filepath.Join("testdata", t.Name()+".golden.txt")
	if *update {
		mkdirErr := os.MkdirAll("testdata", 0700)
		require.Nil(t, mkdirErr)
		writeErr := ioutil.WriteFile(golden, actualBuffer.Bytes(), 0600)
		require.Nil(t, writeErr)
		t.Logf("Wrote: %v\n", golden)
	}

	expectedBytes, readErr := ioutil.ReadFile(golden)
	require.Nil(t, readErr)

	require.Equal(t, expectedBytes, actualBuffer.Bytes())
}

func TestDefaultCommandHelp(t *testing.T) {
	var actualBuffer bytes.Buffer

	app := warg.New(
		"grabbit",
		s.NewSection(
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
						v.StringEmpty,
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
			warg.DefaultSectionHelp,
			warg.DefaultCommandHelp,
		),
	)
	args := []string{"grabbit", "config", "edit", "--help"}
	actualErr := app.Run(args)
	require.Nil(t, actualErr)

	golden := filepath.Join("testdata", t.Name()+".golden.txt")
	if *update {
		mkdirErr := os.MkdirAll("testdata", 0700)
		require.Nil(t, mkdirErr)
		writeErr := ioutil.WriteFile(golden, actualBuffer.Bytes(), 0600)
		require.Nil(t, writeErr)
		t.Logf("Wrote: %v\n", golden)
	}

	expectedBytes, readErr := ioutil.ReadFile(golden)
	require.Nil(t, readErr)

	require.Equal(t, expectedBytes, actualBuffer.Bytes())
}
