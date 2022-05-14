package value

import "fmt"

func boolFromIFace(iFace interface{}) (bool, error) {
	under, ok := iFace.(bool)
	if !ok {
		return false, ErrIncompatibleInterface
	}
	return under, nil
}

func boolFromString(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("expected \"true\" or \"false\", got %s", s)
	}
}

// Bool is updated from "true" or "false"
func Bool() (Value, error) {
	s := newScalarValue(
		false,
		"bool",
		fromIFaceFunc[bool](boolFromIFace),
		fromStringFunc[bool](boolFromString),
	)
	return &s, nil
}
