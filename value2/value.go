package value

import (
	"errors"
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
var ErrInvalidChoice = errors.New("invalid choice for value")
