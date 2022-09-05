package value_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	value "go.bbkane.com/warg/value"
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

func TestIntEnum(t *testing.T) {
	v, err := value.IntEnum(1, 2)()
	require.Nil(t, err)

	err = v.Update("1")
	require.Nil(t, err)
	require.Equal(t, v.Get().(int), 1)

	err = v.Update("-1")
	require.NotNil(t, err)
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
