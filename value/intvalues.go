package value

import (
	"fmt"
	"strconv"
)

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
