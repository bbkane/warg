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
func (v *pathV) Get() interface{} { return string(*v) }
func (v *pathV) String() string   { return fmt.Sprint(string(*v)) }
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

func PathEmpty() (Value, error) {
	return pathNew("")
}

func PathFromInterface(iFace interface{}) (Value, error) {
	under, ok := iFace.(string)
	if !ok {
		return nil, ErrIncompatibleInterface
	}
	return pathNew(under)
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
func PathSliceEmpty() (Value, error) { return &pathSliceV{}, nil }
func PathSliceFromInterface(iFace interface{}) (Value, error) {
	under, ok := iFace.([]string)
	if !ok {
		return nil, ErrIncompatibleInterface
	}
	new, _ := PathSliceEmpty()
	for _, e := range under {
		err := new.Update(e)
		if err != nil {
			return nil, fmt.Errorf("could not expand: %v: %w", e, err)
		}
	}
	return new, nil
}
