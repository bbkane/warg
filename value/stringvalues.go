package value

import "fmt"

type stringV string

func (v *stringV) Get() interface{}    { return string(*v) }
func (v *stringV) String() string      { return fmt.Sprint(string(*v)) }
func (v *stringV) TypeInfo() typeInfo  { return TypeInfoScalar }
func (v *stringV) Description() string { return "string" }

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
func String() (Value, error)        { return stringNew(""), nil }

func (v *stringV) ReplaceFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return fmt.Errorf("can't create StringValue. Expected: string, got: %#v", iFace)
	}
	*v = *stringNew(under)
	return nil
}

// ---

type stringSliceV []string

func stringSliceNew(vals []string) *stringSliceV { return (*stringSliceV)(&vals) }

func (v *stringSliceV) ReplaceFromInterface(iFace interface{}) error {
	under, ok := iFace.([]string)
	if !ok {
		return ErrIncompatibleInterface
	}
	*v = *stringSliceNew(under)
	return nil
}
func (v *stringSliceV) TypeInfo() typeInfo  { return TypeInfoSlice }
func (v *stringSliceV) Description() string { return "string slice" }

func StringSlice() (Value, error)        { return stringSliceNew(nil), nil }
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
