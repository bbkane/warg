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
}

type DictOpt[T comparable] func(*dictValue[T])

// TODO: TEST THIS, but lets try to get more sleep

func New[T comparable](inner contained.TypeInfo[T], opts ...DictOpt[T]) value.EmptyConstructor {
	return func() (value.Value, error) {
		dv := dictValue[T]{
			choices:     []T{},
			defaultVals: make(map[string]T),
			hasDefault:  false,
			inner:       inner,
			vals:        make(map[string]T),
		}
		for _, opt := range opts {
			opt(&dv)
		}
		return &dv, nil
	}
}

func Choices[T comparable](choices ...T) DictOpt[T] {
	return func(v *dictValue[T]) {
		for _, ch := range choices {
			newChoice, err := v.inner.FromInstance(ch)
			if err != nil {
				panic("error constructing default value: " + fmt.Sprint(ch) + " : " + err.Error())
			}
			v.choices = append(v.choices, newChoice)
		}
	}
}

func Default[T comparable](def map[string]T) DictOpt[T] {
	return func(cf *dictValue[T]) {
		for key, d := range def {
			newD, err := cf.inner.FromInstance(d)
			if err != nil {
				panic("error constructing default value: " + fmt.Sprint(d) + " : " + err.Error())
			}
			cf.defaultVals[key] = newD
		}
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

func (v *dictValue[T]) ReplaceFromInterface(iFace interface{}) error {
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

func (v *dictValue[_]) Update(s string) error {
	key, strValue, found := strings.Cut(s, "=")
	if !found {
		return fmt.Errorf("Could not parse key=value for %v", s)
	}
	val, err := v.inner.FromString(strValue)
	if err != nil {
		return err
	}
	return v.update(key, val)
}

func (v *dictValue[_]) ReplaceFromDefault() {
	if v.hasDefault {
		v.vals = v.defaultVals
	}
}
