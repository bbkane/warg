package slice

import (
	"fmt"

	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

type sliceValue[T comparable] struct {
	choices     []T
	defaultVals []T
	hasDefault  bool
	inner       contained.TypeInfo[T]
	vals        []T
	updatedBy   value.UpdatedBy
}

type SliceOpt[T comparable] func(*sliceValue[T])

func New[T comparable](hc contained.TypeInfo[T], opts ...SliceOpt[T]) value.EmptyConstructor {
	return func() value.Value {
		sv := sliceValue[T]{
			choices:     []T{},
			defaultVals: nil,
			hasDefault:  false,
			inner:       hc,
			vals:        nil,
			updatedBy:   value.UpdatedByUnset,
		}
		for _, opt := range opts {
			opt(&sv)
		}
		return &sv
	}
}

func Choices[T comparable](choices ...T) SliceOpt[T] {
	return func(v *sliceValue[T]) {
		v.choices = choices
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

func (v *sliceValue[T]) ReplaceFromInterface(iFace interface{}, u value.UpdatedBy) error {
	under, ok := iFace.([]interface{})
	if !ok {
		return contained.ErrIncompatibleInterface
	}

	newVals := []T{}
	for _, e := range under {
		underE, err := v.inner.FromIFace(e)
		if err != nil {
			// TODO: this won't communicate to the caller *which* element is the wrong type
			return err
		}
		newVals = append(newVals, underE)
	}
	v.updatedBy = u
	v.vals = newVals
	return nil
}

func (v *sliceValue[_]) StringSlice() []string {
	ret := make([]string, 0, len(v.vals))
	for _, e := range v.vals {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
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
		return value.ErrInvalidChoice[T]{Choices: v.choices}
	}
	v.vals = append(v.vals, val)
	return nil
}

func (v *sliceValue[_]) Update(s string, u value.UpdatedBy) error {
	val, err := v.inner.FromString(s)
	if err != nil {
		return err
	}
	err = v.update(val)
	if err != nil {
		return err
	}
	v.updatedBy = u
	return nil
}

func (v *sliceValue[_]) UpdatedBy() value.UpdatedBy {
	return v.updatedBy
}

func (v *sliceValue[_]) ReplaceFromDefault(u value.UpdatedBy) error {
	if v.hasDefault {
		v.vals = v.defaultVals
		v.updatedBy = u
	}
	return nil
}

func (v *sliceValue[_]) AppendFromInterface(iFace interface{}, u value.UpdatedBy) error {
	val, err := v.inner.FromIFace(iFace)
	if err != nil {
		return err
	}
	err = v.update(val)
	if err != nil {
		return err
	}
	v.updatedBy = u
	return nil
}
