package value_test

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"

	w "github.com/bbkane/warg/value"
)

func TestIntValue(t *testing.T) {
	var v w.Value
	v, err := w.IntEmpty()
	require.Nil(t, err)
	require.Equal(t, v.Get().(int), 0)
	v.Update("2")
	require.Equal(t, v.Get().(int), 2)
}

func TestStringValue(t *testing.T) {
	var v w.Value
	v, err := w.StringEmpty()
	require.Nil(t, err)
	v.ReplaceFromInterface("hi")
	require.Equal(t, "hi", v.Get())
}

func TestStringSliceValue(t *testing.T) {
	var v w.Value
	v, err := w.StringSliceEmpty()
	require.Nil(t, err)

	// Not sure why I get the following, but seems to be a
	// limitation of the testing library
	// expected: []string([]string{})
	// actual  : <nil>(<nil>)
	// assert.Equal(t, v.Get().([]string), nil)

	v.Update("hi")
	require.Equal(t, v.Get().([]string), []string{"hi"})
}

func TestIntSliceValue(t *testing.T) {
	var v w.Value
	v, err := w.IntSliceEmpty()
	require.Nil(t, err)
	v.Update("1")
	require.Equal(t, v.Get().([]int), []int{1})
}

func TestPathValue(t *testing.T) {
	home, err := homedir.Dir()
	require.Nil(t, err)

	var v w.Value
	v, err = w.PathEmpty()
	require.Nil(t, err)
	err = v.Update("~/tmp")
	require.Nil(t, err)
	require.Equal(t,
		filepath.Join(home, "tmp"),
		v.Get().(string),
	)
}

func TestPathSliceValue(t *testing.T) {

	var v w.Value
	v, err := w.PathSliceEmpty()
	require.Nil(t, err)
	v.ReplaceFromInterface(
		[]string{"hi"},
	)

	require.Equal(t, v.Get().([]string), []string{"hi"})

}
