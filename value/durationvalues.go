package value

import (
	"fmt"
	"time"
)

type durationV time.Duration

func durationNew(val time.Duration) *durationV {
	return (*durationV)(&val)
}

// Duration is updateable from a string parsed with https://pkg.go.dev/time#ParseDuration
func Duration() (Value, error) { return durationNew(time.Duration(0)), nil }

func (v *durationV) Get() interface{}      { return time.Duration(*v) }
func (v *durationV) String() string        { return fmt.Sprint(time.Duration(*v)) }
func (v *durationV) StringSlice() []string { return nil }
func (v *durationV) TypeInfo() TypeInfo    { return TypeInfoScalar }
func (v *durationV) Description() string   { return "duration" }

func (v *durationV) ReplaceFromInterface(iFace interface{}) error {
	switch under := iFace.(type) {
	case string:
		return v.Update(under)
	default:
		return ErrIncompatibleInterface
	}
}

func (v *durationV) Update(s string) error {
	decoded, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*v = durationV(decoded)
	return nil
}
func (v *durationV) UpdateFromInterface(iFace interface{}) error {
	return v.ReplaceFromInterface(iFace)
}
