package value

import (
	"fmt"
)

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

// String accepts a string from a user. Pretty self explanatory.
func String() (Value, error) { return stringNew(""), nil }

func (v *stringV) ReplaceFromInterface(iFace interface{}) error {
	return v.UpdateFromInterface(iFace)
}

// --

type stringEnumV struct {
	description string
	choices     []string
	current     string
}

func (v *stringEnumV) Get() interface{}    { return v.current }
func (v *stringEnumV) String() string      { return v.current }
func (v *stringEnumV) TypeInfo() typeInfo  { return TypeInfoScalar }
func (v *stringEnumV) Description() string { return v.description }
func (v *stringEnumV) Update(val string) error {
	var updated bool
	for _, choice := range v.choices {
		if val == choice {
			v.current = choice
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("string enum update invalid choice: available: %v: choice: %v", v.choices, val)
	}
	return nil
}

func (v *stringEnumV) UpdateFromInterface(iFace interface{}) error {
	under, ok := iFace.(string)
	if !ok {
		return ErrIncompatibleInterface
	}
	return v.Update(under)
}

func (v *stringEnumV) ReplaceFromInterface(iFace interface{}) error {
	return v.UpdateFromInterface(iFace)
}

// StringEnum acts just like a string, except it only lets the user update from
// the choices provided when creating the EmptyConstructor for it.
func StringEnum(choices ...string) EmptyConstructor {
	return func() (Value, error) {
		return &stringEnumV{
			choices:     choices,
			description: fmt.Sprintf("stringenum with choices: %v", choices),
		}, nil
	}
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

// StringSlice accepts a string from a user and adds it to a slice. Pretty self explanatory.
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
