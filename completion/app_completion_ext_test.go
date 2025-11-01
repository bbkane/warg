package completion_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/completion/internal/testapp"
)

func TestApp_Completions(t *testing.T) {
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
						Name:        "completion",
						Description: "Print shell completion scripts",
					},
					{
						Name:        "section1",
						Description: "section1 help",
					},
					{
						Name:        "--help",
						Description: "Print help",
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
					{
						Name:        "command3",
						Description: "command with AllowForwardedArgs enabled",
					},
					{
						Name:        "--help",
						Description: "Print help",
					},
				},
			},
		},
		{
			name:        "cmdWithAllowForwardedArgs",
			args:        []string{"section1", "command3"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{Name: "--globalFlag", Description: "globalFlag help"},
					{Name: "--help", Description: "Print help"},
					{Name: "--", Description: "Indicates the end of flag parsing and the beginning of forwarded args"},
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
				Type: completion.Type_Values,
				Values: []completion.Candidate{
					{
						Name:        "alpha",
						Description: "",
					},
					{
						Name:        "beta",
						Description: "",
					},
					{
						Name:        "gamma",
						Description: "",
					},
				},
			},
		},
		{
			name:        "cmdFlagScalarValuePassed",
			args:        []string{"command1", "--flag1", "alpha"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_ValuesDescriptions,
				Values: []completion.Candidate{
					{Name: "--globalFlag", Description: "globalFlag help"},
					{Name: "--help", Description: "Print help"}},
			},
		},
		{
			name:        "cmdFlagBool",
			args:        []string{"section1", "command2", "--bool"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_Values,
				Values: []completion.Candidate{
					{
						Name:        "true",
						Description: "",
					},
					{
						Name:        "false",
						Description: "",
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
		{
			name:        "helpPassedNoValue",
			args:        []string{"--help"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type: completion.Type_Values,
				Values: []completion.Candidate{
					{
						Name:        "allcommands",
						Description: "",
					},
					{
						Name:        "default",
						Description: "",
					},
					{
						Name:        "detailed",
						Description: "",
					},
					{
						Name:        "outline",
						Description: "",
					},
				},
			},
		},
		{
			name:        "helpPassedValue",
			args:        []string{"--help", "default"},
			expectedErr: false,
			expectedCandidates: &completion.Candidates{
				Type:   completion.Type_None,
				Values: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			// set it up like os.Args
			args := []string{"appName", "--completion-zsh"}
			args = append(args, tt.args...)
			// add on the blank space the shell would add for us
			args = append(args, "")

			actualCandidates, actualErr := app.Completions(
				warg.ParseWithArgs(args),
				warg.ParseWithLookupEnv(warg.LookupMap(nil)),
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
