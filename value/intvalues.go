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

func (v *intV) Get() interface{}    { return int(*v) }
func (v *intV) String() string      { return fmt.Sprint(int(*v)) }
func (v *intV) TypeInfo() typeInfo  { return TypeInfoScalar }
func (v *intV) Description() string { return "int" }

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

// intSlice is updateable from a float or int. If a float is passed, it will be truncated.
// Example: 4.5 -> 4, 3.99 -> 3
type intSlice []int

func intSliceNew(vals []int) *intSlice {
	return (*intSlice)(&vals)
}

func IntSlice() (Value, error)       { return intSliceNew(nil), nil }
func (v *intSlice) Get() interface{} { return []int(*v) }

func (v *intSlice) ReplaceFromInterface(iFace interface{}) error {
	switch under := iFace.(type) {
	case []int:
		*v = *intSliceNew(under)
	case []float64:
		var ret []int
		for _, e := range under {
			ret = append(ret, int(e))
		}
		*v = *intSliceNew(ret)
	default:
		return ErrIncompatibleInterface
	}
	return nil
}
func (v *intSlice) TypeInfo() typeInfo  { return TypeInfoSlice }
func (v *intSlice) Description() string { return "int slice" }

func (v *intSlice) String() string { return fmt.Sprint([]int(*v)) }
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
