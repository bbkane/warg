package slice

import (
	"time"

	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

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
