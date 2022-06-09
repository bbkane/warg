package value

import (
	"fmt"
	"strconv"
)

func intFromIFace(iFace interface{}) (int, error) {
	switch under := iFace.(type) {
	case int:
		return under, nil
	case float64:
		return int(under), nil
	default:
		return 0, ErrIncompatibleInterface
	}
}

func intFromString(s string) (int, error) {
	i, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// Int is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
func Int() (Value, error) {
	s := newScalarValue(
		0,
		"int",
		fromIFaceFunc[int](intFromIFace),
		fromStringFunc[int](intFromString),
	)
	return &s, nil
}

// IntEnum is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
func IntEnum(choices ...int) EmptyConstructor {
	return func() (Value, error) {
		s := newScalarValue(
			0,
			"int enum with choices "+fmt.Sprint(choices),
			fromIFaceEnum(intFromIFace, choices...),
			fromStringEnum(intFromString, choices...),
		)
		return &s, nil
	}
}

// IntSlice is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
func IntSlice() (Value, error) {
	s := newSliceValue(
		nil,
		"int slice",
		fromIFaceFunc[int](intFromIFace),
		fromStringFunc[int](intFromString),
	)
	return &s, nil
}

// IntEnumSlice is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
func IntEnumSlice(choices ...int) (Value, error) {
	s := newSliceValue(
		nil,
		"int enum slice with choices "+fmt.Sprint(choices),
		fromIFaceEnum(intFromIFace, choices...),
		fromStringEnum(intFromString, choices...),
	)
	return &s, nil
}
