package warg

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// GoldenTest runs the app and and captures stdout and stderr into files.
// If those differ than previously captured stdout/stderr,
// t.Fatalf will be called. Pass updateGolden = true to update captured stdout and stderr files under the ./testdata dir (relative to test location)..
// From the CLI, call with WARG_TEST_UPDATE_GOLDEN=1 go test ./...
func GoldenTest(t *testing.T, app App, args []string, lookup LookupFunc, updateGolden bool) {
	stderrTmpFile, err := os.CreateTemp(os.TempDir(), "warg-test-")
	require.Nil(t, err)

	stdoutTmpFile, err := os.CreateTemp(os.TempDir(), "warg-test-")
	require.Nil(t, err)

	err = app.Validate()
	require.Nil(t, err)

	pr, parseErr := app.Parse(
		OverrideArgs(args),
		OverrideLookupFunc(lookup),
		OverrideStderr(stderrTmpFile),
		OverrideStdout(stdoutTmpFile),
	)
	require.Nil(t, parseErr)

	actionErr := pr.Action(pr.Context)
	require.Nil(t, actionErr)

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

	if updateGolden {
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
