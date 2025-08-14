package warg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
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
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("Help for section", warg.NewChildSection("leaf", "Is empty but shouldn't be")),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		// app.Validate should allow app names with dashes
		{
			name: "appNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewChildCmd("com", "command for validation", warg.DoNothing),
				),
				warg.SkipValidation(),
			),
			expectedErr: false,
		},
		{
			name: "sectionNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewChildSection("-name", "",
						warg.NewChildCmd("com", "command for validation", warg.DoNothing),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "commandNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewChildSection("name", "",
						warg.NewChildCmd("-com", "starts with dash", warg.DoNothing),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "flagNameNoDash",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewChildCmd(
						"c",
						"",
						warg.DoNothing,
						warg.NewChildFlag("f", "", nil),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "aliasNameNoDash",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewChildCmd(
						"c",
						"",
						warg.DoNothing,
						warg.NewChildFlag("-f", "", scalar.Bool(),
							warg.Alias("f"),
						)),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},

		{
			name: "commandFlagAliasCommandFlagNameConflict",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewChildCmd("c", "", warg.DoNothing,
						warg.NewChildFlag("-f", "", scalar.Bool()),
						warg.NewChildFlag("--other", "", scalar.Bool(), warg.Alias("-f")),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagAliasGlobalFlagAliasConflict",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection(
					"help for test",
					warg.NewChildCmd(
						"com",
						"help for com",
						warg.DoNothing,
						warg.NewChildFlag(
							"--commandflag",
							"global flag conflict",
							scalar.String(),
							warg.Alias("--global"),
						),
					),
				),
				warg.SkipValidation(),
				warg.NewGlobalFlag(
					"--globalflag",
					"global flag",
					scalar.String(),
					warg.Alias("--global"),
				),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagAliasGlobalFlagNameConflict",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection(
					"help for test",
					warg.NewChildCmd(
						"com",
						"help for com",
						warg.DoNothing,
						warg.NewChildFlag(
							"--commandflag",
							"global flag conflict",
							scalar.String(),
							warg.Alias("--global"),
						),
					),
				),
				warg.SkipValidation(),
				warg.NewGlobalFlag(
					"--global",
					"global flag",
					scalar.String(),
				),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagNameGlobalFlagNameConflict",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection(
					"help for test",
					warg.NewChildCmd(
						"com",
						"help for com",
						warg.DoNothing,
						warg.NewChildFlag(
							"--global",
							"global flag conflict",
							scalar.String(),
						),
					),
				),
				warg.SkipValidation(),
				warg.NewGlobalFlag(
					"--global",
					"global flag",
					scalar.String(),
				),
			),
			expectedErr: true,
		},
		{
			name: "commandNameSectionNameConflict",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection(
					"help for test",
					warg.NewChildCmd(
						"conflict",
						"help for com",
						warg.DoNothing,
					),
					warg.NewChildSection(
						"conflict",
						"help for section",
					),
				),
				warg.SkipValidation(),
				warg.NewGlobalFlag(
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
