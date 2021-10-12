package value

import (
	"errors"
)

// Value is a "generic" type that lets me store different types into flags
//  ~Stolen from~ "Inspired by" https://golang.org/src/flag/flag.go?#L138
// There are two underlying "type" families designed to fit in Value:
// - scalar types (Int, String, ...)
// - container types (IntSlice, StringMap, ...)
type Value interface {
	// Get returns the underlying value. It's meant to be type asserted against
	// Example: myInt := v.(int)
	Get() interface{}

	// Update appends to container type Values from a string (useful for CLI flags, env vars, default values)
	// and returns ErrCantUpdateScalarType when trying to Update non-empty scalar Values
	Update(string) error

	// UpdateFrom interface updates a container type Value from an interface (useful for configs)
	// It returns ErrCantUpdateScalarType when trying to Update non-empty scalar Values
	// Note that UpdateFromInterface must be called with the "contained" type for container type Values
	// For example, the StringSlice.UpdateFromInterface
	// must be called with a string, not a []string
	// It returns ErrIncompatibleInterface if the interface can't be decoded
	UpdateFromInterface(interface{}) error
}

// EmptyConstructur just builds a new value
// Useful to create new values as well as initialize them
type EmptyConstructor = func() (Value, error)

// FromInterface specifies how to create a Value from an interface
// Useful for reading a value from a config
// It returns ErrIncompatibleInterface if the interface can't be decoded
type FromInterface = func(interface{}) (Value, error)

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")
var ErrCantUpdateScalarType = errors.New("scalar types can not be updated")
