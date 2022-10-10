package types

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/xhit/go-str2duration/v2"
)

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")

type ContainedTypeInfo[T comparable] struct {
	Description string
	FromIFace   func(iFace interface{}) (T, error)
	FromString  func(string) (T, error)
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
		FromString: durationFromString,
	}
}

func Int() ContainedTypeInfo[int] {
	return ContainedTypeInfo[int]{
		Description: "int",
		FromIFace: func(iFace interface{}) (int, error) {
			switch under := iFace.(type) {
			case int:
				return under, nil
			case float64:
				return int(under), nil
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: func(s string) (int, error) {
			i, err := strconv.ParseInt(s, 0, strconv.IntSize)
			if err != nil {
				return 0, err
			}
			return int(i), nil
		},
		Empty: func() int { return 0 },
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
		FromString: pathFromString,
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
		FromString: runeFromString,
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
		FromString: func(s string) (string, error) {
			return s, nil
		},
	}
}
