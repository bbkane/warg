package dict

import (
	"net/netip"
	"time"

	"go.bbkane.com/warg/path"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

// Addr returns an [value.EmptyConstructor] for a dict with [netip.Addr] values.
func Addr(opts ...DictOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.NetIPAddr(), opts...)
}

// AddrPort returns an [value.EmptyConstructor] for a dict with [netip.AddrPort] values.
func AddrPort(opts ...DictOpt[netip.AddrPort]) value.EmptyConstructor {
	return New(contained.AddrPort(), opts...)
}

// Bool returns an [value.EmptyConstructor] for a dict with bool values.
func Bool(opts ...DictOpt[bool]) value.EmptyConstructor {
	return New(contained.Bool(), opts...)
}

// Duration returns an [value.EmptyConstructor] for a dict with [time.Duration] values.
func Duration(opts ...DictOpt[time.Duration]) value.EmptyConstructor {
	return New(contained.Duration(), opts...)
}

// DateTimeRFC3339 returns an [value.EmptyConstructor] for a dict with [time.Time] values in RFC3339 format.
func DateTimeRFC3339(opts ...DictOpt[time.Time]) value.EmptyConstructor {
	return New(contained.DateTimeRFC3339(), opts...)
}

// Int returns an [value.EmptyConstructor] for a dict with int values.
func Int(opts ...DictOpt[int]) value.EmptyConstructor {
	return New(contained.Int(), opts...)
}

// Int8 returns an [value.EmptyConstructor] for a dict with int8 values.
func Int8(opts ...DictOpt[int8]) value.EmptyConstructor {
	return New(contained.Int8(), opts...)
}

// Int16 returns an [value.EmptyConstructor] for a dict with int16 values.
func Int16(opts ...DictOpt[int16]) value.EmptyConstructor {
	return New(contained.Int16(), opts...)
}

// Int32 returns an [value.EmptyConstructor] for a dict with int32 values.
func Int32(opts ...DictOpt[int32]) value.EmptyConstructor {
	return New(contained.Int32(), opts...)
}

// Int64 returns an [value.EmptyConstructor] for a dict with int64 values.
func Int64(opts ...DictOpt[int64]) value.EmptyConstructor {
	return New(contained.Int64(), opts...)
}

// Uint returns an [value.EmptyConstructor] for a dict with uint values.
func Uint(opts ...DictOpt[uint]) value.EmptyConstructor {
	return New(contained.Uint(), opts...)
}

// Uint8 returns an [value.EmptyConstructor] for a dict with uint8 values.
func Uint8(opts ...DictOpt[uint8]) value.EmptyConstructor {
	return New(contained.Uint8(), opts...)
}

// Uint16 returns an [value.EmptyConstructor] for a dict with uint16 values.
func Uint16(opts ...DictOpt[uint16]) value.EmptyConstructor {
	return New(contained.Uint16(), opts...)
}

// Uint32 returns an [value.EmptyConstructor] for a dict with uint32 values.
func Uint32(opts ...DictOpt[uint32]) value.EmptyConstructor {
	return New(contained.Uint32(), opts...)
}

// Uint64 returns an [value.EmptyConstructor] for a dict with uint64 values.
func Uint64(opts ...DictOpt[uint64]) value.EmptyConstructor {
	return New(contained.Uint64(), opts...)
}

// Float32 returns an [value.EmptyConstructor] for a dict with float32 values.
func Float32(opts ...DictOpt[float32]) value.EmptyConstructor {
	return New(contained.Float32(), opts...)
}

// Float64 returns an [value.EmptyConstructor] for a dict with float64 values.
func Float64(opts ...DictOpt[float64]) value.EmptyConstructor {
	return New(contained.Float64(), opts...)
}

// Path returns an [value.EmptyConstructor] for a dict with [path.Path] values.
func Path(opts ...DictOpt[path.Path]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

// Rune returns an [value.EmptyConstructor] for a dict with rune values.
func Rune(opts ...DictOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

// String returns an [value.EmptyConstructor] for a dict with string values.
func String(opts ...DictOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
