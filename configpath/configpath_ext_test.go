package configpath_test

import (
	"testing"

	"github.com/bbkane/warg/configpath"
	"github.com/stretchr/testify/assert"
)

func Test_FollowPath(t *testing.T) {
	tests := []struct {
		name                     string
		path                     string
		data                     configpath.ConfigMap
		expectedFollowPathResult configpath.FollowPathResult
		expectedErr              bool
	}{
		{
			name: "one_key",
			path: "key",
			data: configpath.ConfigMap{
				"key": "value",
			},
			expectedFollowPathResult: configpath.FollowPathResult{
				IFace: "value", Exists: true, Aggregated: false,
			},
			expectedErr: false,
		},
		{
			name: "nil_map",
			path: "hi",
			data: nil,
			expectedFollowPathResult: configpath.FollowPathResult{
				IFace: nil, Exists: false, Aggregated: false,
			},
			expectedErr: false,
		},
		{
			name: "two_keys",
			path: "key1.key2",
			data: configpath.ConfigMap{
				"key1": configpath.ConfigMap{
					"key2": 1,
				},
			},
			expectedFollowPathResult: configpath.FollowPathResult{
				IFace: 1, Exists: true, Aggregated: false,
			},
			expectedErr: false,
		},
		// NOTE: this doesn't appear to be equvalent to w.JSONUnmarshaller
		// so let's go with that for the right behavior
		// TODO: turn this commented code into a real explanation
		// {
		// 	// TODO: make this not fail
		// 	name: "in array",
		// 	path: "subreddits[].name",
		// 	data: configpath.ConfigMap{
		// 		"subreddits": []configpath.ConfigMap{ // This should be a list of interfaces
		// 			{
		// 				"name":  "earthporn",
		// 				"limit": 10,
		// 			},
		// 			{
		// 				"name":  "wallpapers",
		// 				"limit": 5,
		// 			},
		// 		},
		// 	},
		// 	expectedFollowPathResult: configpath.FollowPathResult{
		// 		IFace: []interface{}([]interface{}{"earthporn", "wallpapers"}), Exists: true, Aggregated: true,
		// 	},
		// 	expectedErr: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: toggle back between versions
			actualFPR, err := configpath.FollowPath(
				tt.data,
				tt.path,
			)
			// return early if there's an error
			// don't want to deref a null pr
			if (err != nil) && tt.expectedErr {
				return
			}

			if (err != nil) != tt.expectedErr {
				t.Errorf("FollowPath error = %v, expectedErr = %v", err, tt.expectedErr)
				return
			}

			assert.Equal(t, tt.expectedFollowPathResult, actualFPR)
		})
	}
}
