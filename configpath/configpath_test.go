package configpath_test

import (
	"testing"

	"github.com/bbkane/warg/configpath"
	"github.com/stretchr/testify/assert"
)

func Test_FollowPath(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		data           configpath.ConfigMap
		expectedIface  interface{}
		expectedExists bool
		expectedErr    bool
	}{
		{
			name: "from test app",
			path: "key",
			data: configpath.ConfigMap{
				"key": "mapkeyval",
			},
			expectedIface:  "mapkeyval",
			expectedExists: true,
			expectedErr:    false,
		},
		{
			name:           "nil_map",
			path:           "hi",
			data:           nil,
			expectedIface:  nil,
			expectedExists: false,
			expectedErr:    false,
		},
		{
			name: "nested_keys",
			path: "hi.there",
			data: configpath.ConfigMap{
				"hi": configpath.ConfigMap{
					"there": 1,
				},
			},
			expectedIface:  1,
			expectedExists: true,
			expectedErr:    false,
		},
		{
			// TODO: make this not fail
			name: "in array",
			path: "subreddits[].name",
			data: configpath.ConfigMap{
				"subreddits": []configpath.ConfigMap{
					{
						"name":  "earthporn",
						"limit": 10,
					},
					{
						"name":  "wallpapers",
						"limit": 5,
					},
				},
			},
			expectedIface:  []string{"earthporn", "wallpapers"},
			expectedExists: true,
			expectedErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iface, exists, err := configpath.FollowPath(
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

			assert.Equal(t, tt.expectedExists, exists)
			assert.Equal(t, tt.expectedIface, iface)
		})
	}
}
