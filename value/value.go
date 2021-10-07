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
type EmptyConstructor = func() (Value, error)

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
	// and returns ErrCantUpdateScalarType when trying to Update non-empty scalar Values (TODO)
	Update(string) error

	// UpdateFrom interface updates a container type Value from an interface (useful for configs)
	// and returns ErrCantUpdateScalarType when trying to Update non-empty scalar Values (TODO)
	// Note that UpdateFromInterface must be called with the "contained" type for container type Values
	// For example, the StringSlice.UpdateFromInterface
	// must be called with a string, not a []string
	UpdateFromInterface(interface{}) error
}

var ErrIncompatibleInterface = errors.New("could not decode interface into Value")
var ErrCantUpdateScalarType = errors.New("scalar types can not be updated")

// intV is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
type intV int

func intNew(val int) *intV     { return (*intV)(&val) }
func IntEmpty() (Value, error) { return intNew(0), nil }
func IntFromInterface(iFace interface{}) (Value, error) {
	switch under := iFace.(type) {
	case int:
		return intNew(under), nil
	case float64: // like JSON
		return intNew(int(under)), nil
	default:
		return nil, fmt.Errorf("can't create IntValue. Expected: int or float64, got: %#v", iFace)
	}
}
func (v *intV) Get() interface{} { return int(*v) }
func (v *intV) String() string   { return fmt.Sprint(int(*v)) }

func (v *intV) Update(s string) error {
	decoded, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = intV(decoded)
	return nil
}
func (v *intV) UpdateFromInterface(iFace interface{}) error {
	// TODO: make this accept a float to not panic!
	switch under := iFace.(type) {
	case int:
		*v = intV(under)
	case float64: // like JSON
		*v = intV(int(under))
	default:
		return fmt.Errorf("can't create IntValue. Expected: int or float64, got: %#v", iFace)
	}
	return nil
}

// type path struct {
// 	val     string
// 	updated bool
// }

// func pathNew(val string) (*path, error) {
// 	expanded, err := homedir.Expand(val)
// 	if err != nil {
// 		return fmt.Errorf("Could not expand homedir for %v: err: %v", val, err)
// 	}
// 	return &path{
// 		val:     expanded,
// 		updated: true,
// 	}
// }

// func PathEmpty() Value {
// 	return pathNew("")
// }

// func PathFromInterface(iFace interface{}) (Value, error) {
// 	under, ok := iFace.(string)
// 	if !ok {
// 		return nil, ErrIncompatibleInterface
// 	}
// 	return pathNew(under)
// }

type stringV string

func stringNew(val string) *stringV { return (*stringV)(&val) }
func StringEmpty() (Value, error)   { return stringNew(""), nil }
func StringFromInterface(iFace interface{}) (Value, error) {
	under, ok := iFace.(string)
	if !ok {
		return nil, fmt.Errorf("can't create StringValue. Expected: string, got: %#v", iFace)
	}
	return stringNew(under), nil
}
func (v *stringV) Get() interface{} { return string(*v) }
func (v *stringV) String() string   { return fmt.Sprint(string(*v)) }
func (v *stringV) Update(s string) error {
	*v = stringV(s)
	return nil
}
func (v *stringV) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	*v = stringV(under)
	return nil
}

type stringSliceV []string

func stringSliceNew(vals []string) *stringSliceV { return (*stringSliceV)(&vals) }
func StringSliceFromInterface(iFace interface{}) (Value, error) {
	under, ok := iFace.([]string)
	if !ok {
		return nil, fmt.Errorf("can't create StringSlice. Expected: []string, got: %#v", iFace)
	}
	return stringSliceNew(under), nil
}
func StringSliceEmpty() (Value, error)   { return stringSliceNew(nil), nil }
func (v *stringSliceV) Get() interface{} { return []string(*v) }
func (v *stringSliceV) String() string   { return fmt.Sprint([]string(*v)) }
func (v *stringSliceV) Update(val string) error {
	*v = append(*v, val)
	return nil
}
func (v *stringSliceV) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	*v = append(*v, under)
	return nil
}

// intSlice is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
type intSlice []int

func intSliceNew(vals []int) *intSlice {
	return (*intSlice)(&vals)
}
func IntSliceFromInterface(iFace interface{}) (Value, error) {

	switch under := iFace.(type) {
	case []int:
		return intSliceNew(under), nil
	case []float64:
		var ret []int
		for _, e := range under {
			ret = append(ret, int(e))
		}
		return intSliceNew(ret), nil
	default:
		return nil, ErrIncompatibleInterface
	}
}
func IntSliceEmpty() (Value, error)  { return intSliceNew(nil), nil }
func (v *intSlice) Get() interface{} { return []int(*v) }
func (v *intSlice) String() string   { return fmt.Sprint([]int(*v)) }
func (v *intSlice) Update(s string) error {
	decoded, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = append(*v, int(decoded))
	return nil
}
func (v *intSlice) UpdateFromInterface(iFace interface{}) error {
	switch under := iFace.(type) {
	case int:
		*v = append(*v, under)
	case float64: // like JSON
		*v = append(*v, int(under))
	default:
		return fmt.Errorf("can't update IntSlice. Expected: int or float64, got: %#v", iFace)
	}
	return nil
}
