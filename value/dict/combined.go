package dict

import (
	"net/netip"
	"time"

	"go.bbkane.com/warg/path"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

func Addr(opts ...DictOpt[netip.Addr]) value.EmptyConstructor {
	return New(contained.Addr(), opts...)
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

func Path(opts ...DictOpt[path.Path]) value.EmptyConstructor {
	return New(contained.Path(), opts...)
}

func Rune(opts ...DictOpt[rune]) value.EmptyConstructor {
	return New(contained.Rune(), opts...)
}

func String(opts ...DictOpt[string]) value.EmptyConstructor {
	return New(contained.String(), opts...)
}
