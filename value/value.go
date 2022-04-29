package value

import (
	"errors"
	"fmt"
)

type TypeContainer int64

// These constants describe the container type of a Value.
const (
	TypeContainerScalar TypeContainer = iota + 1
	TypeContainerSlice
	TypeContainerMap
)

// Value is a "generic" type to store different types into flags
// Inspired by https://golang.org/src/flag/flag.go .
// There are two underlying "type" families designed to fit in Value:
// scalar types (Int, String, ...) and container types (IntSlice, StringMap, ...).
type Value interface {

	// Description of the type. useful for help messages. Should not be used as an ID.
	Description() string

	// Get returns the underlying value. It's meant to be type asserted against
	// Example: myInt := v.(int)
	Get() interface{}

	// Len returns 0 for scalar Values and len(underlyingValue) for container Values.
	// TODO: I think this will be useful when/if I start enforcing flag grouping (like grabbits subreddit params).
	// Those should all have the same length and the same source. I don't think I *need* it now, so leaving it out.
	// Len() int

	// ReplaceFromInterface replaces a value with one found in an interface.
	// Useful to update a Value from a config.
	ReplaceFromInterface(interface{}) error

	// String returns a string ready to be printed!
	String() string

	// StringSlice returns a []string ready to be printed for slice values and nil for others
	StringSlice() []string

	// TypeInfo specifies whether what "overall" type of value this is - scalar, slice, etc.
	TypeInfo() TypeContainer

	// Update appends to container type Values from a string (useful for CLI flags, env vars, default values)
	// and replaces scalar Values
	Update(string) error

	// UpdateFromInterface updates a container type Value from an interface (useful for configs)
	// and replaces scalar values (for scalar values, UpdateFromInterface is the same as ReplaceFromInterface).
	// Note that UpdateFromInterface must be called with the "contained" type for container type Values
	// For example, the StringSlice.UpdateFromInterface
	// must be called with a string, not a []string
	// It returns ErrIncompatibleInterface if the interface can't be decoded
	UpdateFromInterface(interface{}) error
}

// EmptyConstructur just builds a new value.
// Useful to create new values as well as initialize them
type EmptyConstructor func() (Value, error)

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")

// -- ScalarValue

type fromIFaceFunc[T any] func(interface{}) (T, error)

func fromIFaceEnum[T comparable](fromIFace fromIFaceFunc[T], choices ...T) fromIFaceFunc[T] {
	return func(iFace interface{}) (T, error) {
		val, err := fromIFace(iFace)
		if err != nil {
			return val, err
		}
		for _, choice := range choices {
			if val == choice {
				return val, nil
			}
		}
		return val, fmt.Errorf("interface enum update invalid choice: available: %v: choice: %v", choices, val)
	}
}

type fromStringFunc[T any] func(string) (T, error)

func fromStringEnum[T comparable](fromString fromStringFunc[T], choices ...T) fromStringFunc[T] {
	return func(s string) (T, error) {
		val, err := fromString(s)
		if err != nil {
			return val, err
		}
		for _, choice := range choices {
			if val == choice {
				return val, nil
			}
		}
		return val, fmt.Errorf("string enum update invalid choice: available: %v: choice: %v", choices, val)
	}
}

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
