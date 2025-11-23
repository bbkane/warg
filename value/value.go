package value

import (
	"fmt"
)

type UpdatedBy string

const (
	UpdatedByUnset   UpdatedBy = ""
	UpdatedByDefault UpdatedBy = "appdefault"
	UpdatedByEnvVar  UpdatedBy = "envvar"
	UpdatedByFlag    UpdatedBy = "passedflag"
	UpdatedByConfig  UpdatedBy = "config"
)

// Value is a "generic" type to store different types into flags
// Inspired by https://golang.org/src/flag/flag.go .
// There are two underlying "type" families designed to fit in Value:
// scalar types (Int, String, ...) and container types (IntSlice, StringMap, ...).
type Value interface {

	// Choices for this value to contain (represented as strings)
	Choices() []string

	// Description of the type. Useful for help messages. Not guaranteed to be unique. Should not be used as an ID.
	Description() string

	// Get returns the underlying value. It's meant to be type asserted against
	// Example: myInt := v.(int)
	Get() interface{}

	// HasDefault returns true if this value has a default
	HasDefault() bool

	// ReplaceFromInterface replaces a value with one found in an interface.
	// Useful to update a Value from a config.
	ReplaceFromInterface(interface{}, UpdatedBy) error

	// Update appends to container type Values from a string (useful for CLI flags, env vars, default values)
	// and replaces scalar Values
	Update(string, UpdatedBy) error

	// UpdatedBy returns the source of the last update
	UpdatedBy() UpdatedBy

	// ReplaceFromDefault updates the Value from a pre-set default, if one exists. use HasDefault to check whether a default exists
	ReplaceFromDefault(u UpdatedBy) error
}

type ScalarValue interface {
	Value

	// DefaultString returns the default underlying value (represented as a string)
	DefaultString() string

	// String returns a string ready to be printed!
	String() string
}

type SliceValue interface {
	Value

	// DefaultStringSlice returns the default underlying value for slice values and nil for others
	DefaultStringSlice() []string

	// StringSlice returns a []string ready to be printed for slice values and nil for others
	StringSlice() []string

	// AppendFromInterface updates a container type Value from an interface (useful for configs)
	// Note that AppendFromInterface must be called with the "contained" type for container type Values
	// For example, the StringSlice.AppendFromInterface
	// must be called with a string, not a []string
	// It returns ErrIncompatibleInterface if the interface can't be decoded
	AppendFromInterface(interface{}, UpdatedBy) error
}

type DictValue interface {
	Value

	// DefaultStringMap returns the default underlying value for the map
	DefaultStringMap() map[string]string

	// StringMap returns what's in the map currently
	StringMap() map[string]string
}

// EmptyConstructur just builds a new value.
// Useful to create new values as well as initialize them
type EmptyConstructor func() Value

type ErrInvalidChoice[T any] struct {
	Choices []T
}

func (e ErrInvalidChoice[T]) Error() string {
	return "invalid choice for value: choices: " + fmt.Sprint(e.Choices)
}

// ErrUpdatedMoreThanOnce is returned when a value is updated more than once. Only applicable to Scalar Values
type ErrUpdatedMoreThanOnce[T any] struct {
	CurrentValue T
	UpdatedBy    UpdatedBy
}

func (e ErrUpdatedMoreThanOnce[T]) Error() string {
	return fmt.Sprintf("value already updated to %#v by %s", e.CurrentValue, e.UpdatedBy)
}
