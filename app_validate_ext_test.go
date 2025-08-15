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
				warg.NewSection("Help for section", warg.NewSubSection("leaf", "Is empty but shouldn't be")),
				warg.SkipValidation(),
			),
			expectedErr: true,
		},
		// app.Validate should allow app names with dashes
		{
			name: "appNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewSubCmd("com", "command for validation", warg.UnimplementedCmd),
				),
				warg.SkipValidation(),
			),
			expectedErr: false,
		},
		{
			name: "sectionNameWithDash",
			app: warg.New("newAppName", "v1.0.0",
				warg.NewSection("",
					warg.NewSubSection("-name", "",
						warg.NewSubCmd("com", "command for validation", warg.UnimplementedCmd),
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
					warg.NewSubSection("name", "",
						warg.NewSubCmd("-com", "starts with dash", warg.UnimplementedCmd),
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
					warg.NewSubCmd(
						"c",
						"",
						warg.UnimplementedCmd,
						warg.NewCmdFlag("f", "", nil),
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
					warg.NewSubCmd(
						"c",
						"",
						warg.UnimplementedCmd,
						warg.NewCmdFlag("-f", "", scalar.Bool(),
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
					warg.NewSubCmd("c", "", warg.UnimplementedCmd,
						warg.NewCmdFlag("-f", "", scalar.Bool()),
						warg.NewCmdFlag("--other", "", scalar.Bool(), warg.Alias("-f")),
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
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.UnimplementedCmd,
						warg.NewCmdFlag(
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
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.UnimplementedCmd,
						warg.NewCmdFlag(
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
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.UnimplementedCmd,
						warg.NewCmdFlag(
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
					warg.NewSubCmd(
						"conflict",
						"help for com",
						warg.UnimplementedCmd,
					),
					warg.NewSubSection(
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
