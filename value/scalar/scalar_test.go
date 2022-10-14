package scalar_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value/scalar"
)

func TestDurationString(t *testing.T) {
	constructor := scalar.Duration(scalar.Default(3 * time.Minute))
	v, err := constructor()
	require.Nil(t, err)

	vStr := v.DefaultString()
	require.Equal(t, "3m0s", vStr)
}
