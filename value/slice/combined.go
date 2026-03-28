package slice

import (
	"net/netip"
	"time"

	"go.bbkane.com/warg/path"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

func Addr(opts ...SliceOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.NetIPAddr(), opts...)
}

func AddrPort(opts ...SliceOpt[netip.AddrPort]) value.EmptyConstructor {
	return New(contained.AddrPort(), opts...)
}

func Bool(opts ...SliceOpt[bool]) value.EmptyConstructor {
	return New(contained.Bool(), opts...)
}

func Duration(opts ...SliceOpt[time.Duration]) value.EmptyConstructor {
	return New(contained.Duration(), opts...)
}

func Int(opts ...SliceOpt[int]) value.EmptyConstructor {
	return New(contained.Int(), opts...)
}

func Int8(opts ...SliceOpt[int8]) value.EmptyConstructor {
	return New(contained.Int8(), opts...)
}

func Int16(opts ...SliceOpt[int16]) value.EmptyConstructor {
	return New(contained.Int16(), opts...)
}

func Int32(opts ...SliceOpt[int32]) value.EmptyConstructor {
	return New(contained.Int32(), opts...)
}

func Int64(opts ...SliceOpt[int64]) value.EmptyConstructor {
	return New(contained.Int64(), opts...)
}

func Uint(opts ...SliceOpt[uint]) value.EmptyConstructor {
	return New(contained.Uint(), opts...)
}

func Uint8(opts ...SliceOpt[uint8]) value.EmptyConstructor {
	return New(contained.Uint8(), opts...)
}

func Uint16(opts ...SliceOpt[uint16]) value.EmptyConstructor {
	return New(contained.Uint16(), opts...)
}

func Uint32(opts ...SliceOpt[uint32]) value.EmptyConstructor {
	return New(contained.Uint32(), opts...)
}

func Uint64(opts ...SliceOpt[uint64]) value.EmptyConstructor {
	return New(contained.Uint64(), opts...)
}

func Float32(opts ...SliceOpt[float32]) value.EmptyConstructor {
	return New(contained.Float32(), opts...)
}

func Float64(opts ...SliceOpt[float64]) value.EmptyConstructor {
	return New(contained.Float64(), opts...)
}

func Path(opts ...SliceOpt[path.Path]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

func Rune(opts ...SliceOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

func String(opts ...SliceOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
