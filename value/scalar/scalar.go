package scalar

import (
	"fmt"

	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

type scalarValue[T comparable] struct {
	choices    []T
	defaultVal *T
	inner      contained.TypeInfo[T]
	val        T
}

type ScalarOpt[T comparable] func(*scalarValue[T])

func newScalarValue[T comparable](
	inner contained.TypeInfo[T],
	opts ...ScalarOpt[T],
) scalarValue[T] {
	sv := scalarValue[T]{
		choices:    []T{},
		defaultVal: nil,
		inner:      inner,
		val:        inner.Empty(),
	}
	for _, opt := range opts {
		opt(&sv)
	}
	return sv
}

func New[T comparable](hc contained.TypeInfo[T], opts ...ScalarOpt[T]) value.EmptyConstructor {
	return func() (value.Value, error) {
		s := newScalarValue(
			hc,
			opts...,
		)
		return &s, nil
	}
}

func Choices[T comparable](choices ...T) ScalarOpt[T] {
	return func(v *scalarValue[T]) {
		for _, ch := range choices {
			newChoice, err := v.inner.FromInstance(ch)
			if err != nil {
				panic("error constructing default value: " + fmt.Sprint(ch) + " : " + err.Error())
			}
			v.choices = append(v.choices, newChoice)
		}
	}
}

func Default[T comparable](def T) ScalarOpt[T] {
	return func(v *scalarValue[T]) {
		newDef, err := v.inner.FromInstance(def)
		// I think it's ok to panic here. This is from author code,
		// not user values
		if err != nil {
			panic("error constructing default value: " + fmt.Sprint(def) + " : " + err.Error())
		}
		v.defaultVal = &newDef
	}
}

func (v *scalarValue[_]) Choices() []string {
	ret := []string{}
	for _, e := range v.choices {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}

func (v *scalarValue[_]) DefaultString() string {
	if v.defaultVal == nil {
		return ""
	}
	// because we're representing the default as a ptr, we need to deref it to get a value
	return fmt.Sprint(*v.defaultVal)
}

func (v *scalarValue[_]) Description() string {
	return v.inner.Description
}

func (v *scalarValue[_]) Get() interface{} {
	return v.val
}

func (v *scalarValue[_]) HasDefault() bool {
	return v.defaultVal != nil
}

func (v *scalarValue[_]) ReplaceFromInterface(iFace interface{}) error {
	val, err := v.inner.FromIFace(iFace)
	if err != nil {
		return err
	}
	v.val = val
	return nil
}

func (v *scalarValue[_]) String() string {
	return fmt.Sprint(v.val)
}

func withinChoices[T comparable](val T, choices []T) bool {
	// User didn't constrain choices
	if len(choices) == 0 {
		return true
	}
	for _, choice := range choices {
		if val == choice {
			return true
		}
	}
	return false
}

func (v *scalarValue[T]) Update(s string) error {
	val, err := v.inner.FromString(s)
	if err != nil {
		return err
	}
	if !withinChoices(val, v.choices) {
		return value.ErrInvalidChoice[T]{Choices: v.choices}
	}
	v.val = val
	return nil
}

func (v *scalarValue[_]) ReplaceFromDefault() {
	if v.defaultVal != nil {
		v.val = *v.defaultVal
	}
}
