package scalar_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
	"go.bbkane.com/warg/value/scalar"
)

func TestDurationString(t *testing.T) {
	constructor := scalar.Duration(scalar.Default(3 * time.Minute))
	v := constructor()
	scalarValue := v.(value.ScalarValue)

	vStr := scalarValue.DefaultString()
	require.Equal(t, "3m0s", vStr)
}

func TestDefaultAndChoices(t *testing.T) {
	typeInfo := contained.Int()
	typeInfo.Description = "Defaults to the perfect number 7"

	constructor := scalar.New(typeInfo, scalar.Default(3), scalar.Choices(1, 2))
	v := constructor()
	scalarValue := v.(value.ScalarValue)

	actualDefaultStr := scalarValue.DefaultString()
	require.Equal(t, "3", actualDefaultStr)
	actualChoices := v.Choices()
	require.Equal(t, []string{"1", "2"}, actualChoices)
}
