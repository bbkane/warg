package value

import (
	"fmt"
	"strconv"
)

// intV is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
type intV int

func intNew(val int) *intV { return (*intV)(&val) }

// Int is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
func Int() (Value, error) { return intNew(0), nil }

func (v *intV) Get() interface{}      { return int(*v) }
func (v *intV) String() string        { return fmt.Sprint(int(*v)) }
func (v *intV) StringSlice() []string { return nil }
func (v *intV) TypeInfo() TypeInfo    { return TypeInfoScalar }
func (v *intV) Description() string   { return "int" }

func (v *intV) ReplaceFromInterface(iFace interface{}) error {
	switch under := iFace.(type) {
	case int:
		*v = *intNew(under)
	case float64: // like JSON
		*v = *intNew(int(under))
	default:
		return fmt.Errorf("can't create IntValue. Expected: int or float64, got: %#v", iFace)
	}
	return nil
}

func (v *intV) Update(s string) error {
	decoded, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = intV(decoded)
	return nil
}
func (v *intV) UpdateFromInterface(iFace interface{}) error {
	return v.ReplaceFromInterface(iFace)
}

// intSliceV is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
type intSliceV []int

func intSliceNew(vals []int) *intSliceV {
	return (*intSliceV)(&vals)
}

// IntSlice is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
func IntSlice() (Value, error)        { return intSliceNew(nil), nil }
func (v *intSliceV) Get() interface{} { return []int(*v) }

func (v *intSliceV) ReplaceFromInterface(iFace interface{}) error {
	// NOTE: this is the most up to date ReplaceFromInterface :)
	// it uses UpdateFromInterface
	under, ok := iFace.([]interface{})
	if !ok {
		return ErrIncompatibleInterface
	}

	new, _ := IntSlice()
	for _, e := range under {
		err := new.UpdateFromInterface(e)
		if err != nil {
			return err
		}
	}
	new2 := new.Get().([]int)
	*v = *(*intSliceV)(&new2)
	return nil
}
func (v *intSliceV) TypeInfo() TypeInfo  { return TypeInfoSlice }
func (v *intSliceV) Description() string { return "int slice" }

func (v *intSliceV) String() string { return fmt.Sprint([]int(*v)) }
func (v *intSliceV) StringSlice() []string {
	var ret []string
	for _, e := range []int(*v) {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}

func (v *intSliceV) Update(s string) error {
	decoded, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = append(*v, int(decoded))
	return nil
}
func (v *intSliceV) UpdateFromInterface(iFace interface{}) error {
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
