package value

import (
	"github.com/mitchellh/go-homedir"
)

func pathFromIFace(iFace interface{}) (string, error) {
	under, ok := iFace.(string)
	if !ok {
		return "", ErrIncompatibleInterface
	}
	return pathFromString(under)
}

func pathFromString(s string) (string, error) {
	expanded, err := homedir.Expand(s)
	if err != nil {
		return "", err
	}
	return expanded, nil
}

// Path autoexpands ~ when updated and otherwise behaves like a string
func Path() (Value, error) {
	s := newScalarValue(
		"",
		"path",
		fromIFaceFunc[string](pathFromIFace),
		fromStringFunc[string](pathFromString),
	)
	return &s, nil
}

// PathSlice autoexpands ~ when updated and otherwise behaves like a []string
func PathSlice() (Value, error) {
	s := newSliceValue(
		nil,
		"path",
		fromIFaceFunc[string](pathFromIFace),
		fromStringFunc[string](pathFromString),
	)
	return &s, nil
}
