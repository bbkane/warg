package main

import (
	"os"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
)

func TestApp_Validate(t *testing.T) {
	app := app()

	if err := app.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestRunHelp(t *testing.T) {
	updateGolden := os.Getenv("WARG_TEST_UPDATE_GOLDEN") != ""
	tests := []struct {
		name   string
		args   []string
		lookup cli.LookupFunc
	}{
		{
			name:   "presentDetailed",
			args:   []string{"butler", "present", "--help", "detailed"},
			lookup: cli.LookupMap(nil),
		},
		{
			name:   "presentBob",
			args:   []string{"butler", "present", "--name", "bob"},
			lookup: cli.LookupMap(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warg.GoldenTest(
				t,
				warg.GoldenTestArgs{
					App:             app(),
					UpdateGolden:    updateGolden,
					ExpectActionErr: false,
				},
				cli.OverrideArgs(tt.args),
				cli.OverrideLookupFunc(tt.lookup),
			)
		})
	}
}
