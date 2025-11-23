package contained_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/value/contained"
)

func TestTypeInfo_ValidateNonNilFuncs(t *testing.T) {
	require := require.New(t)
	require.NoError(contained.NetIPAddr().ValidateNonNilFuncs())
	require.NoError(contained.AddrPort().ValidateNonNilFuncs())
	require.NoError(contained.Bool().ValidateNonNilFuncs())
	require.NoError(contained.Duration().ValidateNonNilFuncs())
	require.NoError(contained.Int().ValidateNonNilFuncs())
	require.NoError(contained.Path().ValidateNonNilFuncs())
	require.NoError(contained.Rune().ValidateNonNilFuncs())
	require.NoError(contained.String().ValidateNonNilFuncs())
}
