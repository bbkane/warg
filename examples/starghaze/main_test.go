package main

import (
	"os"
	"testing"

	"go.bbkane.com/warg"
)

func TestApp_Validate(t *testing.T) {
	app := app()

	if err := app.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestRunHelp(t *testing.T) {
	t.Skip("Machine dependent for now - see https://github.com/bbkane/warg/issues/51")
	updateGolden := os.Getenv("WARG_TEST_UPDATE_GOLDEN") != ""
	tests := []struct {
		name   string
		app    *warg.App
		args   []string
		lookup warg.LookupFunc
	}{
		{
			name:   "starghazeDownloadHelpDetailed",
			app:    app(),
			args:   []string{"starghaze", "download", "--help", "detailed"},
			lookup: warg.LookupMap(nil),
		},
		{
			name:   "starghazeFormatHelpDetailed",
			app:    app(),
			args:   []string{"starghaze", "format", "--help", "detailed"},
			lookup: warg.LookupMap(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warg.GoldenTest(
				t,
				*tt.app,
				updateGolden,
				warg.OverrideArgs(tt.args),
				warg.OverrideLookupFunc(tt.lookup),
			)
		})
	}
}
