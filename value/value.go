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

type Int int

func IntNew(val int) *Int { return (*Int)(&val) }
func IntEmpty() *Int      { return IntNew(0) }
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

type String string

func StringNew(val string) *String { return (*String)(&val) }
func StringEmpty() *String         { return StringNew("") }
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

type StringValue []string

func StringSliceNew(vals []string) *StringValue { return (*StringValue)(&vals) }
func StringSliceEmpty() *StringValue            { return StringSliceNew(nil) }
func (ss *StringValue) Get() interface{}        { return []string(*ss) }
func (ss *StringValue) String() string          { return fmt.Sprint([]string(*ss)) }
func (ss *StringValue) Update(val string) error {
	*ss = append(*ss, val)
	return nil
}
