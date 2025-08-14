package warg

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type GoldenTestArgs struct {
	App *App

	// UpdateGolden files for captured stderr/stdout
	UpdateGolden bool

	// Whether the action should return an error
	ExpectActionErr bool
}

// GoldenTest runs the app and and captures stdout and stderr into files.
// If those differ than previously captured stdout/stderr,
// t.Fatalf will be called.
//
// Passed `parseOpts` should not include OverrideStderr/OverrideStdout as GoldenTest overwrites those
func GoldenTest(
	t *testing.T,
	args GoldenTestArgs,
	parseOpts ...ParseOpt) {
	stderrTmpFile, err := os.CreateTemp(os.TempDir(), "warg-test-")
	require.Nil(t, err)

	stdoutTmpFile, err := os.CreateTemp(os.TempDir(), "warg-test-")
	require.Nil(t, err)

	err = args.App.Validate()
	require.Nil(t, err)

	parseOpts = append(parseOpts, Stderr(stderrTmpFile))
	parseOpts = append(parseOpts, Stdout(stdoutTmpFile))
	pr, parseErr := args.App.Parse(parseOpts...)

	// parseOptHolder := cli.NewParseOptHolder(parseOpts...)
	// parseopt.OverrideStderr(stderrTmpFile)(&parseOptHolder)
	// cli.OverrideStdout(stdoutTmpFile)(&parseOptHolder)
	// pr, parseErr := args.App.parseWithOptHolder2(parseOptHolder)

	require.Nil(t, parseErr)

	actionErr := pr.Action(pr.Context)
	if args.ExpectActionErr {
		require.Error(t, actionErr)
	} else {
		require.NoError(t, actionErr)
	}

	stderrCloseErr := stderrTmpFile.Close()
	require.Nil(t, stderrCloseErr)

	stdoutCloseErr := stdoutTmpFile.Close()
	require.Nil(t, stdoutCloseErr)

	actualStderrBytes, stderrReadErr := os.ReadFile(stderrTmpFile.Name())
	require.Nil(t, stderrReadErr)

	actualStdoutBytes, stoutReadErr := os.ReadFile(stdoutTmpFile.Name())
	require.Nil(t, stoutReadErr)

	goldenDir := filepath.Join("testdata", t.Name())

	stderrGoldenFilePath := filepath.Join(goldenDir, "stderr.golden.txt")
	stderrGoldenFilePath, err = filepath.Abs(stderrGoldenFilePath)
	require.Nil(t, err)

	stdoutGoldenFilePath := filepath.Join(goldenDir, "stdout.golden.txt")
	stdoutGoldenFilePath, err = filepath.Abs(stdoutGoldenFilePath)
	require.Nil(t, err)

	if args.UpdateGolden {
		mkdirErr := os.MkdirAll(goldenDir, 0700)
		require.Nil(t, mkdirErr)

		stderrWriteErr := os.WriteFile(stderrGoldenFilePath, actualStderrBytes, 0600)
		require.Nil(t, stderrWriteErr)
		t.Logf("Wrote: %v\n", stderrGoldenFilePath)

		stdoutWriteErr := os.WriteFile(stdoutGoldenFilePath, actualStdoutBytes, 0600)
		require.Nil(t, stdoutWriteErr)
		t.Logf("Wrote: %v\n", stdoutGoldenFilePath)
	}

	stderrExpectedBytes, stderrExpectedReadErr := os.ReadFile(stderrGoldenFilePath)
	require.Nil(t, stderrExpectedReadErr, "actualBytes: \n%s", string(actualStderrBytes))

	if !bytes.Equal(stderrExpectedBytes, actualStderrBytes) {
		t.Fatalf(
			"expected != actual. See diff:\n  vimdiff %s %s\n",
			stderrGoldenFilePath,
			stderrTmpFile.Name(),
		)
	}

	stdoutExpectedBytes, stdoutExpectedReadErr := os.ReadFile(stdoutGoldenFilePath)
	require.Nil(t, stdoutExpectedReadErr, "actualBytes: \n%s", string(actualStdoutBytes))

	if !bytes.Equal(stdoutExpectedBytes, actualStdoutBytes) {
		t.Fatalf(
			"expected != actual. See diff:\n  vimdiff %s %s\n",
			stdoutGoldenFilePath,
			stdoutTmpFile.Name(),
		)
	}

}
