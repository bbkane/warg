package value

import "fmt"

type commonFieldsOpt[T comparable] func(*commonFields[T])

type commonFields[T comparable] struct {
	choices []T
}

func Choices[T comparable](choices ...T) commonFieldsOpt[T] {
	return func(cf *commonFields[T]) {
		cf.choices = choices
	}
}

func (cf *commonFields[T]) WithinChoices(val T) bool {
	// User didn't constrain choices
	if len(cf.choices) == 0 {
		return true
	}
	for _, choice := range cf.choices {
		if val == choice {
			return true
		}
	}
	return false
}

func (cf *commonFields[T]) Choices() []string {
	ret := []string{}
	for _, e := range cf.choices {
		ret = append(ret, fmt.Sprint(e))
	}
	return ret
}
