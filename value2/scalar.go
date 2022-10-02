package value

import "fmt"

type scalarValue[T comparable] struct {
	common commonFields[T]
	inner  innerTypeInfo[T]
	val    T
}

func newScalarValue[T comparable](
	hardcoded innerTypeInfo[T],
	opts ...commonFieldsOpt[T],
) scalarValue[T] {
	sv := scalarValue[T]{
		common: commonFields[T]{},
		inner:  hardcoded,
		val:    hardcoded.empty(),
	}
	for _, opt := range opts {
		opt(&sv.common)
	}
	return sv
}

func Scalar[T comparable](hc innerTypeInfo[T], opts ...commonFieldsOpt[T]) EmptyConstructor {
	return func() (Value, error) {
		s := newScalarValue(
			hc,
			opts...,
		)
		return &s, nil
	}
}

func (v *scalarValue[_]) Choices() []string {
	return v.common.Choices()
}

func (v *scalarValue[_]) Description() string {
	return v.inner.description // TODO: will need to change this...
}

func (v *scalarValue[_]) Get() interface{} {
	return v.val
}

func (v *scalarValue[_]) ReplaceFromInterface(iFace interface{}) error {
	val, err := v.inner.fromIFace(iFace)
	if err != nil {
		return err
	}
	v.val = val
	return nil
}

func (v *scalarValue[_]) String() string {
	return fmt.Sprint(v.val)
}

func (v *scalarValue[_]) StringSlice() []string {
	return nil
}

func (scalarValue[_]) TypeInfo() TypeContainer {
	return TypeContainerScalar
}

func (v *scalarValue[_]) Update(s string) error {
	val, err := v.inner.fromString(s)
	if err != nil {
		return err
	}
	if !v.common.WithinChoices(val) {
		return ErrInvalidChoice
	}
	v.val = val
	return nil
}

func (v *scalarValue[_]) UpdateFromInterface(iFace interface{}) error {
	return v.ReplaceFromInterface(iFace)
}
