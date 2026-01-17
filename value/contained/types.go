package contained

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/netip"
	"strconv"
	"time"

	"github.com/xhit/go-str2duration/v2"
	"go.bbkane.com/warg/path"
)

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")

// FromZero returns the zero value for type T. Useful for contstructing [TypeInfo] instances.
func FromZero[T any]() T {
	var zero T
	return zero
}

// Equals returns true if a and b are equal. CUseful for contstructing [TypeInfo] instances.
func Equals[T comparable](a, b T) bool {
	return a == b
}

type TypeInfo[T any] struct {
	Description string

	FromIFace func(iFace interface{}) (T, error)

	FromString func(string) (T, error)

	// FromZero returns an initial value for type T. This is used as the intial value for contained types and updated from there. Most types will want to use the [Zero] helper function here.
	FromZero func() T

	// Equals returns true if a and b are equal. Comparable types will want to use the [Equals] helper function here.
	Equals func(a, b T) bool
}

// ValidateNonNilFuncs returns an error if any of the function fields are nil. Used to validate TypeInfo instances in tests
func (ti TypeInfo[T]) ValidateNonNilFuncs() error {
	var errs []error
	if ti.FromIFace == nil {
		errs = append(errs, fmt.Errorf("FromIFace is nil"))
	}
	if ti.FromString == nil {
		errs = append(errs, fmt.Errorf("FromString is nil"))
	}
	if ti.FromZero == nil {
		errs = append(errs, fmt.Errorf("FromZero is nil"))
	}
	if ti.Equals == nil {
		errs = append(errs, fmt.Errorf("Equals is nil"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("nil fields: %w", errors.Join(errs...))
	}

	return nil
}

// WithinChoices returns true if val is within choices according to equals function. Used to update values when passed as strings from flags
func WithinChoices[T any](val T, choices []T, equals func(a, b T) bool) bool {
	// User didn't constrain choices
	if len(choices) == 0 {
		return true
	}
	for _, choice := range choices {
		if equals(val, choice) {
			return true
		}
	}
	return false
}

func NetIPAddr() TypeInfo[netip.Addr] {
	return TypeInfo[netip.Addr]{
		Description: "IP address",
		FromZero:    FromZero[netip.Addr],
		FromIFace: func(iFace interface{}) (netip.Addr, error) {
			switch under := iFace.(type) {
			case netip.Addr:
				return under, nil
			case []byte:
				ip, ok := netip.AddrFromSlice(under)
				if !ok {
					return netip.Addr{}, fmt.Errorf("could not convert %s to netip.Addr", string(under))
				}
				return ip, nil
			case string:
				return netip.ParseAddr(under)
			default:
				return netip.Addr{}, ErrIncompatibleInterface
			}
		},

		FromString: netip.ParseAddr,
		Equals:     Equals[netip.Addr],
	}
}

func AddrPort() TypeInfo[netip.AddrPort] {
	return TypeInfo[netip.AddrPort]{
		Description: "IP and Port number separated by a colon: ip:port ",
		FromZero:    FromZero[netip.AddrPort],
		FromIFace: func(iFace interface{}) (netip.AddrPort, error) {
			switch under := iFace.(type) {
			case netip.AddrPort:
				return under, nil
			case string:
				return netip.ParseAddrPort(under)
			default:
				return netip.AddrPort{}, ErrIncompatibleInterface
			}
		},
		FromString: netip.ParseAddrPort,
		Equals:     Equals[netip.AddrPort],
	}
}

func Bool() TypeInfo[bool] {
	return TypeInfo[bool]{
		Description: "bool",
		FromZero:    FromZero[bool],
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
		Equals: Equals[bool],
	}
}

func durationFromString(s string) (time.Duration, error) {
	decoded, err := str2duration.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return decoded, nil
}

func Duration() TypeInfo[time.Duration] {
	return TypeInfo[time.Duration]{
		Description: "duration",
		FromZero:    FromZero[time.Duration],
		FromIFace: func(iFace interface{}) (time.Duration, error) {
			under, ok := iFace.(string)
			if !ok {
				return 0, ErrIncompatibleInterface
			}
			return durationFromString(under)
		},
		FromString: durationFromString,
		Equals:     Equals[time.Duration],
	}
}

func intFromString(s string) (int, error) {
	i, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func Int() TypeInfo[int] {
	return TypeInfo[int]{
		Description: "int",
		FromIFace: func(iFace interface{}) (int, error) {
			switch under := iFace.(type) {
			case int:
				return under, nil
			// go-yaml may decode all numbers as int64 or uint64
			case int64:
				if under > math.MaxInt || under < math.MinInt {
					return 0, fmt.Errorf("int64 value %d out of range for int", under)
				}
				return int(under), nil
			case uint64:
				if under > math.MaxInt {
					return 0, fmt.Errorf("uint64 value %d out of range for int", under)
				}
				return int(under), nil
			case json.Number:
				return intFromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: intFromString,
		FromZero:   FromZero[int],
		Equals:     Equals[int],
	}
}

func Path() TypeInfo[path.Path] {
	return TypeInfo[path.Path]{
		Description: "path",
		FromZero:    func() path.Path { return path.New("") },
		FromIFace: func(iFace interface{}) (path.Path, error) {
			under, ok := iFace.(string)
			if !ok {
				return path.New(""), ErrIncompatibleInterface
			}
			return path.New(under), nil
		},
		FromString: func(s string) (path.Path, error) { return path.New(s), nil },
		Equals:     func(a, b path.Path) bool { return a.Equals(b) },
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

func Rune() TypeInfo[rune] {
	return TypeInfo[rune]{
		Description: "rune",
		FromZero:    func() rune { return emptyRune },
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
		Equals:     Equals[rune],
	}
}

func String() TypeInfo[string] {
	return TypeInfo[string]{
		Description: "string",
		FromZero:    FromZero[string],
		FromIFace: func(iFace interface{}) (string, error) {
			under, ok := iFace.(string)
			if !ok {
				return "", ErrIncompatibleInterface
			}
			return under, nil
		},
		FromString: func(s string) (string, error) { return s, nil },
		Equals:     Equals[string],
	}
}
