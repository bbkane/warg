package value

import (
	"fmt"
)

type boolV bool

func boolNew(val bool) *boolV { return (*boolV)(&val) }

// Bool is updated from "true" or "false"
func Bool() (Value, error) { return boolNew(false), nil }

func (v *boolV) Get() interface{}      { return bool(*v) }
func (v *boolV) String() string        { return fmt.Sprint(bool(*v)) }
func (v *boolV) StringSlice() []string { return nil }
func (v *boolV) TypeInfo() TypeInfo    { return TypeInfoScalar }
func (v *boolV) Description() string   { return "int" }

func (v *boolV) ReplaceFromInterface(iFace interface{}) error {
	switch under := iFace.(type) {
	case bool:
		*v = *boolNew(under)
	default:
		return fmt.Errorf("can't create boolValue. Expected: bool, got: %#v", iFace)
	}
	return nil
}

func (v *boolV) Update(s string) error {

	switch s {
	case "true":
		*v = boolV(true)
	case "false":
		*v = boolV(false)
	default:
		return fmt.Errorf("expected \"true\" or \"false\", got %s", s)
	}
	return nil
}
func (v *boolV) UpdateFromInterface(iFace interface{}) error {
	return v.ReplaceFromInterface(iFace)
}
