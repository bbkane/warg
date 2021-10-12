package value

import "fmt"

type stringV string

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

func stringNew(val string) *stringV { return (*stringV)(&val) }
func StringEmpty() (Value, error)   { return stringNew(""), nil }
func StringFromInterface(iFace interface{}) (Value, error) {
	under, ok := iFace.(string)
	if !ok {
		return nil, fmt.Errorf("can't create StringValue. Expected: string, got: %#v", iFace)
	}
	return stringNew(under), nil
}

// ---

type stringSliceV []string

func stringSliceNew(vals []string) *stringSliceV { return (*stringSliceV)(&vals) }
func StringSliceFromInterface(iFace interface{}) (Value, error) {
	under, ok := iFace.([]string)
	if !ok {
		return nil, ErrIncompatibleInterface
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
