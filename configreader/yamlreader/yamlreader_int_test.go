package yamlreader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedTokens []token
		expectedErr    bool
	}{
		{
			name:           "one_key",
			path:           "key",
			expectedTokens: []token{{Text: "key", Type: tokenTypeKey}},
			expectedErr:    false,
		},
		{
			name: "two_keys",
			path: "key1.key2",
			expectedTokens: []token{
				{Text: "key1", Type: tokenTypeKey},
				{Text: "key2", Type: tokenTypeKey},
			},
			expectedErr: false,
		},
		{
			name: "key_slice",
			path: "key[].slice_key",
			expectedTokens: []token{
				{Text: "key", Type: tokenTypeKey},
				{Text: "[]", Type: tokenTypeSlice},
				{Text: "slice_key", Type: tokenTypeKey},
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
