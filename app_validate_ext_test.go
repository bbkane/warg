package warg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func TestApp_Validate(t *testing.T) {

	tests := []struct {
		name        string
		app         wargcore.App // NOTE:
		expectedErr bool
	}{
		{
			name: "leafSection",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection("Help for section", wargcore.NewChildSection("leaf", "Is empty but shouldn't be")),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		// app.Validate should allow app names with dashes
		{
			name: "appNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection("",
					wargcore.NewChildCmd("com", "command for validation", wargcore.DoNothing),
				),
				warg.SkipValidation(),
			),
			expectedErr: false,
		},
		{
			name: "sectionNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection("",
					wargcore.NewChildSection("-name", "",
						wargcore.NewChildCmd("com", "command for validation", wargcore.DoNothing),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "commandNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection("",
					wargcore.NewChildSection("name", "",
						wargcore.NewChildCmd("-com", "starts with dash", wargcore.DoNothing),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "flagNameNoDash",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection("",
					wargcore.NewChildCmd(
						"c",
						"",
						wargcore.DoNothing,
						wargcore.NewChildFlag("f", "", nil),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "aliasNameNoDash",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection("",
					wargcore.NewChildCmd(
						"c",
						"",
						wargcore.DoNothing,
						wargcore.NewChildFlag("-f", "", scalar.Bool(),
							wargcore.Alias("f"),
						)),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},

		{
			name: "commandFlagAliasCommandFlagNameConflict",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection("",
					wargcore.NewChildCmd("c", "", wargcore.DoNothing,
						wargcore.NewChildFlag("-f", "", scalar.Bool()),
						wargcore.NewChildFlag("--other", "", scalar.Bool(), wargcore.Alias("-f")),
					),
				),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagAliasGlobalFlagAliasConflict",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection(
					"help for test",
					wargcore.NewChildCmd(
						"com",
						"help for com",
						wargcore.DoNothing,
						wargcore.NewChildFlag(
							"--commandflag",
							"global flag conflict",
							scalar.String(),
							wargcore.Alias("--global"),
						),
					),
				),
				warg.SkipValidation(),
				warg.NewGlobalFlag(
					"--globalflag",
					"global flag",
					scalar.String(),
					wargcore.Alias("--global"),
				),
			),
			expectedErr: true,
		},
		{
			name: "commandFlagAliasGlobalFlagNameConflict",
			app: warg.New("newAppName", "v1.0.0",
				wargcore.NewSection(
					"help for test",
					wargcore.NewChildCmd(
						"com",
						"help for com",
						wargcore.DoNothing,
						wargcore.NewChildFlag(
							"--commandflag",
							"global flag conflict",
							scalar.String(),
							wargcore.Alias("--global"),
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
				wargcore.NewSection(
					"help for test",
					wargcore.NewChildCmd(
						"com",
						"help for com",
						wargcore.DoNothing,
						wargcore.NewChildFlag(
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
				wargcore.NewSection(
					"help for test",
					wargcore.NewChildCmd(
						"conflict",
						"help for com",
						wargcore.DoNothing,
					),
					wargcore.NewChildSection(
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
