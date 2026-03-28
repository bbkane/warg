package dict

import (
	"net/netip"
	"time"

	"go.bbkane.com/warg/path"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

func Addr(opts ...DictOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.NetIPAddr(), opts...)
}

func AddrPort(opts ...DictOpt[netip.AddrPort]) value.EmptyConstructor {
	return New(contained.AddrPort(), opts...)
}

func Bool(opts ...DictOpt[bool]) value.EmptyConstructor {
	return New(contained.Bool(), opts...)
}

func Duration(opts ...DictOpt[time.Duration]) value.EmptyConstructor {
	return New(contained.Duration(), opts...)
}

func Int(opts ...DictOpt[int]) value.EmptyConstructor {
	return New(contained.Int(), opts...)
}

func Int8(opts ...DictOpt[int8]) value.EmptyConstructor {
	return New(contained.Int8(), opts...)
}

func Int16(opts ...DictOpt[int16]) value.EmptyConstructor {
	return New(contained.Int16(), opts...)
}

func Int32(opts ...DictOpt[int32]) value.EmptyConstructor {
	return New(contained.Int32(), opts...)
}

func Int64(opts ...DictOpt[int64]) value.EmptyConstructor {
	return New(contained.Int64(), opts...)
}

func Uint(opts ...DictOpt[uint]) value.EmptyConstructor {
	return New(contained.Uint(), opts...)
}

func Uint8(opts ...DictOpt[uint8]) value.EmptyConstructor {
	return New(contained.Uint8(), opts...)
}

func Uint16(opts ...DictOpt[uint16]) value.EmptyConstructor {
	return New(contained.Uint16(), opts...)
}

func Uint32(opts ...DictOpt[uint32]) value.EmptyConstructor {
	return New(contained.Uint32(), opts...)
}

func Uint64(opts ...DictOpt[uint64]) value.EmptyConstructor {
	return New(contained.Uint64(), opts...)
}

func Float32(opts ...DictOpt[float32]) value.EmptyConstructor {
	return New(contained.Float32(), opts...)
}

func Float64(opts ...DictOpt[float64]) value.EmptyConstructor {
	return New(contained.Float64(), opts...)
}

func Path(opts ...DictOpt[path.Path]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

func Rune(opts ...DictOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

func String(opts ...DictOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
