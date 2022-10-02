package value

import (
	"fmt"
	"strconv"
)

type innerTypeInfo[T comparable] struct {
	description string
	fromIFace   func(iFace interface{}) (T, error)
	fromString  func(string) (T, error)
	// Initalized to the empty value, but used for updating stuff in the container type
	empty func() T
}

func Bool() innerTypeInfo[bool] {
	return innerTypeInfo[bool]{
		description: "bool",
		fromIFace: func(iFace interface{}) (bool, error) {
			under, ok := iFace.(bool)
			if !ok {
				return false, ErrIncompatibleInterface
			}
			return under, nil
		},
		fromString: func(s string) (bool, error) {
			switch s {
			case "true":
				return true, nil
			case "false":
				return false, nil
			default:
				return false, fmt.Errorf("expected \"true\" or \"false\", got %s", s)
			}
		},
		empty: func() bool { return false },
	}
}

func Int() innerTypeInfo[int] {
	return innerTypeInfo[int]{
		description: "int",
		fromIFace: func(iFace interface{}) (int, error) {
			switch under := iFace.(type) {
			case int:
				return under, nil
			case float64:
				return int(under), nil
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		fromString: func(s string) (int, error) {
			i, err := strconv.ParseInt(s, 0, strconv.IntSize)
			if err != nil {
				return 0, err
			}
			return int(i), nil
		},
		empty: func() int { return 0 },
	}
}
