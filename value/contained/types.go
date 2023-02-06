package contained

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/xhit/go-str2duration/v2"
)

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")

// identity simply returns the thing passed and nil
func identity[T comparable](t T) (T, error) {
	return t, nil
}

type ContainedTypeInfo[T comparable] struct {
	Description string

	FromIFace func(iFace interface{}) (T, error)

	// FromInstance updates a T from an instance of itself.
	// This is particularly usefule for paths - when the user sets a scalar.Default of `~`,
	// we want to expand that into /path/to/home the same way we would
	// when updating from a string in the CLI
	FromInstance func(T) (T, error)

	FromString func(string) (T, error)

	// Initalized to the Empty value, but used for updating stuff in the container type
	Empty func() T
}

func Bool() ContainedTypeInfo[bool] {
	return ContainedTypeInfo[bool]{
		Description: "bool",
		Empty:       func() bool { return false },
		FromIFace: func(iFace interface{}) (bool, error) {
			under, ok := iFace.(bool)
			if !ok {
				return false, ErrIncompatibleInterface
			}
			return under, nil
		},
		FromInstance: identity[bool],
		FromString: func(s string) (bool, error) {
			switch s {
			case "true":
				return true, nil
			case "false":
				return false, nil
			default:
				return false, fmt.Errorf("expected \"true\" or \"false\", got %s", s)
			}
		},
	}
}

func durationFromString(s string) (time.Duration, error) {
	decoded, err := str2duration.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return decoded, nil
}

func Duration() ContainedTypeInfo[time.Duration] {
	return ContainedTypeInfo[time.Duration]{
		Description: "duration",
		Empty: func() time.Duration {
			var t time.Duration = 0
			return t
		},
		FromIFace: func(iFace interface{}) (time.Duration, error) {
			under, ok := iFace.(string)
			if !ok {
				return 0, ErrIncompatibleInterface
			}
			return durationFromString(under)
		},
		FromInstance: identity[time.Duration],
		FromString:   durationFromString,
	}
}

func intFromString(s string) (int, error) {
	i, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func Int() ContainedTypeInfo[int] {
	return ContainedTypeInfo[int]{
		Description: "int",
		FromIFace: func(iFace interface{}) (int, error) {
			switch under := iFace.(type) {
			case int:
				return under, nil
			case json.Number:
				return intFromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromInstance: identity[int],
		FromString:   intFromString,
		Empty:        func() int { return 0 },
	}
}

func pathFromString(s string) (string, error) {
	expanded, err := homedir.Expand(s)
	if err != nil {
		return "", err
	}
	return expanded, nil
}

func Path() ContainedTypeInfo[string] {
	return ContainedTypeInfo[string]{
		Description: "path",
		Empty:       func() string { return "" },
		FromIFace: func(iFace interface{}) (string, error) {
			under, ok := iFace.(string)
			if !ok {
				return "", ErrIncompatibleInterface
			}
			return pathFromString(under)
		},
		FromInstance: pathFromString,
		FromString:   pathFromString,
	}
}

// There doesn't seem to be an obvious default value for a rune
const emptyRune rune = -1

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

func Rune() ContainedTypeInfo[rune] {
	return ContainedTypeInfo[rune]{
		Description: "rune",
		Empty:       func() rune { return emptyRune },
		FromIFace: func(iFace interface{}) (rune, error) {
			switch under := iFace.(type) {
			case rune:
				return under, nil
			case string:
				return runeFromString(under)
			default:
				return emptyRune, ErrIncompatibleInterface
			}
		},
		FromInstance: identity[rune],
		FromString:   runeFromString,
	}
}

func String() ContainedTypeInfo[string] {
	return ContainedTypeInfo[string]{
		Description: "string",
		Empty:       func() string { return "" },
		FromIFace: func(iFace interface{}) (string, error) {
			under, ok := iFace.(string)
			if !ok {
				return "", ErrIncompatibleInterface
			}
			return under, nil
		},
		FromInstance: identity[string],
		FromString: func(s string) (string, error) {
			return s, nil
		},
	}
}
