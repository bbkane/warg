package slice

import (
	"fmt"

	value "go.bbkane.com/warg/value2"
	"go.bbkane.com/warg/value2/contained"
)

type sliceValue[T comparable] struct {
	choices     []T
	defaultVals []T
	hasDefault  bool
	inner       contained.ContainedTypeInfo[T]
	vals        []T
}

type SliceOpt[T comparable] func(*sliceValue[T])

func newSliceValue[T comparable](
	inner contained.ContainedTypeInfo[T],
	opts ...SliceOpt[T],
) sliceValue[T] {
	sv := sliceValue[T]{
		choices: []T{},
		inner:   inner,
		vals:    nil,
	}
	for _, opt := range opts {
		opt(&sv)
	}
	return sv
}

func New[T comparable](hc contained.ContainedTypeInfo[T], opts ...SliceOpt[T]) value.EmptyConstructor {
	return func() (value.Value, error) {
		s := newSliceValue(
			hc,
			opts...,
		)
		return &s, nil
	}
}

func Choices[T comparable](choices ...T) SliceOpt[T] {
	return func(cf *sliceValue[T]) {
		cf.choices = choices
	}
}

func Default[T comparable](def []T) SliceOpt[T] {
	return func(cf *sliceValue[T]) {
		cf.defaultVals = def
		cf.hasDefault = true
	}
}

func (v *sliceValue[_]) Choices() []string {
	ret := []string{}
	for _, e := range v.choices {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}

func (v *sliceValue[_]) DefaultString() string {
	if !v.hasDefault {
		return ""
	}
	return fmt.Sprint(v.defaultVals)
}

func (v *sliceValue[_]) DefaultStringSlice() []string {
	// TODO: no copy paste
	ret := make([]string, 0, len(v.defaultVals))
	for _, e := range v.defaultVals {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}

func (v *sliceValue[_]) Description() string {
	return v.inner.Description
}

func (v *sliceValue[_]) Get() interface{} {
	return v.vals
}

func (v *sliceValue[_]) HasDefault() bool {
	return v.hasDefault
}

func (v *sliceValue[T]) ReplaceFromInterface(iFace interface{}) error {
	under, ok := iFace.([]interface{})
	if !ok {
		return contained.ErrIncompatibleInterface
	}

	new := []T{}
	for _, e := range under {
		underE, err := v.inner.FromIFace(e)
		if err != nil {
			// TODO: this won't communicate to the caller *which* element is the wrong type
			return err
		}
		new = append(new, underE)
	}
	v.vals = new
	return nil
}

func (v *sliceValue[_]) String() string {
	return fmt.Sprint(v.vals)
}

func (v *sliceValue[_]) StringSlice() []string {
	ret := make([]string, 0, len(v.vals))
	for _, e := range v.vals {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}

func (sliceValue[_]) TypeInfo() value.TypeContainer {
	return value.TypeContainerSlice
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

func (v *sliceValue[T]) update(val T) error {
	if !withinChoices(val, v.choices) {
		return value.ErrInvalidChoice
	}
	v.vals = append(v.vals, val)
	return nil
}

func (v *sliceValue[_]) Update(s string) error {
	val, err := v.inner.FromString(s)
	if err != nil {
		return err
	}
	return v.update(val)
}

func (v *sliceValue[_]) UpdateFromDefault() {
	if v.hasDefault {
		v.vals = v.defaultVals
	}
}

func (v *sliceValue[_]) UpdateFromInterface(iFace interface{}) error {
	val, err := v.inner.FromIFace(iFace)
	if err != nil {
		return err
	}
	return v.update(val)
}