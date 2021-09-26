package value

import (
	"errors"
	"fmt"
	"strconv"
)

// FromInterface specifies how to create a Value from an interface
// Useful for reading a value from a config
type FromInterface = func(interface{}) (Value, error)

// EmptyConstructur just builds a new value
// Useful to create new values as well as initialize them
// TODO: better name :)
type EmptyConstructor = func() Value

// Value is a "generic" type that lets me store different types into flags
//  ~Stolen from~ "Inspired by" https://golang.org/src/flag/flag.go?#L138
type Value interface {
	// Get returns the underlying value. It's meant to be type asserted against
	Get() interface{}

	// Update updates the underlying value from a string
	// It replaces single values and appends to list values
	// TODO: throw an error for single values
	Update(string) error

	// UpdateFromInterface updates the underlying value from an interface
	// It replaces single values and appends to container type values,
	// so the interface MUST BE the 'single' part of the aggreate
	// For exampple, the StringSlice.UpdateFromInterface
	// must be called with a string, not a []string
	// TODO: return an error for already initialized single values
	// This function is needed for configpath handling
	UpdateFromInterface(interface{}) error

	// Make it printable!
	String() string
}

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")

type Int int

func IntNew(val int) *Int { return (*Int)(&val) }
func IntEmpty() Value     { return IntNew(0) }
func IntFromFloatOrIntInterface(val interface{}) (Value, error) {
	switch under := val.(type) {
	case int:
		return IntNew(under), nil
	case float64: // like JSON
		return IntNew(int(under)), nil
	default:
		return nil, fmt.Errorf("can't create IntValue. Expected: int or float64, got: %#v", val)
	}

}
func IntFromInterface(val interface{}) (Value, error) {
	under, ok := val.(int)
	if !ok {
		return nil, fmt.Errorf("can't create IntValue. Expected: int, got: %#v", val)
	}
	return IntNew(under), nil
}
func (v *Int) Get() interface{} { return int(*v) }
func (v *Int) String() string   { return fmt.Sprint(int(*v)) }

func (v *Int) Update(s string) error {
	decoded, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = Int(decoded)
	return nil
}
func (v *Int) UpdateFromInterface(iFace interface{}) error {
	// TODO: make this accept a float to not panic!
	under, ok := iFace.(int)
	if !ok {
		return ErrIncompatibleInterface
	}
	*v = Int(under)
	return nil
}

type String string

func StringNew(val string) *String { return (*String)(&val) }
func StringEmpty() Value           { return StringNew("") }
func StringFromInterface(val interface{}) (Value, error) {
	under, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("can't create StringValue. Expected: string, got: %#v", val)
	}
	return StringNew(under), nil
}
func (v *String) Get() interface{} { return string(*v) }
func (v *String) String() string   { return fmt.Sprint(string(*v)) }
func (v *String) Update(s string) error {
	*v = String(s)
	return nil
}
func (v *String) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	*v = String(under)
	return nil
}

type StringSlice []string

func StringSliceNew(vals []string) *StringSlice { return (*StringSlice)(&vals) }
func StringSliceFromInterface(val interface{}) (Value, error) {
	under, ok := val.([]string)
	if !ok {
		return nil, fmt.Errorf("can't create StringSlice. Expected: []string, got: %#v", val)
	}
	return StringSliceNew(under), nil
}
func StringSliceEmpty() Value           { return StringSliceNew(nil) }
func (v *StringSlice) Get() interface{} { return []string(*v) }
func (v *StringSlice) String() string   { return fmt.Sprint([]string(*v)) }
func (v *StringSlice) Update(val string) error {
	*v = append(*v, val)
	return nil
}
func (v *StringSlice) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	*v = append(*v, under)
	return nil
}

type IntSlice []int

func IntSliceNew(vals []int) *IntSlice {
	return (*IntSlice)(&vals)
}
func IntSliceFromInterface(val interface{}) (Value, error) {
	under, ok := val.([]int)
	if !ok {
		return nil, ErrIncompatibleInterface
	}
	return IntSliceNew(under), nil
}
func IntSliceEmpty() Value           { return IntSliceNew(nil) }
func (v *IntSlice) Get() interface{} { return []int(*v) }
func (v *IntSlice) String() string   { return fmt.Sprint([]int(*v)) }
func (v *IntSlice) Update(s string) error {
	decoded, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = append(*v, int(decoded))
	return nil
}
func (v *IntSlice) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(int)
	if !ok {
		return ErrIncompatibleInterface
	}
	*v = append(*v, under)
	return nil
}
