package value

import (
	"fmt"
	"strconv"
)

// FromInterface specifies how to create a Value from an interface
// Useful for reading a value from a config
type FromInterface = func(interface{}) (Value, error)

type ValueMap = map[string]Value

// Value is a "generic" type that lets me store different types into flags
//  ~Stolen from~ "Inspired by" https://golang.org/src/flag/flag.go?#L138
type Value interface {
	// Get returns the underlying value. It's meant to be type asserted against
	Get() interface{}

	// Update updates the underlying value from a string
	// It replaces single values and appends to list values
	Update(string) error

	// Make it printable!
	String() string
}

// TODO: should I be returning pointers?

type IntValue int

func IntValueNew(val int) *IntValue { return (*IntValue)(&val) }
func IntValueEmpty() *IntValue      { return IntValueNew(0) }
func IntValueFromInterface(val interface{}) (Value, error) {
	under, ok := val.(int)
	if !ok {
		return nil, fmt.Errorf("can't create IntValue. Expected: int, got: %#v", val)
	}
	return IntValueNew(under), nil
}
func (i *IntValue) Get() interface{} { return int(*i) }
func (i *IntValue) String() string   { return fmt.Sprint(int(*i)) }

func (i *IntValue) Update(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*i = IntValue(v)
	return nil
}

type StringValue string

func StringValueNew(val string) *StringValue { return (*StringValue)(&val) }
func StringValueEmpty() *StringValue         { return StringValueNew("") }
func StringValueFromInterface(val interface{}) (Value, error) {
	under, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("can't create StringValue. Expected: string, got: %#v", val)
	}
	return StringValueNew(under), nil
}
func (v *StringValue) Get() interface{} { return string(*v) }
func (v *StringValue) String() string   { return fmt.Sprint(string(*v)) }
func (v *StringValue) Update(s string) error {
	*v = StringValue(s)
	return nil
}

type StringSliceValue []string

func StringSliceValueNew(vals []string) *StringSliceValue { return (*StringSliceValue)(&vals) }
func StringSliceValueEmpty() *StringSliceValue            { return StringSliceValueNew(nil) }
func (ss *StringSliceValue) Get() interface{}             { return []string(*ss) }
func (ss *StringSliceValue) String() string               { return fmt.Sprint([]string(*ss)) }
func (ss *StringSliceValue) Update(val string) error {
	*ss = append(*ss, val)
	return nil
}
