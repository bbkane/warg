package value

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
)

type pathV string

func pathNew(val string) (*pathV, error) {
	expanded, err := homedir.Expand(val)
	if err != nil {
		return nil, fmt.Errorf("could not expand homedir for %v: err: %v", val, err)
	}
	return (*pathV)(&expanded), nil
}
func (v *pathV) Get() interface{}    { return string(*v) }
func (v *pathV) String() string      { return fmt.Sprint(string(*v)) }
func (v *pathV) TypeInfo() typeInfo  { return TypeInfoScalar }
func (v *pathV) Description() string { return "path" }

func (v *pathV) Update(s string) error {
	new, err := pathNew(s)
	if err != nil {
		return fmt.Errorf("could not expand homedir for %v: err: %v", new, err)
	}
	*v = *new
	return nil
}
func (v *pathV) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	new, err := pathNew(under)
	if err != nil {
		return fmt.Errorf("could not expand homedir for %v: err: %v", new, err)
	}
	*v = *new
	return nil
}

// Path autoexpands ~ when updated and otherwise behaves like a string
func Path() (Value, error) {
	return pathNew("")
}

func (v *pathV) ReplaceFromInterface(iFace interface{}) error {
	return v.UpdateFromInterface(iFace)
}

// ---

type pathSliceV []string

func (v *pathSliceV) Get() interface{} { return []string(*v) }
func (v *pathSliceV) String() string   { return fmt.Sprint([]string(*v)) }
func (v *pathSliceV) Update(val string) error {
	expanded, err := homedir.Expand(val)
	if err != nil {
		return fmt.Errorf("could not expand homedir for %v: err: %v", val, err)
	}
	*v = append(*v, expanded)
	return nil
}
func (v *pathSliceV) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	expanded, err := homedir.Expand(under)
	if err != nil {
		return fmt.Errorf("could not expand homedir for %v: err: %v", under, err)
	}
	*v = append(*v, expanded)
	return nil
}

func (v *pathSliceV) TypeInfo() typeInfo  { return TypeInfoSlice }
func (v *pathSliceV) Description() string { return "path slice" }

// PathSlice autoexpands ~ when updated and otherwise behaves like a []string
func PathSlice() (Value, error) { return &pathSliceV{}, nil }

func (v *pathSliceV) ReplaceFromInterface(iFace interface{}) error {
	under, ok := iFace.([]string)
	if !ok {
		return ErrIncompatibleInterface
	}
	new, _ := PathSlice()
	for _, e := range under {
		err := new.Update(e)
		if err != nil {
			return fmt.Errorf("could not expand: %v: %w", e, err)
		}
	}
	// TODO: does this even work?
	new2 := new.Get().([]string)
	*v = *(*pathSliceV)(&new2)
	return nil
}
