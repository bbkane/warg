package value_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	value "go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

func TestIntScalar(t *testing.T) {
	v, err := scalar.Int(
		scalar.Choices(1, 2),
		scalar.Default(2),
	)()
	require.Nil(t, err)

	err = v.Update("1")
	require.Nil(t, err)
	require.Equal(t, v.Get().(int), 1)

	err = v.Update("-1")
	require.NotNil(t, err)
	require.Equal(t, v.Get().(int), 1)

	v.ReplaceFromDefault()
	require.Equal(t, v.Get().(int), 2)
}

func TestIntSlice(t *testing.T) {
	var v value.Value

	v, err := slice.Int(
		slice.Choices(1, 2),
		slice.Default([]int{1, 1, 1}),
	)()
	require.Nil(t, err)

	err = v.Update("1")
	require.Nil(t, err)
	require.Equal(
		t,
		[]int{1},
		v.Get().([]int),
	)

	err = v.Update("-1")
	require.NotNil(t, err)
	require.Equal(
		t,
		v.Get().([]int),
		[]int{1},
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

	v.ReplaceFromDefault()
	require.Equal(
		t,
		[]int{1, 1, 1},
		v.Get().([]int),
	)
}
