package configpath

import (
	"fmt"
	"strings"
)

type ConfigMap = map[string]interface{}

// FollowPath takes a map and a path with elements separated by dots
// and retrieves the interface at the end of it. If the interface
// doesn't exist, then the bool value is false
func FollowPath(m ConfigMap, path string) (interface{}, bool, error) {
	pathSlice := strings.Split(path, ".")
	lastIndex := len(pathSlice) - 1
	var err error
	// step down the path
	for _, step := range pathSlice[:lastIndex] {
		nextIface, exists := m[step]
		if !exists {
			return nil, false, nil
		}
		nextMap, isMap := nextIface.(map[string]interface{})
		if !isMap {
			err = fmt.Errorf(
				"error: expected map[string]interface{} at %#v: got %#v",
				step,
				nextIface,
			)
			return nil, false, err
		}
		m = nextMap
	}

	step := pathSlice[lastIndex]
	val, exists := m[step]
	if !exists {
		return nil, false, err
	}

	return val, true, nil
}
