package scalar

import (
	"fmt"

	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

type scalarValue[T any] struct {
	choices    []T
	defaultVal *T
	inner      contained.TypeInfo[T]
	val        *T
	updatedBy  value.UpdatedBy
}

type ScalarOpt[T any] func(*scalarValue[T])

func newScalarValue[T any](
	inner contained.TypeInfo[T],
	opts ...ScalarOpt[T],
) scalarValue[T] {
	empty := inner.FromZero()
	sv := scalarValue[T]{
		choices:    []T{},
		defaultVal: nil,
		inner:      inner,
		val:        &empty,
		updatedBy:  value.UpdatedByUnset,
	}
	for _, opt := range opts {
		opt(&sv)
	}
	return sv
}

func New[T any](hc contained.TypeInfo[T], opts ...ScalarOpt[T]) value.EmptyConstructor {
	return func() value.Value {
		s := newScalarValue(
			hc,
			opts...,
		)
		return &s
	}
}

func PointerTo[T any](addr *T) ScalarOpt[T] {
	return func(v *scalarValue[T]) {
		v.val = addr
	}
}

func Choices[T any](choices ...T) ScalarOpt[T] {
	return func(v *scalarValue[T]) {
		v.choices = choices
	}
}

func Default[T any](def T) ScalarOpt[T] {
	return func(v *scalarValue[T]) {
		v.defaultVal = &def
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
	return *v.val
}

func (v *scalarValue[_]) HasDefault() bool {
	return v.defaultVal != nil
}

func (v *scalarValue[T]) ReplaceFromInterface(iFace interface{}, u value.UpdatedBy) error {
	if v.updatedBy != value.UpdatedByUnset {
		return value.ErrUpdatedMoreThanOnce[T]{CurrentValue: *v.val, UpdatedBy: v.updatedBy}
	}
	val, err := v.inner.FromIFace(iFace)
	if err != nil {
		return err
	}
	*v.val = val
	v.updatedBy = u
	return nil
}

func (v *scalarValue[_]) String() string {
	return fmt.Sprint(*v.val)
}

func (v *scalarValue[T]) Update(s string, u value.UpdatedBy) error {
	if v.updatedBy != value.UpdatedByUnset {
		return value.ErrUpdatedMoreThanOnce[T]{CurrentValue: *v.val, UpdatedBy: v.updatedBy}
	}
	val, err := v.inner.FromString(s)
	if err != nil {
		return err
	}
	if !contained.WithinChoices(val, v.choices, v.inner.Equals) {
		return value.ErrInvalidChoice[T]{Choices: v.choices}
	}
	*v.val = val
	v.updatedBy = u
	return nil
}

func (v *scalarValue[_]) UpdatedBy() value.UpdatedBy {
	return v.updatedBy
}

func (v *scalarValue[T]) ReplaceFromDefault(u value.UpdatedBy) error {
	if v.updatedBy != value.UpdatedByUnset {
		return value.ErrUpdatedMoreThanOnce[T]{CurrentValue: *v.val, UpdatedBy: v.updatedBy}
	}
	if v.defaultVal != nil {
		v.updatedBy = u
		*v.val = *v.defaultVal
	}
	return nil
}
