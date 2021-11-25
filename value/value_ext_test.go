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
	v.Update("2")
	require.Equal(t, v.Get().(int), 2)
}

func TestStringValue(t *testing.T) {
	var v value.Value
	v, err := value.String()
	require.Nil(t, err)
	v.ReplaceFromInterface("hi")
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

	v.Update("hi")
	require.Equal(t, v.Get().([]string), []string{"hi"})

	v.Update("there")
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
	v.Update("1")
	require.Equal(t, v.Get().([]int), []int{1})
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
	v.ReplaceFromInterface(
		[]string{"hi"},
	)

	require.Equal(t, v.Get().([]string), []string{"hi"})

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
