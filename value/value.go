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
	// It replaces single values and appends to aggregate values,
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
func IntFromInterface(val interface{}) (Value, error) {
	under, ok := val.(int)
	if !ok {
		return nil, fmt.Errorf("can't create IntValue. Expected: int, got: %#v", val)
	}
	return IntNew(under), nil
}
func (i *Int) Get() interface{} { return int(*i) }
func (i *Int) String() string   { return fmt.Sprint(int(*i)) }

func (i *Int) Update(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*i = Int(v)
	return nil
}
func (i *Int) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(int)
	if !ok {
		return ErrIncompatibleInterface
	}
	*i = Int(under)
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
func StringSliceEmpty() Value            { return StringSliceNew(nil) }
func (ss *StringSlice) Get() interface{} { return []string(*ss) }
func (ss *StringSlice) String() string   { return fmt.Sprint([]string(*ss)) }
func (ss *StringSlice) Update(val string) error {
	*ss = append(*ss, val)
	return nil
}
func (ss *StringSlice) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	*ss = append(*ss, under)
	return nil
}
