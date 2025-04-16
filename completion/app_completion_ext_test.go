package completion_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/completion/internal/testapp"
	"go.bbkane.com/warg/parseopt"
)

func TestApp_CompletionCandidates(t *testing.T) {
	// To try to make this more concise, these tests are gonna share an app...
	app := testapp.BuildApp()

	globalFlagcompletion := completion.Candidate{
		Name:        "--globalFlag",
		Description: "globalFlag help",
	}
	helpCompletion := completion.Candidate{
		Name:        "--help",
		Description: "Print help",
	}
	tests := []struct {
		name               string
		args               []string
		expectedErr        bool
		expectedCandidates *completion.Candidates
	}{
		{
			name:        "noArgs",
			args:        []string{},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{
						Name:        "command1",
						Description: "command1 help",
					},
					{
						Name:        "manual",
						Description: "commands with flags using all completion types for manual testing",
					},
					{
						Name:        "section1",
						Description: "section1 help",
					},
				},
			},
		},
		{
			name:               "moreArgsThanSections",
			args:               []string{"bob"},
			expectedErr:        true,
			expectedCandidates: nil,
		},
		{
			name:        "childSectionCommands",
			args:        []string{"section1"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{
						Name:        "command2",
						Description: "command2 help",
					},
				},
			},
		},
		{
			name:        "cmdFlagName",
			args:        []string{"command1"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{
						Name:        "--flag1",
						Description: "flag1 help",
					},
					globalFlagcompletion,
					helpCompletion,
				},
			},
		},
		{
			name:        "cmdFlagValue",
			args:        []string{"command1", "--flag1"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{
						Name:        "alpha",
						Description: "NO DESCRIPTION",
					},
					{
						Name:        "beta",
						Description: "NO DESCRIPTION",
					},
					{
						Name:        "gamma",
						Description: "NO DESCRIPTION",
					},
				},
			},
		},
		{
			name:        "cmdFlagCustomCompletionDefault",
			args:        []string{"section1", "command2", "--flag2"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{
						Name:        "default",
						Description: "default completion",
					},
				},
			},
		},
		{
			name:        "cmdFlagCustomCompletionNonDefault",
			args:        []string{"section1", "command2", "--globalFlag", "nondefault", "--flag2"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{
						Name:        "nondefault",
						Description: "nondefault completion",
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
				parseopt.Args(args),
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
