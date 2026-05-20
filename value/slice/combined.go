package slice

import (
	"net/netip"
	"time"

	"go.bbkane.com/warg/path"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

// Addr returns an [value.EmptyConstructor] for a slice of [netip.Addr] values.
func Addr(opts ...SliceOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.NetIPAddr(), opts...)
}

// AddrPort returns an [value.EmptyConstructor] for a slice of [netip.AddrPort] values.
func AddrPort(opts ...SliceOpt[netip.AddrPort]) value.EmptyConstructor {
	return New(contained.AddrPort(), opts...)
}

// Bool returns an [value.EmptyConstructor] for a slice of bool values.
func Bool(opts ...SliceOpt[bool]) value.EmptyConstructor {
	return New(contained.Bool(), opts...)
}

// Duration returns an [value.EmptyConstructor] for a slice of [time.Duration] values.
func Duration(opts ...SliceOpt[time.Duration]) value.EmptyConstructor {
	return New(contained.Duration(), opts...)
}

// Int returns an [value.EmptyConstructor] for a slice of int values.
func Int(opts ...SliceOpt[int]) value.EmptyConstructor {
	return New(contained.Int(), opts...)
}

// Int8 returns an [value.EmptyConstructor] for a slice of int8 values.
func Int8(opts ...SliceOpt[int8]) value.EmptyConstructor {
	return New(contained.Int8(), opts...)
}

// Int16 returns an [value.EmptyConstructor] for a slice of int16 values.
func Int16(opts ...SliceOpt[int16]) value.EmptyConstructor {
	return New(contained.Int16(), opts...)
}

// Int32 returns an [value.EmptyConstructor] for a slice of int32 values.
func Int32(opts ...SliceOpt[int32]) value.EmptyConstructor {
	return New(contained.Int32(), opts...)
}

// Int64 returns an [value.EmptyConstructor] for a slice of int64 values.
func Int64(opts ...SliceOpt[int64]) value.EmptyConstructor {
	return New(contained.Int64(), opts...)
}

// Uint returns an [value.EmptyConstructor] for a slice of uint values.
func Uint(opts ...SliceOpt[uint]) value.EmptyConstructor {
	return New(contained.Uint(), opts...)
}

// Uint8 returns an [value.EmptyConstructor] for a slice of uint8 values.
func Uint8(opts ...SliceOpt[uint8]) value.EmptyConstructor {
	return New(contained.Uint8(), opts...)
}

// Uint16 returns an [value.EmptyConstructor] for a slice of uint16 values.
func Uint16(opts ...SliceOpt[uint16]) value.EmptyConstructor {
	return New(contained.Uint16(), opts...)
}

// Uint32 returns an [value.EmptyConstructor] for a slice of uint32 values.
func Uint32(opts ...SliceOpt[uint32]) value.EmptyConstructor {
	return New(contained.Uint32(), opts...)
}

// Uint64 returns an [value.EmptyConstructor] for a slice of uint64 values.
func Uint64(opts ...SliceOpt[uint64]) value.EmptyConstructor {
	return New(contained.Uint64(), opts...)
}

// Float32 returns an [value.EmptyConstructor] for a slice of float32 values.
func Float32(opts ...SliceOpt[float32]) value.EmptyConstructor {
	return New(contained.Float32(), opts...)
}

// Float64 returns an [value.EmptyConstructor] for a slice of float64 values.
func Float64(opts ...SliceOpt[float64]) value.EmptyConstructor {
	return New(contained.Float64(), opts...)
}

// Path returns an [value.EmptyConstructor] for a slice of [path.Path] values.
func Path(opts ...SliceOpt[path.Path]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

// Rune returns an [value.EmptyConstructor] for a slice of rune values.
func Rune(opts ...SliceOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

// String returns an [value.EmptyConstructor] for a slice of string values.
func String(opts ...SliceOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
