package colerr

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"go.bbkane.com/warg/styles"
)

func TestErrorWithColorStyle(t *testing.T) {
	s := styles.NewEnabledStyles()
	err := NewWrappedf(
		errors.New("inner error"),
		"Start of message: %s: end of message",
		"middle of message",
	)
	goldenTest(t, goldenTestParams{
		TmpFilePrefix: "colerr-errorwithcolorstyle",
		FileNames:     []string{"stderr.txt"},
		GoldenDir:     filepath.Join("testdata", t.Name()),
		UpdateEnvVar:  "WARG_TEST_UPDATE_GOLDEN",
		WorkFunc: func(files map[string]*os.File) {
			Stacktrace(files["stderr.txt"], &s, err)
		},
	})
}

func TestStacktrace(t *testing.T) {
	s := styles.NewEmptyStyles()

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "single_error",
			err:  errors.New("this is an error"),
		},
		{
			name: "wrapped_custom_error",
			err:  NewWrapped(errors.New("wrapped err"), "wrapper msg"),
		},
		{
			name: "wrappedf_custom_error",
			err:  NewWrappedf(errors.New("wrappedf err"), "wrapperf msg: %s", "with arg"),
		},
		{
			name: "wrapped with errors.Join under",
			err: NewWrappedf(
				errors.Join(
					errors.New("first error"),
					errors.New("second error"),
				),
				"wrapperf msg: %s", "with arg",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenTest(t, goldenTestParams{
				TmpFilePrefix: "colerr-stacktrace",
				FileNames:     []string{"stderr.txt"},
				GoldenDir:     filepath.Join("testdata", t.Name()),
				UpdateEnvVar:  "WARG_TEST_UPDATE_GOLDEN",
				WorkFunc: func(files map[string]*os.File) {
					Stacktrace(files["stderr.txt"], &s, tt.err)
				},
			})
		})
	}
}
