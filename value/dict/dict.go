package dict

import (
	"fmt"
	"strings"

	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
)

type dictValue[T comparable] struct {
	choices     []T
	defaultVals map[string]T
	hasDefault  bool
	inner       contained.TypeInfo[T]
	vals        map[string]T
	updatedBy   value.UpdatedBy
}

type DictOpt[T comparable] func(*dictValue[T])

func New[T comparable](inner contained.TypeInfo[T], opts ...DictOpt[T]) value.EmptyConstructor {
	return func() value.Value {
		dv := dictValue[T]{
			choices:     []T{},
			defaultVals: make(map[string]T),
			hasDefault:  false,
			inner:       inner,
			vals:        make(map[string]T),
			updatedBy:   value.UpdatedByUnset,
		}
		for _, opt := range opts {
			opt(&dv)
		}
		return &dv
	}
}

func Choices[T comparable](choices ...T) DictOpt[T] {
	return func(v *dictValue[T]) {
		v.choices = choices
	}
}

func Default[T comparable](def map[string]T) DictOpt[T] {
	return func(cf *dictValue[T]) {
		cf.defaultVals = def
		cf.hasDefault = true
	}
}

func (v *dictValue[_]) Choices() []string {
	ret := []string{}
	for _, e := range v.choices {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}

func (v *dictValue[_]) DefaultStringMap() map[string]string {
	ret := make(map[string]string, len(v.defaultVals))
	for k, e := range v.defaultVals {
		ret[k] = fmt.Sprint(e)
	}
	return ret
}

func (v *dictValue[_]) Description() string {
	return v.inner.Description
}

func (v *dictValue[_]) Get() interface{} {
	return v.vals
}

func (v *dictValue[_]) HasDefault() bool {
	return v.hasDefault
}

func (v *dictValue[T]) ReplaceFromInterface(iFace interface{}, u value.UpdatedBy) error {
	under, ok := iFace.(map[string]interface{})
	if !ok {
		return contained.ErrIncompatibleInterface // TODO: should ErrIncompatibleInterface be in value?
	}

	newVals := make(map[string]T)
	for k, e := range under {
		underE, err := v.inner.FromIFace(e)
		if err != nil {
			// TODO: this won't communicate to the caller *which* element is the wrong type
			return err
		}
		newVals[k] = underE
	}
	v.vals = newVals
	v.updatedBy = u
	return nil
}

func (v *dictValue[_]) StringMap() map[string]string {
	ret := make(map[string]string, len(v.vals))
	for k, e := range v.vals {
		ret[k] = fmt.Sprint(e)
	}
	return ret
}

func withinChoices[T comparable](val T, choices []T) bool {
	// User didn't constrain choices
	if len(choices) == 0 {
		return true
	}
	for _, choice := range choices {
		if val == choice {
			return true
		}
	}
	return false
}

func (v *dictValue[T]) update(key string, val T) error {
	if !withinChoices(val, v.choices) {
		return value.ErrInvalidChoice[T]{Choices: v.choices}
	}
	v.vals[key] = val
	return nil
}

func (v *dictValue[_]) Update(s string, u value.UpdatedBy) error {
	key, strValue, found := strings.Cut(s, "=")
	if !found {
		return fmt.Errorf("Could not parse key=value for %v", s)
	}
	val, err := v.inner.FromString(strValue)
	if err != nil {
		return err
	}
	err = v.update(key, val)
	if err != nil {
		return err
	}
	v.updatedBy = u
	return nil
}

func (v *dictValue[_]) ReplaceFromDefault(u value.UpdatedBy) {
	if v.hasDefault {
		v.vals = v.defaultVals
		v.updatedBy = u
	}
}
