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
	"go.bbkane.com/warg/colerr"
	"go.bbkane.com/warg/path"
)

// ErrIncompatibleInterface is returned when a value cannot be decoded from an interface{}
// (typically from a config file whose type doesn't match the expected Go type).
var ErrIncompatibleInterface = errors.New("Could not decode interface into Value")

// FromZero returns the zero value for type T. Use as the FromZero field in [TypeInfo].
func FromZero[T any]() T {
	var zero T
	return zero
}

// Equals returns true if a and b are equal. Use as the Equals field in [TypeInfo]
// for comparable types.
func Equals[T comparable](a, b T) bool {
	return a == b
}

// TypeInfo defines how to parse, compare, and initialize values of type T.
// Provide one to [scalar.New], [slice.New], or [dict.New] for custom types.
type TypeInfo[T any] struct {
	Description string

	FromIFace func(iFace interface{}) (T, error)

	FromString func(string) (T, error)

	// FromZero returns an initial value for type T. This is used as the intial value for contained types and updated from there. Most types will want to use the [Zero] helper function here.
	FromZero func() T

	// Equals returns true if a and b are equal. Comparable types will want to use the [Equals] helper function here.
	Equals func(a, b T) bool
}

// ValidateNonNilFuncs returns an error if any function fields are nil.
// Useful in tests to catch incomplete [TypeInfo] definitions.
func (ti TypeInfo[T]) ValidateNonNilFuncs() error {
	var errs []error
	if ti.FromIFace == nil {
		errs = append(errs, errors.New("FromIFace is nil"))
	}
	if ti.FromString == nil {
		errs = append(errs, errors.New("FromString is nil"))
	}
	if ti.FromZero == nil {
		errs = append(errs, errors.New("FromZero is nil"))
	}
	if ti.Equals == nil {
		errs = append(errs, errors.New("Equals is nil"))
	}

	if len(errs) > 0 {
		return colerr.NewWrapped(errors.Join(errs...), "Nil fields")
	}

	return nil
}

// WithinChoices reports whether val is in choices using the provided equals function.
// Returns true if choices is empty (no constraint).
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

// NetIPAddr returns a [TypeInfo] for [netip.Addr] values.
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
					return netip.Addr{}, colerr.NewWrappedf(nil, "Could not convert %s to netip.Addr", string(under))
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

// AddrPort returns a [TypeInfo] for [netip.AddrPort] values (ip:port format).
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

// Bool returns a [TypeInfo] for boolean values. Accepts "true" or "false" strings.
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
				return false, colerr.NewWrappedf(nil, "Expected \"true\" or \"false\", got %s", s)
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

// Duration returns a [TypeInfo] for [time.Duration] values.
// Accepts Go duration strings (e.g., "1h30m") as well as extended formats like "1d2h".
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

// DateTimeRFC3339 returns a [TypeInfo] for [time.Time] values in RFC3339 format.
func DateTimeRFC3339() TypeInfo[time.Time] {
	return TypeInfo[time.Time]{
		Description: "datetime in RFC3339 format",
		FromZero:    FromZero[time.Time],
		FromIFace: func(iFace interface{}) (time.Time, error) {
			under, ok := iFace.(string)
			if !ok {
				return time.Time{}, ErrIncompatibleInterface
			}
			return time.Parse(time.RFC3339, under)
		},
		FromString: func(s string) (time.Time, error) {
			return time.Parse(time.RFC3339, s)
		},
		Equals: func(a, b time.Time) bool {
			return a.Equal(b)
		},
	}
}

func intFromString(s string) (int, error) {
	i, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// Int returns a [TypeInfo] for int values. Accepts decimal, hex (0x), octal (0o), and binary (0b) strings.
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
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for int", fmt.Sprintf("%d", under))
				}
				return int(under), nil
			case uint64:
				if under > math.MaxInt {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for int", fmt.Sprintf("%d", under))
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

func int8FromString(s string) (int8, error) {
	i, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

// Int8 returns a [TypeInfo] for int8 values (range -128 to 127).
func Int8() TypeInfo[int8] {
	return TypeInfo[int8]{
		Description: "int8",
		FromIFace: func(iFace interface{}) (int8, error) {
			switch under := iFace.(type) {
			case int8:
				return under, nil
			case int:
				if under > math.MaxInt8 || under < math.MinInt8 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for int8", fmt.Sprintf("%d", under))
				}
				return int8(under), nil
			case int64:
				if under > math.MaxInt8 || under < math.MinInt8 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for int8", fmt.Sprintf("%d", under))
				}
				return int8(under), nil
			case uint64:
				if under > math.MaxInt8 {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for int8", fmt.Sprintf("%d", under))
				}
				return int8(under), nil
			case json.Number:
				return int8FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: int8FromString,
		FromZero:   FromZero[int8],
		Equals:     Equals[int8],
	}
}

func int16FromString(s string) (int16, error) {
	i, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

// Int16 returns a [TypeInfo] for int16 values (range -32768 to 32767).
func Int16() TypeInfo[int16] {
	return TypeInfo[int16]{
		Description: "int16",
		FromIFace: func(iFace interface{}) (int16, error) {
			switch under := iFace.(type) {
			case int16:
				return under, nil
			case int:
				if under > math.MaxInt16 || under < math.MinInt16 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for int16", fmt.Sprintf("%d", under))
				}
				return int16(under), nil
			case int64:
				if under > math.MaxInt16 || under < math.MinInt16 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for int16", fmt.Sprintf("%d", under))
				}
				return int16(under), nil
			case uint64:
				if under > math.MaxInt16 {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for int16", fmt.Sprintf("%d", under))
				}
				return int16(under), nil
			case json.Number:
				return int16FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: int16FromString,
		FromZero:   FromZero[int16],
		Equals:     Equals[int16],
	}
}

func int32FromString(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// Int32 returns a [TypeInfo] for int32 values (range -2147483648 to 2147483647).
func Int32() TypeInfo[int32] {
	return TypeInfo[int32]{
		Description: "int32",
		FromIFace: func(iFace interface{}) (int32, error) {
			switch under := iFace.(type) {
			case int32:
				return under, nil
			case int:
				if under > math.MaxInt32 || under < math.MinInt32 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for int32", fmt.Sprintf("%d", under))
				}
				return int32(under), nil
			case int64:
				if under > math.MaxInt32 || under < math.MinInt32 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for int32", fmt.Sprintf("%d", under))
				}
				return int32(under), nil
			case uint64:
				if under > math.MaxInt32 {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for int32", fmt.Sprintf("%d", under))
				}
				return int32(under), nil
			case json.Number:
				return int32FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: int32FromString,
		FromZero:   FromZero[int32],
		Equals:     Equals[int32],
	}
}

func int64FromString(s string) (int64, error) {
	return strconv.ParseInt(s, 0, 64)
}

// Int64 returns a [TypeInfo] for int64 values.
func Int64() TypeInfo[int64] {
	return TypeInfo[int64]{
		Description: "int64",
		FromIFace: func(iFace interface{}) (int64, error) {
			switch under := iFace.(type) {
			case int64:
				return under, nil
			case int:
				return int64(under), nil
			case uint64:
				if under > math.MaxInt64 {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for int64", fmt.Sprintf("%d", under))
				}
				return int64(under), nil
			case json.Number:
				return int64FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: int64FromString,
		FromZero:   FromZero[int64],
		Equals:     Equals[int64],
	}
}

func uintFromString(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

// Uint returns a [TypeInfo] for uint values.
func Uint() TypeInfo[uint] {
	return TypeInfo[uint]{
		Description: "uint",
		FromIFace: func(iFace interface{}) (uint, error) {
			switch under := iFace.(type) {
			case uint:
				return under, nil
			case uint64:
				if under > math.MaxUint {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for uint", fmt.Sprintf("%d", under))
				}
				return uint(under), nil
			case int:
				if under < 0 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for uint", fmt.Sprintf("%d", under))
				}
				return uint(under), nil
			case int64:
				if under < 0 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for uint", fmt.Sprintf("%d", under))
				}
				return uint(under), nil
			case json.Number:
				return uintFromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: uintFromString,
		FromZero:   FromZero[uint],
		Equals:     Equals[uint],
	}
}

func uint8FromString(s string) (uint8, error) {
	i, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

// Uint8 returns a [TypeInfo] for uint8 values (range 0 to 255).
func Uint8() TypeInfo[uint8] {
	return TypeInfo[uint8]{
		Description: "uint8",
		FromIFace: func(iFace interface{}) (uint8, error) {
			switch under := iFace.(type) {
			case uint8:
				return under, nil
			case uint64:
				if under > math.MaxUint8 {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for uint8", fmt.Sprintf("%d", under))
				}
				return uint8(under), nil
			case int:
				if under < 0 || under > math.MaxUint8 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for uint8", fmt.Sprintf("%d", under))
				}
				return uint8(under), nil
			case int64:
				if under < 0 || under > math.MaxUint8 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for uint8", fmt.Sprintf("%d", under))
				}
				return uint8(under), nil
			case json.Number:
				return uint8FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: uint8FromString,
		FromZero:   FromZero[uint8],
		Equals:     Equals[uint8],
	}
}

func uint16FromString(s string) (uint16, error) {
	i, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

// Uint16 returns a [TypeInfo] for uint16 values (range 0 to 65535).
func Uint16() TypeInfo[uint16] {
	return TypeInfo[uint16]{
		Description: "uint16",
		FromIFace: func(iFace interface{}) (uint16, error) {
			switch under := iFace.(type) {
			case uint16:
				return under, nil
			case uint64:
				if under > math.MaxUint16 {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for uint16", fmt.Sprintf("%d", under))
				}
				return uint16(under), nil
			case int:
				if under < 0 || under > math.MaxUint16 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for uint16", fmt.Sprintf("%d", under))
				}
				return uint16(under), nil
			case int64:
				if under < 0 || under > math.MaxUint16 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for uint16", fmt.Sprintf("%d", under))
				}
				return uint16(under), nil
			case json.Number:
				return uint16FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: uint16FromString,
		FromZero:   FromZero[uint16],
		Equals:     Equals[uint16],
	}
}

func uint32FromString(s string) (uint32, error) {
	i, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

// Uint32 returns a [TypeInfo] for uint32 values (range 0 to 4294967295).
func Uint32() TypeInfo[uint32] {
	return TypeInfo[uint32]{
		Description: "uint32",
		FromIFace: func(iFace interface{}) (uint32, error) {
			switch under := iFace.(type) {
			case uint32:
				return under, nil
			case uint64:
				if under > math.MaxUint32 {
					return 0, colerr.NewWrappedf(nil, "Uint64 value %s out of range for uint32", fmt.Sprintf("%d", under))
				}
				return uint32(under), nil
			case int:
				if under < 0 || under > math.MaxUint32 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for uint32", fmt.Sprintf("%d", under))
				}
				return uint32(under), nil
			case int64:
				if under < 0 || under > math.MaxUint32 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for uint32", fmt.Sprintf("%d", under))
				}
				return uint32(under), nil
			case json.Number:
				return uint32FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: uint32FromString,
		FromZero:   FromZero[uint32],
		Equals:     Equals[uint32],
	}
}

func uint64FromString(s string) (uint64, error) {
	return strconv.ParseUint(s, 0, 64)
}

// Uint64 returns a [TypeInfo] for uint64 values.
func Uint64() TypeInfo[uint64] {
	return TypeInfo[uint64]{
		Description: "uint64",
		FromIFace: func(iFace interface{}) (uint64, error) {
			switch under := iFace.(type) {
			case uint64:
				return under, nil
			case int:
				if under < 0 {
					return 0, colerr.NewWrappedf(nil, "Int value %s out of range for uint64", fmt.Sprintf("%d", under))
				}
				return uint64(under), nil
			case int64:
				if under < 0 {
					return 0, colerr.NewWrappedf(nil, "Int64 value %s out of range for uint64", fmt.Sprintf("%d", under))
				}
				return uint64(under), nil
			case json.Number:
				return uint64FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: uint64FromString,
		FromZero:   FromZero[uint64],
		Equals:     Equals[uint64],
	}
}

func float32FromString(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// Float32 returns a [TypeInfo] for float32 values.
func Float32() TypeInfo[float32] {
	return TypeInfo[float32]{
		Description: "float32",
		FromIFace: func(iFace interface{}) (float32, error) {
			switch under := iFace.(type) {
			case float32:
				return under, nil
			case float64:
				if under > math.MaxFloat32 || under < -math.MaxFloat32 {
					return 0, colerr.NewWrappedf(nil, "Float64 value %s out of range for float32", fmt.Sprintf("%g", under))
				}
				return float32(under), nil
			case int:
				return float32(under), nil
			case int64:
				return float32(under), nil
			case uint64:
				return float32(under), nil
			case json.Number:
				return float32FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: float32FromString,
		FromZero:   FromZero[float32],
		Equals:     Equals[float32],
	}
}

func float64FromString(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Float64 returns a [TypeInfo] for float64 values.
func Float64() TypeInfo[float64] {
	return TypeInfo[float64]{
		Description: "float64",
		FromIFace: func(iFace interface{}) (float64, error) {
			switch under := iFace.(type) {
			case float64:
				return under, nil
			case float32:
				return float64(under), nil
			case int:
				return float64(under), nil
			case int64:
				return float64(under), nil
			case uint64:
				return float64(under), nil
			case json.Number:
				return float64FromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: float64FromString,
		FromZero:   FromZero[float64],
		Equals:     Equals[float64],
	}
}

// Path returns a [TypeInfo] for [path.Path] values (file paths with ~ expansion support).
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
		return -1, errors.New("Empty string passed")
	}
	var r rune
	if rs := []rune(s); len(rs) != 1 {
		return emptyRune, errors.New("Runes shuld only be one character")
	} else {
		return r, nil
	}
}

// Rune returns a [TypeInfo] for single-rune values. Accepts exactly one character.
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

// String returns a [TypeInfo] for string values.
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
