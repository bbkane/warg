package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	w "github.com/bbkane/warg/value"
)

func TestIntValue(t *testing.T) {
	var v w.Value = w.IntNew(1)
	assert.Equal(t, v.Get().(int), 1)
	v.Update("2")
	assert.Equal(t, v.Get().(int), 2, "IntValue should be equal")
}

func TestStringSliceValue(t *testing.T) {
	var v w.Value = w.StringSliceNew([]string{})

	// Not sure why I get the following, but seems to be a
	// limitation of the testing library
	// expected: []string([]string{})
	// actual  : <nil>(<nil>)
	// assert.Equal(t, v.Get().([]string), nil)

	v.Update("hi")
	assert.Equal(t, v.Get().([]string), []string{"hi"})
}

func TestIntSliceValue(t *testing.T) {
	var v w.Value = w.IntSliceNew([]int{})
	v.Update("1")
	assert.Equal(t, v.Get().([]int), []int{1})
}
