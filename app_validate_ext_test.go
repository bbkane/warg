package warg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

func TestApp_Validate(t *testing.T) {

	tests := []struct {
		name        string
		app         warg.App // NOTE:
		expectedErr bool
	}{
		{
			name: "leafSection",
			app: warg.New(
				"newAppName",
				section.New("Help for section"),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		// app.Validate should allow app names with dashes
		{
			name: "appNameWithDash",
			app: warg.New(
				"newAppName",
				section.New("",
					section.Command("com", "command for validation", command.DoNothing),
				),
				warg.SkipValidation(),
			),
			expectedErr: false,
		},
		{
			name: "sectionNameWithDash",
			app: warg.New(
				"newAppName",
				section.New("",
					section.Section("-name", "",
						section.Command("com", "command for validation", command.DoNothing),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "commandNameWithDash",
			app: warg.New(
				"newAppName",
				section.New("",
					section.Section("name", "",
						section.Command("-com", "starts with dash", command.DoNothing),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "flagNameNoDash",
			app: warg.New(
				"newAppName",
				section.New("",
					section.Flag("f", "", nil),
					section.Command("c", "", nil),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "aliasNameNoDash",
			app: warg.New(
				"newAppName",
				section.New("",
					section.Flag("-f", "", value.Bool,
						flag.Alias("f"),
					),
					section.Command("c", "", nil),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},

		{
			name: "flagNameAliasConflict",
			app: warg.New(
				"newAppName",
				section.New("",
					section.Flag("-f", "", value.Bool),
					section.Command("c", "", command.DoNothing,
						command.Flag("--other", "", value.Bool, flag.Alias("-f")),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualErr := tt.app.Validate()

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
		})
	}
}
