package warg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
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
					section.Command(
						"c",
						"",
						command.DoNothing,
						command.Flag("f", "", nil),
					),
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
					section.Command(
						"c",
						"",
						command.DoNothing,
						command.Flag("-f", "", scalar.Bool(),
							flag.Alias("f"),
						)),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},

		{
			name: "commandFlagAliasCommandFlagNameConflict",
			app: warg.New(
				"newAppName",
				section.New("",
					section.Command("c", "", command.DoNothing,
						command.Flag("-f", "", scalar.Bool()),
						command.Flag("--other", "", scalar.Bool(), flag.Alias("-f")),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagAliasGlobalFlagAliasConflict",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
						command.Flag(
							"--commandflag",
							"global flag conflict",
							scalar.String(),
							flag.Alias("--global"),
						),
					),
				),
				warg.SkipValidation(),
				warg.GlobalFlag(
					"--globalflag",
					"global flag",
					scalar.String(),
					flag.Alias("--global"),
				),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagAliasGlobalFlagNameConflict",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
						command.Flag(
							"--commandflag",
							"global flag conflict",
							scalar.String(),
							flag.Alias("--global"),
						),
					),
				),
				warg.SkipValidation(),
				warg.GlobalFlag(
					"--global",
					"global flag",
					scalar.String(),
				),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagNameGlobalFlagNameConflict",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
						command.Flag(
							"--global",
							"global flag conflict",
							scalar.String(),
						),
					),
				),
				warg.SkipValidation(),
				warg.GlobalFlag(
					"--global",
					"global flag",
					scalar.String(),
				),
			),
			expectedErr: true,
		},
		{
			name: "commandNameSectionNameConflict",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Command(
						"conflict",
						"help for com",
						command.DoNothing,
					),
					section.Section(
						"conflict",
						"help for section",
					),
				),
				warg.SkipValidation(),
				warg.GlobalFlag(
					"--global",
					"global flag",
					scalar.String(),
				),
			),
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualErr := tt.app.Validate()

			if tt.expectedErr {
				require.Error(t, actualErr)
				return
			} else {
				require.NoError(t, actualErr)
			}
		})
	}
}
