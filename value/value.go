package value

import (
	"fmt"
	"strings"

	"go.bbkane.com/warg/colerr"
	"go.bbkane.com/warg/styles"
)

// UpdatedBy identifies the source that last set a flag's value.
type UpdatedBy string

const (
	UpdatedByUnset   UpdatedBy = ""
	UpdatedByDefault UpdatedBy = "appdefault"
	UpdatedByEnvVar  UpdatedBy = "envvar"
	UpdatedByFlag    UpdatedBy = "passedflag"
	UpdatedByConfig  UpdatedBy = "config"
)

// Value is the interface for all flag value types (scalar, slice, dict).
// Implementations hold the current value, track how it was set, and handle
// parsing from strings (CLI/env) and interfaces (config files).
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

// ScalarValue extends [Value] for single-valued flags (e.g., string, int, bool).
type ScalarValue interface {
	Value

	// DefaultString returns the default underlying value (represented as a string)
	DefaultString() string

	// String returns a string ready to be printed!
	String() string
}

// SliceValue extends [Value] for list-valued flags that accumulate multiple values.
type SliceValue interface {
	Value

	// DefaultStringSlice returns the default underlying value for slice values and nil for others
	DefaultStringSlice() []string

	// StringSlice returns a []string ready to be printed for slice values and nil for others
	StringSlice() []string
}

// DictValue extends [Value] for key=value map flags.
type DictValue interface {
	Value

	// DefaultStringMap returns the default underlying value for the map
	DefaultStringMap() map[string]string

	// StringMap returns what's in the map currently
	StringMap() map[string]string
}

// EmptyConstructor creates a new zero-valued [Value] instance.
// Used both for initialization and to produce fresh values during parsing.
type EmptyConstructor func() Value

// ErrInvalidChoice is returned when an Update value is not within the allowed choices.
type ErrInvalidChoice[T any] struct {
	Choices []T
}

func (e ErrInvalidChoice[T]) Error() string {
	return "invalid choice for value: choices: " + fmt.Sprint(e.Choices)
}

func (e ErrInvalidChoice[T]) ColorError(s *styles.Styles) string {
	var buf strings.Builder

	buf.WriteString("invalid choice for value: choices:\n")
	for _, c := range e.Choices {
		buf.WriteString(string(s.ErrorAltCode) + "  " + fmt.Sprint(c) + "\n")
	}

	return s.Error(buf.String())
}

// ErrUpdatedMoreThanOnce is returned when a scalar value is set more than once
// from the same priority level (e.g., two CLI flags for the same scalar).
type ErrUpdatedMoreThanOnce[T any] struct {
	CurrentValue T
	UpdatedBy    UpdatedBy
}

func (e ErrUpdatedMoreThanOnce[T]) Error() string {
	return fmt.Sprintf("value already updated to %#v by %s", e.CurrentValue, e.UpdatedBy)
}

func (e ErrUpdatedMoreThanOnce[T]) ColorError(s *styles.Styles) string {
	err := colerr.NewWrappedf(nil, "Value already updated to %s by %v", fmt.Sprintf("%v", e.CurrentValue), string(e.UpdatedBy))
	return err.ColorError(s)
}
