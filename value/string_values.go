package value

import "fmt"

func stringFromIFace(iFace interface{}) (string, error) {
	under, ok := iFace.(string)
	if !ok {
		return "", ErrIncompatibleInterface
	}
	return under, nil
}

func stringFromString(s string) (string, error) {
	return s, nil
}

// String accepts a string from a user. Pretty self explanatory.
func String() (Value, error) {
	s := newScalarValue(
		"",
		"string",
		fromIFaceFunc[string](stringFromIFace),
		fromStringFunc[string](stringFromString),
	)
	return &s, nil
}

// StringEnum acts just like a string, except it only lets the user update from
// the choices provided when creating the EmptyConstructor for it.
func StringEnum(choices ...string) EmptyConstructor {
	return func() (Value, error) {
		s := newScalarValue(
			"",
			"string enum with choices "+fmt.Sprint(choices),
			fromIFaceEnum(stringFromIFace, choices...),
			fromStringEnum(stringFromString, choices...),
		)
		return &s, nil
	}
}

// StringSlice accepts a string from a user and adds it to a slice. Pretty self explanatory.
func StringSlice() (Value, error) {
	s := newSliceValue(
		nil,
		"string list",
		fromIFaceFunc[string](stringFromIFace),
		fromStringFunc[string](stringFromString),
	)
	return &s, nil
}

func StringEnumSlice(choices ...string) EmptyConstructor {
	return func() (Value, error) {
		s := newSliceValue(
			nil,
			"string enum slice with choices "+fmt.Sprint(choices),
			fromIFaceEnum(stringFromIFace, choices...),
			fromStringEnum(stringFromString, choices...),
		)
		return &s, nil
	}
}
