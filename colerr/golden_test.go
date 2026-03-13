package colerr

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type goldenTestParams struct {
	// TmpFilePrefix is prepended to tmpfiles goldenTest creates as debugging convenience.
	TmpFilePrefix string

	// FileNames for goldenTest to create and work to write to. Example:
	//	[]string{"stdout.txt", "stderr.txt"}
	FileNames []string

	// GoldenDir holds saved GoldenTeests.
	// Example:
	//	filepath.Join("testdata", t.Name())
	GoldenDir string

	// UpdateEnvVar is checked, and, is set to any value, will update golden files
	UpdateEnvVar string

	// WorkFunc runs the code to be tested. Workfunc should write to file handles retrieved from the passed map with values from FileNames.
	// Example:
	//	func(files map[string]*os.File){
	//		f := files["stdout.txt"]
	//		fmt.Fprint(f, "hello")
	//	}
	WorkFunc func(map[string]*os.File)
}

// goldenTest provides files to p.WorkFunc, then compares those files to previously
// saved ones and provides vimdiff commands to inspect any diffences found.
// In pseudocode for a one-file version of this function:
//
//	tmpFile := NewTmpFile()
//	Write(tmpFile)  // work
//	Close(tmpFile)
//	actualBytes := Read(tmpFile.Name())
//	var goldenFilePath
//	if update {
//	    Write(actualBytes, goldenFilePath)
//	}
//	expectedBytes := Read(goldenFilePath)
//	Compare(expectedBytes, actualBytes)
func goldenTest(t *testing.T, p goldenTestParams) {

	update := os.Getenv(p.UpdateEnvVar) != ""
	if !update {
		t.Logf("To update golden files, run:\n  %s=1 go test ./...", p.UpdateEnvVar)
	}

	var tmpFiles = make(map[string]*os.File, len(p.FileNames))
	for _, name := range p.FileNames {
		file, err := os.CreateTemp(os.TempDir(), p.TmpFilePrefix+"-"+name)
		require.Nil(t, err)
		t.Logf("wrote tmpfile: %#v", file.Name())
		tmpFiles[name] = file
	}

	p.WorkFunc(tmpFiles)

	for _, name := range p.FileNames {
		tmpFile := tmpFiles[name]
		err := tmpFile.Close()
		require.Nil(t, err)

		actualBytes, err := os.ReadFile(tmpFile.Name())
		require.Nil(t, err)

		goldenFilePath := filepath.Join(p.GoldenDir, "golden-"+name)
		goldenFilePath, err = filepath.Abs(goldenFilePath)
		require.Nil(t, err)

		if update {
			err = os.MkdirAll(p.GoldenDir, 0700)
			require.Nil(t, err)

			err = os.WriteFile(goldenFilePath, actualBytes, 0600)
			require.Nil(t, err)
			t.Logf("wrote golden file: %#v\n", goldenFilePath)
		}

		expectedBytes, err := os.ReadFile(goldenFilePath)
		require.Nil(t, err)

		if !bytes.Equal(expectedBytes, actualBytes) {
			t.Logf(
				"%s: expected != actual. See diff:\n  vimdiff %s %s\n\n",
				name,
				goldenFilePath,
				tmpFile.Name(),
			)
			t.Fail()
		}
	}

}
