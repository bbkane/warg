package value_test

import (
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/mitchellh/go-homedir"

	"go.bbkane.com/warg/value"
)

func TestIntValue(t *testing.T) {
	var v value.Value
	v, err := value.Int()
	assert.Nil(t, err)
	assert.Equal(t, v.Get().(int), 0)

	err = v.Update("2")
	assert.Nil(t, err)
	assert.Equal(t, v.Get().(int), 2)
}

func TestStringValue(t *testing.T) {
	var v value.Value
	v, err := value.String()
	assert.Nil(t, err)

	err = v.ReplaceFromInterface("hi")
	assert.Nil(t, err)
	assert.Equal(t, "hi", v.Get())
}

func TestStringSliceValue(t *testing.T) {
	var v value.Value
	v, err := value.StringSlice()
	assert.Nil(t, err)

	// Not sure why I get the following, but seems to be a
	// limitation of the testing library
	// expected: []string([]string{})
	// actual  : <nil>(<nil>)
	// assert.Equal(t, v.Get().([]string), nil)

	err = v.Update("hi")
	assert.Nil(t, err)
	assert.Equal(t, v.Get().([]string), []string{"hi"})

	err = v.Update("there")
	assert.Nil(t, err)
	assert.Equal(
		t,
		v.StringSlice(),
		[]string{"hi", "there"},
	)
}

func TestIntSliceValue(t *testing.T) {
	var v value.Value

	v, err := value.IntSlice()
	assert.Nil(t, err)

	err = v.Update("1")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]int{1},
		v.Get().([]int),
	)

	err = v.ReplaceFromInterface(
		[]interface{}{1, 2},
	)
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]int{1, 2},
		v.Get().([]int),
	)
}

func TestPathValue(t *testing.T) {
	home, err := homedir.Dir()
	assert.Nil(t, err)

	var v value.Value
	v, err = value.Path()
	assert.Nil(t, err)

	err = v.Update("~/tmp")
	assert.Nil(t, err)
	assert.Equal(t,
		filepath.Join(home, "tmp"),
		v.Get().(string),
	)
}

func TestPathSliceValue(t *testing.T) {

	var v value.Value
	v, err := value.PathSlice()
	assert.Nil(t, err)

	err = v.ReplaceFromInterface(
		[]interface{}{"hi", "there"},
	)
	assert.Nil(t, err)
	assert.Equal(t, v.Get().([]string), []string{"hi", "there"})
}

func TestStringEnumV(t *testing.T) {
	v, err := value.StringEnum("a", "b", "c")()
	assert.Nil(t, err)

	err = v.Update("a")
	assert.Nil(t, err)
	assert.Equal(t, v.Get().(string), "a")

	err = v.Update("notachoice")
	assert.NotNil(t, err)
}

func TestBoolV(t *testing.T) {
	v, err := value.Bool()
	assert.Nil(t, err)

	err = v.Update("true")
	assert.Nil(t, err)
	assert.Equal(t, v.Get().(bool), true)

	err = v.Update("bob")
	assert.NotNil(t, err)
}

func TestDurationV(t *testing.T) {
	v, err := value.Duration()
	assert.Nil(t, err)

	err = v.Update("2m")
	assert.Nil(t, err)

	err = v.ReplaceFromInterface("4m")
	assert.Nil(t, err)
}
