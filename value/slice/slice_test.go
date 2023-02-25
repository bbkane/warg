package slice_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
	"go.bbkane.com/warg/value/slice"
)

func TestDefaultAndChoices(t *testing.T) {
	typeInfo := contained.Int()
	typeInfo.Description = "Defaults to the perfect number 7"
	typeInfo.FromInstance = func(i int) (int, error) {
		return 7, nil
	}

	constructor := slice.New(typeInfo, slice.Default([]int{3}), slice.Choices(1, 2))
	v, err := constructor()
	require.Nil(t, err)
	sliceVal := v.(value.SliceValue)

	actualDefaultStr := sliceVal.DefaultStringSlice()
	require.Equal(t, []string{"7"}, actualDefaultStr)
	actualChoices := v.Choices()
	require.Equal(t, []string{"7", "7"}, actualChoices)
}
