package scalar

import (
	"net/netip"
	"time"

	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

func Addr(opts ...ScalarOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.Addr(), opts...)
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

func Path(opts ...ScalarOpt[string]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

func Rune(opts ...ScalarOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

func String(opts ...ScalarOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
