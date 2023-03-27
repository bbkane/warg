package slice

import (
	"net/netip"
	"time"

	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

func Addr(opts ...SliceOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.Addr(), opts...)
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

func Path(opts ...SliceOpt[string]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

func Rune(opts ...SliceOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

func String(opts ...SliceOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
