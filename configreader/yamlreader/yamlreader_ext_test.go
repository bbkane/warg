package yamlreader_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bbkane/warg/configreader"
	"github.com/bbkane/warg/configreader/yamlreader"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name                 string
		filePath             string
		searchPath           string
		expectedCreationErr  bool
		expectedSearchResult configreader.ConfigSearchResult
		expectedSearchErr    bool
	}{
		{
			name:                "one_key",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "key",
			expectedCreationErr: false,
			expectedSearchResult: configreader.ConfigSearchResult{
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
			expectedSearchResult: configreader.ConfigSearchResult{
				IFace:        nil,
				Exists:       false,
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "two_keys",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "key1.key2",
			expectedCreationErr: false,
			expectedSearchResult: configreader.ConfigSearchResult{
				IFace:        1,
				Exists:       true,
				IsAggregated: false,
			},
			expectedSearchErr: false,
		},
		{
			name:                "in_array",
			filePath:            "testdata/TestSearch.yaml",
			searchPath:          "subreddits[].name",
			expectedCreationErr: false,
			expectedSearchResult: configreader.ConfigSearchResult{
				IFace:        []interface{}([]interface{}{"earthporn", "wallpapers"}),
				Exists:       true,
				IsAggregated: true,
			},
			expectedSearchErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr, err := yamlreader.NewYAMLConfigReader(tt.filePath)

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
