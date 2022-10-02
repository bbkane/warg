package value

import "fmt"

type sliceValue[T comparable] struct {
	common commonFields[T]
	inner  innerTypeInfo[T]
	vals   []T
}

func newSliceValue[T comparable](
	inner innerTypeInfo[T],
	opts ...commonFieldsOpt[T],
) sliceValue[T] {
	sv := sliceValue[T]{
		common: commonFields[T]{},
		inner:  inner,
		vals:   nil,
	}
	for _, opt := range opts {
		opt(&sv.common)
	}
	return sv
}

func Slice[T comparable](hc innerTypeInfo[T], opts ...commonFieldsOpt[T]) EmptyConstructor {
	return func() (Value, error) {
		s := newSliceValue(
			hc,
			opts...,
		)
		return &s, nil
	}
}

func (v *sliceValue[_]) Choices() []string {
	return v.common.Choices()
}

func (v *sliceValue[_]) Description() string {
	return v.inner.description
}

func (v *sliceValue[_]) Get() interface{} {
	return v.vals
}

func (v *sliceValue[T]) ReplaceFromInterface(iFace interface{}) error {
	under, ok := iFace.([]interface{})
	if !ok {
		return ErrIncompatibleInterface
	}

	new := []T{}
	for _, e := range under {
		underE, err := v.inner.fromIFace(e)
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

func (sliceValue[_]) TypeInfo() TypeContainer {
	return TypeContainerSlice
}

func (v *sliceValue[T]) update(val T) error {
	if !v.common.WithinChoices(val) {
		return ErrInvalidChoice
	}
	v.vals = append(v.vals, val)
	return nil
}

func (v *sliceValue[_]) Update(s string) error {
	val, err := v.inner.fromString(s)
	if err != nil {
		return err
	}
	return v.update(val)
}

func (v *sliceValue[_]) UpdateFromInterface(iFace interface{}) error {
	val, err := v.inner.fromIFace(iFace)
	if err != nil {
		return err
	}
	return v.update(val)
}
