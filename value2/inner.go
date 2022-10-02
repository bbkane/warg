package value

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/xhit/go-str2duration/v2"
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
		empty:       func() bool { return false },
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
	}
}

func durationFromString(s string) (time.Duration, error) {
	decoded, err := str2duration.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return decoded, nil
}

func Duration() innerTypeInfo[time.Duration] {
	return innerTypeInfo[time.Duration]{
		description: "duration",
		empty: func() time.Duration {
			var t time.Duration = 0
			return t
		},
		fromIFace: func(iFace interface{}) (time.Duration, error) {
			under, ok := iFace.(string)
			if !ok {
				return 0, ErrIncompatibleInterface
			}
			return durationFromString(under)
		},
		fromString: durationFromString,
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

func pathFromString(s string) (string, error) {
	expanded, err := homedir.Expand(s)
	if err != nil {
		return "", err
	}
	return expanded, nil
}

func Path() innerTypeInfo[string] {
	return innerTypeInfo[string]{
		description: "path",
		empty:       func() string { return "" },
		fromIFace: func(iFace interface{}) (string, error) {
			under, ok := iFace.(string)
			if !ok {
				return "", ErrIncompatibleInterface
			}
			return pathFromString(under)
		},
		fromString: pathFromString,
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

func Rune() innerTypeInfo[rune] {
	return innerTypeInfo[rune]{
		description: "rune",
		empty:       func() rune { return emptyRune },
		fromIFace: func(iFace interface{}) (rune, error) {
			switch under := iFace.(type) {
			case rune:
				return under, nil
			case string:
				return runeFromString(under)
			default:
				return emptyRune, ErrIncompatibleInterface
			}
		},
		fromString: runeFromString,
	}
}

func String() innerTypeInfo[string] {
	return innerTypeInfo[string]{
		description: "string",
		empty:       func() string { return "" },
		fromIFace: func(iFace interface{}) (string, error) {
			under, ok := iFace.(string)
			if !ok {
				return "", ErrIncompatibleInterface
			}
			return under, nil
		},
		fromString: func(s string) (string, error) {
			return s, nil
		},
	}
}
