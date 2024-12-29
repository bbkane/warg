package value_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

func TestIntScalar(t *testing.T) {
	v := scalar.Int(
		scalar.Choices(1, 2),
		scalar.Default(2),
	)()

	err := v.Update("1", value.UpdatedByFlag)
	require.Nil(t, err)
	require.Equal(t, v.Get().(int), 1)

	err = v.Update("-1", value.UpdatedByFlag)
	require.NotNil(t, err)
	require.Equal(t, v.Get().(int), 1)

	v.ReplaceFromDefault(value.UpdatedByDefault)
	require.Equal(t, v.Get().(int), 2)
}

func TestIntSlice(t *testing.T) {

	v := slice.Int(
		slice.Choices(1, 2),
		slice.Default([]int{1, 1, 1}),
	)()

	err := v.Update("1", value.UpdatedByFlag)
	require.Nil(t, err)
	require.Equal(
		t,
		[]int{1},
		v.Get().([]int),
	)

	err = v.Update("-1", value.UpdatedByFlag)
	require.NotNil(t, err)
	require.Equal(
		t,
		v.Get().([]int),
		[]int{1},
	)

	err = v.ReplaceFromInterface(
		[]interface{}{1, 2},
		value.UpdatedByFlag,
	)
	require.Nil(t, err)
	require.Equal(
		t,
		[]int{1, 2},
		v.Get().([]int),
	)

	v.ReplaceFromDefault(value.UpdatedByFlag)
	require.Equal(
		t,
		[]int{1, 1, 1},
		v.Get().([]int),
	)
}
