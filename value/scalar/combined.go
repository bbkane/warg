package scalar

import (
	"net/netip"
	"time"

	"go.bbkane.com/warg/path"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

// Addr returns an [value.EmptyConstructor] for a scalar [netip.Addr] flag.
func Addr(opts ...ScalarOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.NetIPAddr(), opts...)
}

// AddrPort returns an [value.EmptyConstructor] for a scalar [netip.AddrPort] flag (ip:port).
func AddrPort(opts ...ScalarOpt[netip.AddrPort]) value.EmptyConstructor {
	return New(contained.AddrPort(), opts...)
}

// Bool returns an [value.EmptyConstructor] for a scalar bool flag.
func Bool(opts ...ScalarOpt[bool]) value.EmptyConstructor {
	return New(contained.Bool(), opts...)
}

// Duration returns an [value.EmptyConstructor] for a scalar [time.Duration] flag.
func Duration(opts ...ScalarOpt[time.Duration]) value.EmptyConstructor {
	return New(contained.Duration(), opts...)
}

// DateTimeRFC3339 returns an [value.EmptyConstructor] for a scalar [time.Time] flag in RFC3339 format.
func DateTimeRFC3339(opts ...ScalarOpt[time.Time]) value.EmptyConstructor {
	return New(contained.DateTimeRFC3339(), opts...)
}

// Int returns an [value.EmptyConstructor] for a scalar int flag.
func Int(opts ...ScalarOpt[int]) value.EmptyConstructor {
	return New(contained.Int(), opts...)
}

// Int8 returns an [value.EmptyConstructor] for a scalar int8 flag.
func Int8(opts ...ScalarOpt[int8]) value.EmptyConstructor {
	return New(contained.Int8(), opts...)
}

// Int16 returns an [value.EmptyConstructor] for a scalar int16 flag.
func Int16(opts ...ScalarOpt[int16]) value.EmptyConstructor {
	return New(contained.Int16(), opts...)
}

// Int32 returns an [value.EmptyConstructor] for a scalar int32 flag.
func Int32(opts ...ScalarOpt[int32]) value.EmptyConstructor {
	return New(contained.Int32(), opts...)
}

// Int64 returns an [value.EmptyConstructor] for a scalar int64 flag.
func Int64(opts ...ScalarOpt[int64]) value.EmptyConstructor {
	return New(contained.Int64(), opts...)
}

// Uint returns an [value.EmptyConstructor] for a scalar uint flag.
func Uint(opts ...ScalarOpt[uint]) value.EmptyConstructor {
	return New(contained.Uint(), opts...)
}

// Uint8 returns an [value.EmptyConstructor] for a scalar uint8 flag.
func Uint8(opts ...ScalarOpt[uint8]) value.EmptyConstructor {
	return New(contained.Uint8(), opts...)
}

// Uint16 returns an [value.EmptyConstructor] for a scalar uint16 flag.
func Uint16(opts ...ScalarOpt[uint16]) value.EmptyConstructor {
	return New(contained.Uint16(), opts...)
}

// Uint32 returns an [value.EmptyConstructor] for a scalar uint32 flag.
func Uint32(opts ...ScalarOpt[uint32]) value.EmptyConstructor {
	return New(contained.Uint32(), opts...)
}

// Uint64 returns an [value.EmptyConstructor] for a scalar uint64 flag.
func Uint64(opts ...ScalarOpt[uint64]) value.EmptyConstructor {
	return New(contained.Uint64(), opts...)
}

// Float32 returns an [value.EmptyConstructor] for a scalar float32 flag.
func Float32(opts ...ScalarOpt[float32]) value.EmptyConstructor {
	return New(contained.Float32(), opts...)
}

// Float64 returns an [value.EmptyConstructor] for a scalar float64 flag.
func Float64(opts ...ScalarOpt[float64]) value.EmptyConstructor {
	return New(contained.Float64(), opts...)
}

// Path returns an [value.EmptyConstructor] for a scalar [path.Path] flag (supports ~ expansion).
func Path(opts ...ScalarOpt[path.Path]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

// Rune returns an [value.EmptyConstructor] for a scalar rune flag.
func Rune(opts ...ScalarOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

// String returns an [value.EmptyConstructor] for a scalar string flag.
func String(opts ...ScalarOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
