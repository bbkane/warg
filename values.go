package clide

import "strconv"

// Value is a "generic" type that lets me store different types into flags
//  ~Stolen from~ "Inspired" by https://golang.org/src/flag/flag.go?#L138
type Value interface {
	// Get returns the underlying value. It's meant to be type asserted against
	Get() interface{}

	// Update updates the underlying value from a string
	// It replaces single values and appends to list values
	Update(string) error
}

// IntValue
type IntValue int

func NewIntValue(val int) *IntValue {
	return (*IntValue)(&val)
}

func (i *IntValue) Get() interface{} {
	return int(*i)
}

func (i *IntValue) Update(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*i = IntValue(v)
	return nil
}

type StringSliceValue []string

func NewStringSliceValue(vals []string) *StringSliceValue {
	return (*StringSliceValue)(&vals)
}

func (ss *StringSliceValue) Get() interface{} {
	return []string(*ss)
}

func (ss *StringSliceValue) Update(val string) error {
	*ss = append(*ss, val)
	return nil
}
