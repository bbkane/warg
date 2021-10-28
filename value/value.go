package value

import (
	"errors"
)

type typeInfo int64

// These constants describe the container type of a Value.
const (
	TypeInfoScalar typeInfo = iota + 1
	TypeInfoSlice
	TypeInfoMap
)

// Value is a "generic" type to store different types into flags
// Inspired by https://golang.org/src/flag/flag.go .
// There are two underlying "type" families designed to fit in Value:
//
// - scalar types (Int, String, ...)
//
// - container types (IntSlice, StringMap, ...)
type Value interface {

	// Description of the type. useful for help messages
	Description() string

	// Get returns the underlying value. It's meant to be type asserted against
	// Example: myInt := v.(int)
	Get() interface{}

	// ReplaceFromInterface replaces a value with one found in an interface (useful for configs)
	ReplaceFromInterface(interface{}) error

	// TypeInfo specifies whether what "overall" type of value this is - scalar, slice, etc.
	TypeInfo() typeInfo

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
type EmptyConstructor = func() (Value, error)

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")
