package configpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedTokens []Token
		expectedErr    bool
	}{
		{
			name:           "one_key",
			path:           "key",
			expectedTokens: []Token{{Text: "key", Type: TokenTypeKey}},
			expectedErr:    false,
		},
		{
			name: "two_keys",
			path: "key1.key2",
			expectedTokens: []Token{
				{Text: "key1", Type: TokenTypeKey},
				{Text: "key2", Type: TokenTypeKey},
			},
			expectedErr: false,
		},
		{
			name: "key_slice",
			path: "key[].slice_key",
			expectedTokens: []Token{
				{Text: "key", Type: TokenTypeKey},
				{Text: "[]", Type: TokenTypeSlice},
				{Text: "slice_key", Type: TokenTypeKey},
			},
		},
		// TODO: a slice in the middle of a key will break it. Test that once I care enough :)
	}

	for _, tt := range tests {
		gotTokens, gotErr := tokenize(tt.path)
		// return early if there's an error
		// don't want to deref a null pr
		if (gotErr != nil) && tt.expectedErr {
			return
		}

		if (gotErr != nil) != tt.expectedErr {
			t.Errorf("tokenize error = %v, expectedErr = %v", gotErr, tt.expectedErr)
			return
		}
		assert.Equal(t, tt.expectedTokens, gotTokens)
	}
}
