package value

import "fmt"

type scalarValue[
	T any,
	FI fromIFaceFunc[T],
	FS fromStringFunc[T],
] struct {
	val         T
	description string
	fromIFace   FI
	fromString  FS
}

func newScalarValue[
	T any,
	FI fromIFaceFunc[T],
	FS fromStringFunc[T],
](
	val T,
	description string,
	fromIFace FI,
	fromString FS,
) scalarValue[T, FI, FS] {
	return scalarValue[T, FI, FS]{
		val:         val,
		description: description,
		fromIFace:   fromIFace,
		fromString:  fromString,
	}
}

func (v *scalarValue[_, _, _]) Description() string {
	return v.description
}

func (v *scalarValue[_, _, _]) Get() interface{} {
	return v.val
}

func (v *scalarValue[_, _, _]) ReplaceFromInterface(iFace interface{}) error {
	val, err := v.fromIFace(iFace)
	if err != nil {
		return err
	}
	v.val = val
	return nil
}

func (v *scalarValue[_, _, _]) String() string {
	return fmt.Sprint(v.val)
}

func (scalarValue[_, _, _]) StringSlice() []string {
	return nil
}

func (scalarValue[_, _, _]) TypeInfo() TypeContainer {
	return TypeContainerScalar
}

func (v *scalarValue[_, _, _]) Update(s string) error {
	val, err := v.fromString(s)
	if err != nil {
		return err
	}
	v.val = val
	return nil
}

func (v *scalarValue[_, _, _]) UpdateFromInterface(iFace interface{}) error {
	return v.ReplaceFromInterface(iFace)
}
