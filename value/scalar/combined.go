package scalar

import (
	"net/netip"
	"time"

	"go.bbkane.com/warg/path"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

func Addr(opts ...ScalarOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.NetIPAddr(), opts...)
}

func AddrPort(opts ...ScalarOpt[netip.AddrPort]) value.EmptyConstructor {
	return New(contained.AddrPort(), opts...)
}

func Bool(opts ...ScalarOpt[bool]) value.EmptyConstructor {
	return New(contained.Bool(), opts...)
}

func Duration(opts ...ScalarOpt[time.Duration]) value.EmptyConstructor {
	return New(contained.Duration(), opts...)
}

func Int(opts ...ScalarOpt[int]) value.EmptyConstructor {
	return New(contained.Int(), opts...)
}

func Int8(opts ...ScalarOpt[int8]) value.EmptyConstructor {
	return New(contained.Int8(), opts...)
}

func Int16(opts ...ScalarOpt[int16]) value.EmptyConstructor {
	return New(contained.Int16(), opts...)
}

func Int32(opts ...ScalarOpt[int32]) value.EmptyConstructor {
	return New(contained.Int32(), opts...)
}

func Int64(opts ...ScalarOpt[int64]) value.EmptyConstructor {
	return New(contained.Int64(), opts...)
}

func Uint(opts ...ScalarOpt[uint]) value.EmptyConstructor {
	return New(contained.Uint(), opts...)
}

func Uint8(opts ...ScalarOpt[uint8]) value.EmptyConstructor {
	return New(contained.Uint8(), opts...)
}

func Uint16(opts ...ScalarOpt[uint16]) value.EmptyConstructor {
	return New(contained.Uint16(), opts...)
}

func Uint32(opts ...ScalarOpt[uint32]) value.EmptyConstructor {
	return New(contained.Uint32(), opts...)
}

func Uint64(opts ...ScalarOpt[uint64]) value.EmptyConstructor {
	return New(contained.Uint64(), opts...)
}

func Float32(opts ...ScalarOpt[float32]) value.EmptyConstructor {
	return New(contained.Float32(), opts...)
}

func Float64(opts ...ScalarOpt[float64]) value.EmptyConstructor {
	return New(contained.Float64(), opts...)
}

func Path(opts ...ScalarOpt[path.Path]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

func Rune(opts ...ScalarOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

func String(opts ...ScalarOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
