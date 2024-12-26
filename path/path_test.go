package path

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath_expand(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		homedir          string
		expectedExpanded string
		expectedErr      bool
	}{
		{
			name:             "empty",
			path:             "",
			homedir:          "",
			expectedExpanded: "",
			expectedErr:      false,
		},
		{
			name:             "nonTildePrefix",
			path:             "bob",
			homedir:          "",
			expectedExpanded: "bob",
			expectedErr:      false,
		},
		{
			name:             "secondCharNotCharacter",
			path:             "~BAD",
			homedir:          "",
			expectedExpanded: "",
			expectedErr:      true,
		},
		{
			name:             "expanded",
			path:             "~/name",
			homedir:          "/homedir",
			expectedExpanded: "/homedir/name",
			expectedErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.path)
			actualExpanded, actualErr := p.expand(tt.homedir)
			if tt.expectedErr {
				require.Error(t, actualErr)
			} else {
				require.NoError(t, actualErr)
			}
			require.Equal(t, tt.expectedExpanded, actualExpanded)
		})
	}
}
