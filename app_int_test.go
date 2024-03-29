package warg

// internal tests - part of the warg package

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGatherArgs2(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		helpFlagNames  []string
		expectedResult *gatherArgsResult
		expectedErr    bool
	}{
		{
			name:          "empty",
			args:          []string{t.Name()},
			helpFlagNames: []string{"-h"},
			expectedResult: &gatherArgsResult{
				Path:       nil,
				FlagStrs:   nil,
				HelpPassed: false,
			},
			expectedErr: false,
		},
		{
			name:           "helpFirst",
			args:           []string{t.Name(), "-h", "default", "other"},
			helpFlagNames:  []string{"--help", "-h"},
			expectedResult: nil,
			expectedErr:    true,
		},
		{
			name:          "helpNoVal",
			args:          []string{t.Name(), "-h"},
			helpFlagNames: []string{"--help", "-h"},
			expectedResult: &gatherArgsResult{
				Path:       nil,
				FlagStrs:   nil,
				HelpPassed: true,
			},
			expectedErr: false,
		},
		{
			name:          "helpWithVal",
			args:          []string{t.Name(), "-h", "val"},
			helpFlagNames: []string{"--help", "-h"},
			expectedResult: &gatherArgsResult{
				Path: nil,
				FlagStrs: []flagStr{
					{NameOrAlias: "-h", Value: "val", Consumed: false}},
				HelpPassed: true,
			},
			expectedErr: false,
		},
		{
			name:           "noFlagVal",
			args:           []string{t.Name(), "cmd", "-f"},
			helpFlagNames:  []string{"--help", "-h"},
			expectedResult: nil,
			expectedErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualResult, actualErr := gatherArgs(tt.args, tt.helpFlagNames)
			if tt.expectedErr {
				require.NotNil(t, actualErr)
			} else {
				require.Nil(t, actualErr)
			}
			require.Equal(t, tt.expectedResult, actualResult)
		})
	}

}
