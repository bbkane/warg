package colerr

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go.bbkane.com/warg/styles"
)

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
			name: "wrapped_fmt_errorf",
			err:  fmt.Errorf("this is a wrapped error: %w", errors.New("this is the inner error")),
		},
		{
			name: "wrapped_custom_error",
			err:  NewWrapped(errors.New("wrapped err"), "wrapper msg"),
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
					Stacktrace(files["stderr.txt"], s, tt.err)
				},
			})
		})
	}
}
