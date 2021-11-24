package jsonreader_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bbkane/warg/config"
	"github.com/bbkane/warg/config/jsonreader"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name                 string
		filePath             string
		searchPath           string
		expectedCreationErr  bool
		expectedSearchResult config.SearchResult
		expectedSearchErr    bool
	}{
		{
			name:                "one_key",
			filePath:            "testdata/TestSearch.json",
			searchPath:          "key",
			expectedCreationErr: false,
			expectedSearchResult: config.SearchResult{
				IFace:        "value",
				Exists:       true,
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "nil_map",
			filePath:            "non-existant",
			searchPath:          "non-existant",
			expectedCreationErr: false, // It's ok to not have a config file
			expectedSearchResult: config.SearchResult{
				IFace:        nil,
				Exists:       false,
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "two_keys",
			filePath:            "testdata/TestSearch.json",
			searchPath:          "key1.key2",
			expectedCreationErr: false,
			expectedSearchResult: config.SearchResult{
				IFace:        float64(1),
				Exists:       true,
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "in_array",
			filePath:            "testdata/TestSearch.json",
			searchPath:          "subreddits[].name",
			expectedCreationErr: false,
			expectedSearchResult: config.SearchResult{
				IFace:        []interface{}([]interface{}{"earthporn", "wallpapers"}),
				Exists:       true,
				IsAggregated: true,
			},
			expectedSearchErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr, err := jsonreader.New(tt.filePath)

			if tt.expectedCreationErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

			res, err := cr.Search(tt.searchPath)

			if tt.expectedSearchErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tt.expectedSearchResult, res)
		})
	}
}
