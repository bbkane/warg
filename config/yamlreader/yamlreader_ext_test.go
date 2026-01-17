package yamlreader_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/config/yamlreader"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name                 string
		filePath             string
		searchPath           string
		expectedCreationErr  bool
		expectedSearchResult *config.SearchResult
		expectedSearchErr    bool
	}{
		{
			name:                "one_key",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "key",
			expectedCreationErr: false,
			expectedSearchResult: &config.SearchResult{
				IFace:        "value",
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "uint64_key",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "uint64_key",
			expectedCreationErr: false,
			expectedSearchResult: &config.SearchResult{
				IFace:        uint64(42),
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "int64_key_negative",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "int64_key_negative",
			expectedCreationErr: false,
			expectedSearchResult: &config.SearchResult{
				IFace:        int64(-42),
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                 "nil_map",
			filePath:             "non-existant",
			searchPath:           "non-existant",
			expectedCreationErr:  false, // It's ok to not have a config file
			expectedSearchResult: nil,
			expectedSearchErr:    false,
		},
		{
			name:                "two_keys",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "key1.key2",
			expectedCreationErr: false,
			expectedSearchResult: &config.SearchResult{
				IFace:        uint64(1),
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "in_array",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "subreddits[].name",
			expectedCreationErr: false,
			expectedSearchResult: &config.SearchResult{
				IFace:        []interface{}([]interface{}{"earthporn", "wallpapers"}),
				IsAggregated: true,
			},
			expectedSearchErr: false,
		},
		{
			name:                "map_val",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "map_val",
			expectedCreationErr: false,
			expectedSearchResult: &config.SearchResult{
				IFace:        map[string]interface{}{"a": uint64(1)},
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr, err := yamlreader.New(tt.filePath)

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
