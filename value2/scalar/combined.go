package scalar

import (
	"time"

	value "go.bbkane.com/warg/value2"
	"go.bbkane.com/warg/value2/contained"
)

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
