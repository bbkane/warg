package warg

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// GoldenTestArgs holds configuration for [GoldenTest].
type GoldenTestArgs struct {
	App *App

	// UpdateGolden overwrites existing golden files with actual output when true.
	UpdateGolden bool

	// ExpectActionErr asserts the action returns a non-nil error when true.
	ExpectActionErr bool

	// Args are the command-line arguments to parse (without program name).
	Args []string
}

// GoldenTest runs the app with the given args and compares stdout/stderr against
// golden files in testdata/<TestName>/. If the output differs, t.Fatalf is called.
// Set GoldenTestArgs.UpdateGolden to overwrite golden files with actual output.
//
// Do not pass [ParseWithStderr] or [ParseWithStdout] in parseOpts, as GoldenTest
// overrides those to capture output.
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

	parseOpts = append(parseOpts, ParseWithStderr(stderrTmpFile))
	parseOpts = append(parseOpts, ParseWithStdout(stdoutTmpFile))
	pr, parseErr := args.App.Parse(args.Args, parseOpts...)

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
