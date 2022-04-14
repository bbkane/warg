package warg_test

import (
	"testing"

	"github.com/alecthomas/assert"
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
				section.New("Help for section"),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		// app.Validate should allow app names with dashes
		{
			name: "appNameWithDash",
			app: warg.New(
				section.New("",
					section.Command("com", "command for validation", command.DoNothing),
				),
				warg.SkipValidation(),
				warg.Name("-"+t.Name()),
			),
			expectedErr: false,
		},
		{
			name: "sectionNameWithDash",
			app: warg.New(
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
				assert.NotNil(t, actualErr)
				return
			} else {
				assert.Nil(t, actualErr)
			}
		})
	}
}
