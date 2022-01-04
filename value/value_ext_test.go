package value_test

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"

	"github.com/bbkane/warg/value"
)

func TestIntValue(t *testing.T) {
	var v value.Value
	v, err := value.Int()
	require.Nil(t, err)
	require.Equal(t, v.Get().(int), 0)

	err = v.Update("2")
	require.Nil(t, err)
	require.Equal(t, v.Get().(int), 2)
}

func TestStringValue(t *testing.T) {
	var v value.Value
	v, err := value.String()
	require.Nil(t, err)

	err = v.ReplaceFromInterface("hi")
	require.Nil(t, err)
	require.Equal(t, "hi", v.Get())
}

func TestStringSliceValue(t *testing.T) {
	var v value.Value
	v, err := value.StringSlice()
	require.Nil(t, err)

	// Not sure why I get the following, but seems to be a
	// limitation of the testing library
	// expected: []string([]string{})
	// actual  : <nil>(<nil>)
	// assert.Equal(t, v.Get().([]string), nil)

	err = v.Update("hi")
	require.Nil(t, err)
	require.Equal(t, v.Get().([]string), []string{"hi"})

	err = v.Update("there")
	require.Nil(t, err)
	require.Equal(
		t,
		v.StringSlice(),
		[]string{"hi", "there"},
	)
}

func TestIntSliceValue(t *testing.T) {
	var v value.Value

	v, err := value.IntSlice()
	require.Nil(t, err)

	err = v.Update("1")
	require.Nil(t, err)
	require.Equal(
		t,
		[]int{1},
		v.Get().([]int),
	)

	err = v.ReplaceFromInterface(
		[]interface{}{1, 2},
	)
	require.Nil(t, err)
	require.Equal(
		t,
		[]int{1, 2},
		v.Get().([]int),
	)
}

func TestPathValue(t *testing.T) {
	home, err := homedir.Dir()
	require.Nil(t, err)

	var v value.Value
	v, err = value.Path()
	require.Nil(t, err)

	err = v.Update("~/tmp")
	require.Nil(t, err)
	require.Equal(t,
		filepath.Join(home, "tmp"),
		v.Get().(string),
	)
}

func TestPathSliceValue(t *testing.T) {

	var v value.Value
	v, err := value.PathSlice()
	require.Nil(t, err)

	err = v.ReplaceFromInterface(
		[]interface{}{"hi", "there"},
	)
	require.Nil(t, err)
	require.Equal(t, v.Get().([]string), []string{"hi", "there"})
}

func TestStringEnumV(t *testing.T) {
	v, err := value.StringEnum("a", "b", "c")()
	require.Nil(t, err)

	err = v.Update("a")
	require.Nil(t, err)
	require.Equal(t, v.Get().(string), "a")

	err = v.Update("notachoice")
	require.NotNil(t, err)
}

func TestBoolV(t *testing.T) {
	v, err := value.Bool()
	require.Nil(t, err)

	err = v.Update("true")
	require.Nil(t, err)
	require.Equal(t, v.Get().(bool), true)

	err = v.Update("bob")
	require.NotNil(t, err)
}

func TestDurationV(t *testing.T) {
	v, err := value.Duration()
	require.Nil(t, err)

	err = v.Update("2m")
	require.Nil(t, err)

	err = v.ReplaceFromInterface("4m")
	require.Nil(t, err)
}
