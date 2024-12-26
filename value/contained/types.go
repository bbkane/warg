package contained

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
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

type TypeInfo[T comparable] struct {
	Description string

	FromIFace func(iFace interface{}) (T, error)

	FromString func(string) (T, error)

	// Initalized to the Empty value, but used for updating stuff in the container type
	Empty func() T
}

func Addr() TypeInfo[netip.Addr] {
	return TypeInfo[netip.Addr]{
		Description: "IP address",
		Empty: func() netip.Addr {
			return netip.Addr{}
		},
		FromIFace: func(iFace interface{}) (netip.Addr, error) {
			switch under := iFace.(type) {
			case netip.Addr:
				return under, nil
			case []byte:
				ip, ok := netip.AddrFromSlice(under)
				if !ok {
					return netip.Addr{}, fmt.Errorf("Could not convert %s to netip.Addr", string(under))
				}
				return ip, nil
			case string:
				return netip.ParseAddr(under)
			default:
				return netip.Addr{}, ErrIncompatibleInterface
			}
		},

		FromString: netip.ParseAddr,
	}
}

func AddrPort() TypeInfo[netip.AddrPort] {
	return TypeInfo[netip.AddrPort]{
		Description: "IP and Port number separated by a colon: ip:port ",
		Empty: func() netip.AddrPort {
			return netip.AddrPort{}
		},
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
	}
}

func Bool() TypeInfo[bool] {
	return TypeInfo[bool]{
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

func Duration() TypeInfo[time.Duration] {
	return TypeInfo[time.Duration]{
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
			case json.Number:
				return intFromString(string(under))
			default:
				return 0, ErrIncompatibleInterface
			}
		},
		FromString: intFromString,
		Empty:      func() int { return 0 },
	}
}

func pathFromString(s string) (string, error) {
	expanded, err := homedir.Expand(s)
	if err != nil {
		return "", err
	}
	return expanded, nil
}

func Path() TypeInfo[string] {
	return TypeInfo[string]{
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

func Rune() TypeInfo[rune] {
	return TypeInfo[rune]{
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

func String() TypeInfo[string] {
	return TypeInfo[string]{
		Description: "string",
		Empty:       func() string { return "" },
		FromIFace: func(iFace interface{}) (string, error) {
			under, ok := iFace.(string)
			if !ok {
				return "", ErrIncompatibleInterface
			}
			return under, nil
		},
		FromString: identity[string],
	}
}
