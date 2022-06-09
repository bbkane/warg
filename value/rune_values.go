package value

import (
	"errors"
	"fmt"
)

// There doesn't seem to be an obvious default value for a rune
const emptyRune rune = -1

func runeFromIFace(iFace interface{}) (rune, error) {
	switch under := iFace.(type) {
	case rune:
		return under, nil
	case string:
		return runeFromString(under)
	default:
		return emptyRune, ErrIncompatibleInterface
	}
}

func runeFromString(s string) (rune, error) {
	if s == "" {
		return -1, errors.New("empty string passed")
	}
	var r rune
	if rs := []rune(s); len(rs) != 1 {
		return emptyRune, fmt.Errorf("runes shuld only be one character")
	} else {
		return r, nil
	}
}

// Rune is updateable from a single character string or a rune interface
func Rune() (Value, error) {
	s := newScalarValue(
		0,
		"rune",
		fromIFaceFunc[rune](runeFromIFace),
		fromStringFunc[rune](runeFromString),
	)
	return &s, nil
}
