package warg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/section"
)

func TestApp_CompletionCandidates(t *testing.T) {
	app := warg.NewApp(
		"newAppName",
		"v1.0.0",
		section.NewSectionT(
			"root section help",
			section.NewCommand(
				"command1",
				"command1 help",
				command.DoNothing,
			),
		),
	)
	nilLookup := cli.LookupMap(nil)
	// lm := func(key, value string) cli.LookupFunc {
	// 	return cli.LookupMap(map[string]string{key: value})
	// }
	tests := []struct {
		name               string
		args               []string
		lookupFunc         cli.LookupFunc
		expectedErr        bool
		expectedCandidates *completion.Candidates
	}{
		{
			name:        "no args",
			args:        []string{},
			lookupFunc:  nilLookup,
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValueDescription,
				Values: []completion.Candidate{
					{
						Name:        "command1",
						Description: "command1 help",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			// set it up like os.Args
			args := []string{"appName", "--completion-zsh"}
			args = append(args, tt.args...)
			args = append(args, "")

			actualCandidates, actualErr := app.CompletionCandidates(
				cli.OverrideArgs(args),
				cli.OverrideLookupFunc(tt.lookupFunc),
			)

			if tt.expectedErr {
				require.Error(actualErr)
				return
			} else {
				require.NoError(actualErr)
			}
			require.Equal(tt.expectedCandidates, actualCandidates)
		})
	}

}
