package value

import "fmt"

//  -- SliceValue

// It doesn't really make sense to use the Value type as the type param because then I get a list of values, when I want a list of the type the value contains. Also, if I use the Value interface, I'll need to initialize it before I can call any of it's methods I think. Using the the TRU
// type params means that I can just use the constructor

type sliceValue[
	T any,
	FI fromIFaceFunc[T],
	FS fromStringFunc[T],
] struct {
	vals        []T
	description string
	fromIFace   FI
	fromString  FS
}

func newSliceValue[
	T any,
	FI fromIFaceFunc[T],
	FS fromStringFunc[T],
](
	vals []T,
	description string,
	fromIFace FI,
	fromString FS,
) sliceValue[T, FI, FS] {
	return sliceValue[T, FI, FS]{
		vals:        vals,
		description: description,
		fromIFace:   fromIFace,
		fromString:  fromString,
	}
}

func (v *sliceValue[_, _, _]) Description() string {
	return v.description
}

func (v *sliceValue[_, _, _]) Get() interface{} {
	return v.vals
}

func (v *sliceValue[T, _, _]) ReplaceFromInterface(iFace interface{}) error {
	under, ok := iFace.([]interface{})
	if !ok {
		return ErrIncompatibleInterface
	}

	new := []T{}

	// NOTE: in the previous impllementation, this method was able to use v.UpdateFromInterface  because it literally replaced all of the type's data, which was fine because the type was an alias for []T
	// For this one, we should use the type parameter's fromIFace method because we're only wanting to work on the contained []T slice

	for _, e := range under {
		underE, err := v.fromIFace(e)
		if err != nil {
			// TODO: this won't communicate to the caller *which* element is the wrong type
			return err
		}
		new = append(new, underE)
	}
	v.vals = new
	return nil
}

func (v *sliceValue[_, _, _]) String() string {
	return fmt.Sprint(v.vals)
}

func (v *sliceValue[_, _, _]) StringSlice() []string {
	ret := make([]string, 0, len(v.vals))
	for _, e := range v.vals {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}

func (sliceValue[_, _, _]) TypeInfo() TypeContainer {
	return TypeContainerSlice
}

func (v *sliceValue[_, _, _]) Update(s string) error {
	val, err := v.fromString(s)
	if err != nil {
		return err
	}
	v.vals = append(v.vals, val)
	return nil
}

func (v *sliceValue[_, _, _]) UpdateFromInterface(iFace interface{}) error {
	under, err := v.fromIFace(iFace)
	if err != nil {
		return ErrIncompatibleInterface
	}

	v.vals = append(v.vals, under)
	return nil
}
