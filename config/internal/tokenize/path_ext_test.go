package tokenize_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.bbkane.com/warg/config/internal/tokenize"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedTokens []tokenize.Token
		expectedErr    bool
	}{
		{
			name:           "one_key",
			path:           "key",
			expectedTokens: []tokenize.Token{{Text: "key", Type: tokenize.TokenTypeKey}},
			expectedErr:    false,
		},
		{
			name: "two_keys",
			path: "key1.key2",
			expectedTokens: []tokenize.Token{
				{Text: "key1", Type: tokenize.TokenTypeKey},
				{Text: "key2", Type: tokenize.TokenTypeKey},
			},
			expectedErr: false,
		},
		{
			name: "key_slice",
			path: "key[].slice_key",
			expectedTokens: []tokenize.Token{
				{Text: "key", Type: tokenize.TokenTypeKey},
				{Text: "[]", Type: tokenize.TokenTypeSlice},
				{Text: "slice_key", Type: tokenize.TokenTypeKey},
			},
			expectedErr: false,
		},
		// TODO: a slice in the middle of a key will break it. Test that once I care enough :)
	}

	for _, tt := range tests {
		gotTokens, gotErr := tokenize.Tokenize(tt.path)
		// return early if there's an error
		// don't want to deref a null pr
		if (gotErr != nil) && tt.expectedErr {
			return
		}

		if (gotErr != nil) != tt.expectedErr {
			t.Errorf("tokenize error = %v, expectedErr = %v", gotErr, tt.expectedErr)
			return
		}
		require.Equal(t, tt.expectedTokens, gotTokens)
	}
}
